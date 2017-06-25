package main

import (
	"encoding/xml"
	"flag"
	"log"
	"net/http"
)

type Resp struct {
	XMLName xml.Name `xml:"Response"`
	Message Message
}

type Message struct {
	Body string `xml:"Body"`
	From string `xml:"From,attr"`
	To   string `xml:"To,attr"`
}

func main() {

	// Config

	var accountSid string
	var num1, num2 string
	flag.StringVar(&accountSid, "accountsid", "", "Twilio Account Sid")
	flag.StringVar(&num1, "num1", "", "First phone number")
	flag.StringVar(&num2, "num2", "", "Second phone number")
	flag.Parse()
	if accountSid == "" || num1 == "" || num2 == "" {
		log.Fatalf("accountsid, token, num1, num2 are required CLI parameters\n")
	}

	//
	// Handlers
	//

	// Handle and forward an incoming SMS message
	smsHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported method", 405)
			return
		}
		log.Println(r.URL.Query())
		log.Println(r.Header)

		from := r.URL.Query().Get("From")
		to := r.URL.Query().Get("To")
		body := r.URL.Query().Get("Body")
		toAccountSid := r.URL.Query().Get("AccountSid")
		if accountSid != toAccountSid {
			http.Error(w, "bad twilio account", 400)
			return
		}

		var forwardNum string
		switch from {
		case num1:
			forwardNum = num2
		case num2:
			forwardNum = num1
		default:
			http.Error(w, "Bad request", 400)
			return
		}
		resp := Resp{Message: Message{Body: body, From: to, To: forwardNum}}
		bytes, err := xml.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/xml")
		w.Write(bytes)
	}

	// Service healthcheck
	healthcheckHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported method", 405)
			return
		}
		w.Write([]byte("HEALTHCHECK OK"))
	}

	http.HandleFunc("/sms", smsHandler)
	http.HandleFunc("/healthcheck", healthcheckHandler)
	http.ListenAndServeTLS(":443", "cert.pem", "privkey.pem", nil)
}
