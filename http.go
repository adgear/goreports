// Copyright (c) 2014 Datacratic. All rights reserved.

package report

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

// RequestHTTP creates the binary data that represents the HTTP request.
func RequestHTTP(r *http.Request, body []byte) (data Data) {
	buffer := new(bytes.Buffer)

	s := fmt.Sprintf("%s %s %s\r\n", r.Method, r.RequestURI, r.Proto)
	if _, err := buffer.WriteString(s); err != nil {
		log.Panic(err.Error())
	}

	if err := r.Header.Write(buffer); err != nil {
		log.Panic(err.Error())
	}

	if _, err := buffer.WriteString("\r\n"); err != nil {
		log.Panic(err.Error())
	}

	if _, err := buffer.Write(body); err != nil {
		log.Panic(err.Error())
	}

	data.Name = "request"
	data.Blob = buffer.Bytes()
	return
}

// ResponseHTTP creates the binary data that represents the HTTP response.
func ResponseHTTP(r *http.Response, body []byte) (data Data) {
	buffer := new(bytes.Buffer)

	s := fmt.Sprintf("%s %s\r\n", r.Status, r.Proto)
	if _, err := buffer.WriteString(s); err != nil {
		log.Panic(err.Error())
	}

	if err := r.Header.Write(buffer); err != nil {
		log.Panic(err.Error())
	}

	if _, err := buffer.WriteString("\r\n"); err != nil {
		log.Panic(err.Error())
	}

	if _, err := buffer.Write(body); err != nil {
		log.Panic(err.Error())
	}

	data.Name = "response"
	data.Blob = buffer.Bytes()
	return
}
