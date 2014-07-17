// Copyright (c) 2014 Datacratic. All rights reserved.

package report

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Publisher is an interface to send reports.
type Publisher interface {
	Send(r *Report)
}

// Reporter implements a service where reports can be published.
type Reporter struct {
	// Name contains a friendly identifier for reporting.
	Name string
	// Publisher represents the handler that will publish reports.
	Publisher Publisher

	reports chan *Report
}

// Start creates the background service that publish reports.
func (reporter *Reporter) Start() (err error) {
	if reporter.Name == "" {
		err = fmt.Errorf("missing reporter name")
		return
	}

	if reporter.Publisher == nil {
		reporter.PublishFunc(func(r *Report) {
			r.Write(os.Stdout)
		})
	}

	reporter.reports = make(chan *Report)

	go func() {
		for item := range reporter.reports {
			reporter.Publisher.Send(item)
		}
	}()

	return
}

// Log creates a report based on the supplied text along with some associated binary data.
func (reporter *Reporter) Log(name, s string, data ...Data) {
	reporter.write(name, s, data)
}

// Error creates a report based on the supplied error along with some associated binary data.
func (reporter *Reporter) Error(name string, e error, data ...Data) {
	reporter.write(name, e.Error(), data)
}

func (reporter *Reporter) write(name, s string, data []Data) {
	r := Report{
		Reporter:  reporter.Name,
		Time:      time.Now(),
		Component: name,
		Status:    s,
		Content:   data,
	}

	reporter.reports <- &r
}

// Publish sets the publish handler.
func (reporter *Reporter) Publish(handler Publisher) {
	reporter.Publisher = handler
}

// PublishFunc defines an helper to support the Publisher interface.
type PublishFunc func(r *Report)

// Send invokes the function literal with the report.
func (f PublishFunc) Send(r *Report) {
	f(r)
}

// PublishFunc is an helper that sets a function literal as the publish handler.
func (reporter *Reporter) PublishFunc(handler func(r *Report)) {
	reporter.Publish(PublishFunc(handler))
}

// NewJSONReporter creates and starts a publisher that POST the report as a JSON object to the specified URL followed by the associated binary data.
func NewJSONReporter(name string, url string) (reporter *Reporter) {
	reporter = &Reporter{
		Name: name,
	}

	reporter.PublishFunc(func(r *Report) {
		buf := new(bytes.Buffer)
		r.WriteJSON(buf)

		rep, err := http.Post(url, "application/json", buf)
		if err != nil {
			log.Printf("failed to send report: %s", err.Error())
			return
		}

		ioutil.ReadAll(rep.Body)
		rep.Body.Close()
	})

	reporter.Start()
	return
}
