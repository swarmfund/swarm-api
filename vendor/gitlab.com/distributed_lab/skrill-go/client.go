package skrill

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"encoding/xml"

	"github.com/pkg/errors"
)

const (
	skrillPayEndpoint         = "https://pay.skrill.com"
	skrillSendEndpoint        = "https://www.skrill.com/app/pay.pl"
	skrillQueryEndpoint       = "https://www.skrill.com/app/query.pl"
	contentTypeJSON           = "application/json"
	contentTypeFormUrlencoded = "application/x-www-form-urlencoded"
	precision                 = 4
)

var (
	ErrBadRequest = errors.New("bad request")
)

type MerchantFielder interface {
	MerchantFields() map[string]string
}

type Client struct {
	options map[string]string

	email         string
	password      string
	requestTicker *time.Ticker
}

func NewClient() *Client {
	return &Client{
		options:       map[string]string{},
		requestTicker: time.NewTicker(10 * time.Second),
	}
}

func NewMerchantClient(email, password string) *Client {
	return &Client{
		options:       map[string]string{},
		email:         email,
		password:      password,
		requestTicker: time.NewTicker(10 * time.Second),
	}
}

func (c *Client) Transaction(params *TransactionParams) (*StatusTransaction, error) {
	params.Action = "status_trn"
	if params.Password == "" {
		params.Password = c.password
	}
	if params.Email == "" {
		params.Email = c.email
	}
	body, err := c.post(skrillQueryEndpoint, params)
	if err != nil {
		return nil, err
	}

	return NewStatusTransaction(bytes.NewReader(body))
}

func (c *Client) PrepareSend(params *PrepareSendParams) (*ExecuteSendParams, error) {
	params.Action = "prepare"
	if params.Password == "" {
		params.Password = c.password
	}
	if params.Email == "" {
		params.Email = c.email
	}
	body, err := c.post(skrillSendEndpoint, params)
	if err != nil {
		return nil, err
	}

	response := ExecuteSendParams{}
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) ExecuteSend(params *ExecuteSendParams) error {
	params.Action = "transfer"
	body, err := c.post(skrillSendEndpoint, params)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	return nil
}

func (c *Client) History(params *HistoryParams) (chan HistoryTransaction, chan error) {
	txCh := make(chan HistoryTransaction)
	errCh := make(chan error)

	go func() {
		func() {
			params.Action = "history"
			if params.Email == "" {
				params.Email = c.email
			}
			if params.Password == "" {
				params.Password = c.password
			}

			body, err := c.get(skrillQueryEndpoint, params)
			if err != nil {
				errCh <- err
				return
			}
			defer body.Close()

			bodyBytes, err := ioutil.ReadAll(body)
			if err != nil {
				errCh <- err
				return
			}

			r := csv.NewReader(bytes.NewReader(bodyBytes))
			r.LazyQuotes = true

			// skip header
			_, err = r.Read()
			if err == io.EOF {
				// kinda unexpected EOF
				errCh <- nil
				return
			}
			if err != nil {
				errCh <- err
				return
			}

			for {
				record, err := r.Read()
				if err == io.EOF {
					errCh <- nil
					break
				}
				if err != nil {
					errCh <- err
					break
				}

				// Time
				txTime, err := time.Parse(time.RFC822, fmt.Sprintf("%s CET", record[1]))
				if err != nil {
					errCh <- err
					return
				}

				// Type
				txType, ok := SkrillTypeTransactionType[record[2]]
				if !ok {
					errCh <- fmt.Errorf("unknown tx type %s", record[2])
					return
				}

				// SendUSD
				var txSendUSD int64 = 0
				if record[4] != "" {
					txSendUSD, err = ParseAmount(record[4], precision)
					if err != nil {
						errCh <- err
						return
					}
				}

				// ReceivedUSD
				var txReceivedUSD int64 = 0
				if record[5] != "" {
					txReceivedUSD, err = ParseAmount(record[5], precision)
					if err != nil {
						errCh <- err
						return
					}
				}

				// Status
				txStatus, ok := SkrillStatusTransactionStatus[record[6]]
				if !ok {
					errCh <- fmt.Errorf("unknown tx status %s", record[6])
					return
				}

				// Balance
				var txBalance int64 = 0
				if record[7] != "" {
					txBalance, err = ParseAmount(record[7], precision)
					if err != nil {
						errCh <- err
						return
					}
				}

				// AmountSent
				var txAmountSent int64 = 0
				if record[9] != "" {
					txAmountSent, err = ParseAmount(record[9], precision)
					if err != nil {
						errCh <- err
						return
					}
				}

				tx := HistoryTransaction{
					record[0],
					txTime,
					txType,
					record[3],
					txSendUSD,
					txReceivedUSD,
					txStatus,
					txBalance,
					record[8],
					txAmountSent,
					record[10],
					record[11],
					record[12],
					record[13],
				}
				txCh <- tx
			}
		}()
		close(errCh)
		close(txCh)
	}()

	return txCh, errCh
}

func (c *Client) get(endpoint string, params interface{}) (io.ReadCloser, error) {
	<-c.requestTicker.C
	// going struct -> json -> map -> form -> string
	paramsBytes, err := json.Marshal(&params)
	if err != nil {
		return nil, err
	}

	paramsMap := map[string]interface{}{}
	err = json.Unmarshal(paramsBytes, &paramsMap)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	for key, value := range paramsMap {
		// TODO type cast non-strings
		switch v := value.(type) {
		case string:
			form.Add(key, v)
		default:
			return nil, ErrBadRequest
		}
	}

	url, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	url.RawQuery = form.Encode()

	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200:

		return res.Body, nil
	default:
		fmt.Println(res.Status)
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(body))
		return nil, ErrBadRequest
	}
}

func (c *Client) post(endpoint string, params interface{}) ([]byte, error) {
	<-c.requestTicker.C
	// going struct -> json -> map -> form -> string
	paramsBytes, err := json.Marshal(&params)
	if err != nil {
		return nil, err
	}

	paramsMap := map[string]interface{}{}
	err = json.Unmarshal(paramsBytes, &paramsMap)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	for key, value := range paramsMap {
		// TODO type cast non-strings
		switch v := value.(type) {
		case string:
			form.Add(key, v)
		default:
			return nil, ErrBadRequest
		}
	}

	for key, value := range c.options {
		form.Add(key, value)
	}

	if fielder, ok := params.(MerchantFielder); ok {
		for key, value := range fielder.MerchantFields() {
			form.Add(key, value)
		}
	}

	res, err := http.Post(endpoint, contentTypeFormUrlencoded, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		// let's check response for interesting errors
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var errorResponse ErrorResponse
		// silence error since we are just guessing here and not sure of actual response schema
		_ = xml.Unmarshal(body, &errorResponse)
		switch errorResponse.Error.Message {
		case ErrorMessageAlreadyExecuted:
			return nil, ErrAlreadyExecuted
		case ErrorMessageBalanceNotEnough:
			return nil, ErrBalanceNotEnough
		}

		fmt.Println(string(body))

		return body, nil
	default:
		return nil, fmt.Errorf("request failed: %s", res.Status)
	}
}

type PrepareSendParams struct {
	Action   string `json:"action"`
	Email    string `json:"email"` // merchant email
	Password string `json:"password"`
	// required fields
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	BnfEmail string `json:"bnf_email"` // recipient email
	Subject  string `json:"subject"`
	Note     string `json:"note"`
	// optional
	TransactionID string `json:"frn_trn_id"`
}

type ExecuteSendParams struct {
	Action string `json:"action"`
	SID    string `json:"sid" xml:"sid"`
}

func (p *QuickCheckoutParams) verify() bool {
	// TODO verify all required fields are set
	return true
}

type HistoryParams struct {
	// required fields
	Email    string `json:"email"`
	Password string `json:"password"`
	Start    string `json:"start_date"`
	// optional
	End    string `json:"end_date,omitempty"`
	Action string `json:"action"`
}

type TransactionParams struct {
	// required fields
	Email    string `json:"email"`
	Password string `json:"password"`
	// one of
	TransactionID       string `json:"trn_id"`
	SkrillTransactionID string `json:"mb_trn_id"`
	//optional
	Action string `json:"action"`
}

type ErrorMessage string

var (
	ErrorMessageBalanceNotEnough ErrorMessage = "BALANCE_NOT_ENOUGH"
	ErrorMessageAlreadyExecuted  ErrorMessage = "ALREADY_EXECUTED"
)

var (
	ErrBalanceNotEnough = errors.New(string(ErrorMessageBalanceNotEnough))
	ErrAlreadyExecuted  = errors.New(string(ErrorMessageAlreadyExecuted))
)

type ErrorResponse struct {
	Error struct {
		Message ErrorMessage `xml:"error_msg"`
	} `xml:"error"`
}
