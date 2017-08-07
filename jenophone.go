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
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
}

func main() {

	// Config

	var num1, num2 string
	flag.StringVar(&num1, "num1", "", "First phone number")
	flag.StringVar(&num2, "num2", "", "Second phone number")
	flag.Parse()
	if num1 == "" || num2 == "" {
		log.Fatalf("num1, num2 are required CLI parameters\n")
	}

	// Handlers

	smsHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported method", 405)
			return
		}
		from := r.URL.Query().Get("From")
		to := r.URL.Query().Get("To")
		body := r.URL.Query().Get("Body")
		log.Println(r.URL.Query())
		log.Println(r.Header)

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
		log.Printf("Received from %v and going to %v\n", from, forwardNum)
		resp := Resp{Message: Message{Body: body, From: to, To: forwardNum}}
		bytes, err := xml.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println("Sending back", string(bytes))
		w.Header().Set("Content-Type", "application/xml")
		w.Write(bytes)
	}

	healthcheckHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Unsupported method", 405)
			return
		}
		w.Write([]byte("HEALTHCHECK OK"))
	}

	http.HandleFunc("/sms", smsHandler)
	http.HandleFunc("/healthcheck", healthcheckHandler)
	http.ListenAndServeTLS(":443", "fullchain.pem", "privkey.pem", nil)
}
