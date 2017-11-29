package skrill

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
)

type StatusTransaction struct {
	TransactionID       string            // transaction_id=2116016081 tx id provided during by horizon
	MerchantAmount      int64             // mb_amount=10.2 tx amount in merchant account currency
	PaymentType         string            // payment_type=VSA
	Amount              int64             // amount=10.20
	SkrillTransactionID string            // mb_transaction_id=2116016081 skrill tx id
	MerchantCurrency    string            // mb_currency=USD currency of merchant account
	PayFromEmail        string            // pay_from_email=foo%40bar.com
	MD5sig              string            // md5sig = E3BCE44891C970A295FE062C25005F4D TODO verify
	PayToEmail          string            // pay_to_email=comrad.awsum%2Bmerchant%40gmail.com
	Currency            string            // currency=USD
	MerchantID          string            // merchant_id=93554219
	Status              TransactionStatus // status=2

	customFields map[string]string // as a convention we are using fields starting with x-
}

func NewStatusTransaction(body io.Reader) (*StatusTransaction, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// skrill kindly include super convenient `200 OK` in *some* responses
	split := strings.Split(string(data), "\n")
	var dataString string
	switch len(split) {
	case 1:
		dataString = split[0]
	case 3:
		dataString = split[1]
	default:
		return nil, errors.New("unknown body format")
	}

	values, err := url.ParseQuery(dataString)
	if err != nil {
		return nil, err
	}

	var merchantAmount int64 = 0
	if values.Get("mb_amount") != "" {
		merchantAmount, err = ParseAmount(values.Get("mb_amount"), precision)
		if err != nil {
			return nil, err
		}
	}

	var amount int64 = 0
	if values.Get("amount") != "" {
		amount, err = ParseAmount(values.Get("amount"), precision)
		if err != nil {
			return nil, err

		}
	}

	status, ok := SkrillStatusTransactionStatus[values.Get("status")]
	if !ok {
		return nil, fmt.Errorf("unknown tx status %s", values.Get("status"))
	}

	transaction := StatusTransaction{
		values.Get("transaction_id"),
		merchantAmount,
		values.Get("payment_type"),
		amount,
		values.Get("mb_transaction_id"),
		values.Get("mb_currency"),
		values.Get("pay_from_email"),
		values.Get("md5sig"),
		values.Get("pay_to_email"),
		values.Get("currency"),
		values.Get("merchant_id"),
		status,
		map[string]string{},
	}

	for key, _ := range values {
		if strings.HasPrefix(key, "x-") {
			transaction.customFields[strings.TrimPrefix(key, "x-")] = values.Get(key)
		}
	}

	return &transaction, nil
}

func (tx *StatusTransaction) CustomField(key string) string {
	return tx.customFields[key]
}
