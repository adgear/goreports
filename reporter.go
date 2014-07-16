// Copyright (c) 2014 Datacratic. All rights reserved.

package report

import (
	"bytes"
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
}

// String creates a report based on the supplied text.
func (reporter *Reporter) String(s string) {
	reporter.write(s, nil)
}

// StringWithData creates a report based on the supplied text along with some associated binary data.
func (reporter *Reporter) StringWithData(s string, data ...Data) {
	reporter.write(s, data)
}

// Error creates a report based on the supplied error.
func (reporter *Reporter) Error(e error) {
	reporter.String(e.Error())
}

// ErrorWithData creates a report based on the supplied error along with some associated binary data.
func (reporter *Reporter) ErrorWithData(e error, data ...Data) {
	reporter.StringWithData(e.Error(), data...)
}

func (reporter *Reporter) write(s string, data []Data) {
	r := Report{
		Time:      time.Now(),
		Component: reporter.Name,
		Status:    s,
		Content:   data,
	}

	if reporter.Publisher == nil {
		r.Write(os.Stdout)
		return
	}

	reporter.Publisher.Send(&r)
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

// NewJSONReporter creates a publisher that POST the report as a JSON object to the specified URL followed by the associated binary data.
func NewJSONReporter(name string, url string) *Reporter {
	reporter := &Reporter{
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

	return reporter
}
