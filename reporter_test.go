// Copyright (c) 2014 Datacratic. All rights reserved.

package report

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReporter(t *testing.T) {
	result := make(chan string)

	endpoint := func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err.Error())
		}

		t.Logf("received:\n%s", string(b))

		report, err := ReadJSON(bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err.Error())
		}

		result <- report.Status
		w.WriteHeader(http.StatusOK)
	}

	h := httptest.NewServer(http.HandlerFunc(endpoint))
	r := NewJSONReporter("test", h.URL)

	go func() {
		r.Log("test", "hello")

		i := struct {
			ID   int
			Name string
		}{
			ID:   0,
			Name: "some data",
		}

		some, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err.Error())
		}

		i.ID = 1
		i.Name = "some more data"

		more, err := json.Marshal(i)
		if err != nil {
			t.Fatal(err.Error())
		}

		r.Log("test", "brave", Data{"some", some})
		r.Log("test", "world", Data{"some", some}, Data{"some more", more})
	}()

	if s := <-result; s != "hello" {
		t.Fail()
	}

	if s := <-result; s != "brave" {
		t.Fail()
	}

	if s := <-result; s != "world" {
		t.Fail()
	}
}
