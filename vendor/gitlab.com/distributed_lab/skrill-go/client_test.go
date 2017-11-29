package skrill

import (
	"encoding/xml"
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

var (
	merchantEmail = "comrad.awsum+merchant@gmail.com"
	password      = "1cd5ea39086126ed4528e320a9282b2d"
)

func TestErrorResponse(t *testing.T) {
	var response ErrorResponse

	data := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><response><error><error_msg>BALANCE_NOT_ENOUGH</error_msg></error></response>`)
	err := xml.Unmarshal(data, &response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error.Message != ErrorMessageBalanceNotEnough {
		t.Fatalf("expected %s got %s", ErrorMessageBalanceNotEnough, response.Error.Message)
	}

	data = []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><response><error><error_msg>ALREADY_EXECUTED</error_msg></error></response>`)
	err = xml.Unmarshal(data, &response)
	if err != nil {
		t.Fatal(err)
	}
	if response.Error.Message != ErrorMessageAlreadyExecuted {
		t.Fatalf("expected %s got %s", ErrorMessageAlreadyExecuted, response.Error.Message)
	}
}

func TestSend(t *testing.T) {
	client := NewMerchantClient(merchantEmail, password)

	sid, err := client.PrepareSend(&PrepareSendParams{
		Amount:        "100",
		Currency:      "USD",
		BnfEmail:      "comrad.awsum+client1@gmail.com",
		Subject:       "bc fr",
		Note:          "note",
		TransactionID: "1",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.ExecuteSend(sid)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(sid)
}

func TestClient(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	txID := fmt.Sprintf("%d", rand.Int63n(math.MaxInt64))

	client := NewClient()

	sid, err := client.QuickCheckoutSession(&QuickCheckoutParams{
		PayToEmail:    merchantEmail,
		Amount:        "1000",
		StatusURL:     "http://006470c5.ngrok.io/skrill/status",
		Currency:      "USD",
		TransactionID: "SOOOOOOOOOOOOOOOOOOOOOOOKA",
		CustomFields: map[string]string{
			"x-receiver": "BDFI3OE5XRDYZIGPR75V4FNFNVRXP4PBDO6PXXCIG2UJYBUVXHBH2YT2",
			"x-asset":    "XUSD",
		},
	})

	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(txID, sid)

	txCh, doneCh := client.History(&HistoryParams{
		Email:    merchantEmail,
		Password: password,
		Start:    "01-01-1970",
	})

	func() {
		for {
			select {
			case htx := <-txCh:
				t.Logf("history: %+v\n", htx)
				stx, err := client.Transaction(&TransactionParams{
					SkrillTransactionID: htx.TransactionID,
					Email:               merchantEmail,
					Password:            password,
				})
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("status: %+v\n", stx)
				if htx.Type == TxTypeSend {
					fmt.Printf("YOBA %+v\n", htx)
					fmt.Printf("YOBA %+v\n", stx)
				}
			case err := <-doneCh:
				if err != nil {
					t.Fatal(err)
				}
				return
			}
		}
	}()
}

func TestAccountHistory(t *testing.T) {
	_ = `"ID","Time (CET)","Type","Transaction Details","[-] USD","[+] USD","Status","balance","Reference","Amount Sent","Currency sent","More information","ID of the corresponding Skrill transaction","Payment Type"
"2108363574","25 May 17 09:23","Receive Money","from testfundsreserve@skrill.com","","10000","processed","10000","","10000","USD","-","2108363559","WLT"
"2108363575","25 May 17 09:23","Receive Money","Fee","350","","processed","9650","","","","","2108363559",""
"2115886559","03 Jun 17 12:20","Receive Money","from comrad.awsum@gmail.com","","10.2","processed","9660.2","","10.2","USD"," ","2115886554","DIN"
"2115886560","03 Jun 17 12:20","Receive Money","Fee","1.1281","","processed","9659.0719","","","","","2115886554",""
"2115886561","03 Jun 17 12:20","Receive Money","Per Transaction Fee",".327149","","processed","9658.744751","","","","","2115886554",""
"2115898283","03 Jun 17 12:36","Receive Money","from comrad.awsum@gmail.com","","10.2","processed","9668.944751","","10.2","USD"," ","2115898279","VSA"
"2115898285","03 Jun 17 12:36","Receive Money","Fee","1.1281","","processed","9667.816651","","","","","2115898279",""
"2115898287","03 Jun 17 12:36","Receive Money","Per Transaction Fee",".327149","","processed","9667.489502","","","","","2115898279",""
"2115915599","03 Jun 17 13:00","Receive Money","from 11@11.com","","10.2","processed","9677.689502","","10.2","USD"," ","2115915596","VSA"
"2115915603","03 Jun 17 13:00","Receive Money","Fee","1.1281","","processed","9676.561402","","","","","2115915596",""
"2115915604","03 Jun 17 13:00","Receive Money","Per Transaction Fee",".327149","","processed","9676.234253","","","","","2115915596",""
"2116016090","03 Jun 17 15:14","Receive Money","from foo@bar.com","","10.2","processed","9686.434253","","10.2","USD"," ","2116016081","VSA"
"2116016094","03 Jun 17 15:14","Receive Money","Fee","1.1281","","processed","9685.306153","","","","","2116016081",""
"2116016095","03 Jun 17 15:14","Receive Money","Per Transaction Fee",".327149","","processed","9684.979004","","","","","2116016081",""
"2116024874","03 Jun 17 15:26","Receive Money","from comrad.awsum@gmail.com","","10.0001","processed","9694.979104","","10.01","USD"," ","2116024869","VSA"
"2116024879","03 Jun 17 15:26","Receive Money","Fee","1.1281","","processed","9693.851004","","","","","2116024869",""
"2116024886","03 Jun 17 15:26","Receive Money","Per Transaction Fee",".327149","","processed","9693.523855","","","","","2116024869",""
"2117429560","05 Jun 17 15:18","Receive Money","from foo@bar.com","","1000","processed","10693.523855","","1000","USD"," ","2117429557","VSA"
"2117429561","05 Jun 17 15:18","Receive Money","Fee","35","","processed","10658.523855","","","","","2117429557",""
"2117429562","05 Jun 17 15:18","Receive Money","Per Transaction Fee",".325989","","processed","10658.197866","","","","","2117429557",""
"2117430450","05 Jun 17 15:20","Receive Money","from foo@bar.com","","1000","processed","11658.197866","","1000","USD"," ","2117430443","AMX"
"2117430455","05 Jun 17 15:20","Receive Money","Fee","35","","processed","11623.197866","","","","","2117430443",""
"2117430458","05 Jun 17 15:20","Receive Money","Per Transaction Fee",".325989","","processed","11622.871877","","","","","2117430443",""
"2117437443","05 Jun 17 15:28","Receive Money","from dd@dd.dd","","1000","processed","12622.871877","","1000","USD"," ","2117437438","DIN"
"2117437444","05 Jun 17 15:28","Receive Money","Fee","35","","processed","12587.871877","","","","","2117437438",""
"2117437445","05 Jun 17 15:28","Receive Money","Per Transaction Fee",".325989","","processed","12587.545888","","","","","2117437438",""
"2117446614","05 Jun 17 15:40","Receive Money","from dd@dd.dd","","1000","processed","13587.545888","","1000","USD"," ","2117446609","DIN"
"2117446615","05 Jun 17 15:40","Receive Money","Fee","35","","processed","13552.545888","","","","","2117446609",""
"2117446616","05 Jun 17 15:40","Receive Money","Per Transaction Fee",".326149","","processed","13552.219739","","","","","2117446609",""
"2117448156","05 Jun 17 15:42","Receive Money","from dd@dd.dd","","1000","processed","14552.219739","","1000","USD"," ","2117448153","DIN"
"2117448157","05 Jun 17 15:42","Receive Money","Fee","35","","processed","14517.219739","","","","","2117448153",""
"2117448158","05 Jun 17 15:42","Receive Money","Per Transaction Fee",".326149","","processed","14516.89359","","","","","2117448153",""
"2118995190","07 Jun 17 11:37","Receive Money","from dd@dd.dd","","729.988453","processed","15246.882043","","700.83","EUR"," ","2118995184","MSC"
"2118995191","07 Jun 17 11:37","Receive Money","Fee","25.549596","","processed","15221.332447","","","","","2118995184",""
"2118995194","07 Jun 17 11:37","Receive Money","Per Transaction Fee",".326758","","processed","15221.005689","","","","","2118995184",""`
}
