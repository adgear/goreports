// Copyright (c) 2014 Datacratic. All rights reserved.

package report

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Report defines basic informations that are required for reports.
type Report struct {
	// Time contains the time the report was created.
	Time time.Time `json:"time"`
	// Component contains the name of the reporter producing the report.
	Component string `json:"component"`
	// Status contains the text of the report.
	Status string `json:"status"`
	// Content contains the binary data that can be associated with the report.
	Content []Data `json:"content,omitempty"`
}

// Data defines a binary data entry that can be associated with reports.
type Data struct {
	// ID contains an identifier for this entry.
	Name string
	// Bytes contains the binary data associated with the report.
	Bytes []byte
}

// Write outputs the report as text followed by the associated binary data separated by new lines.
func (report *Report) Write(w io.Writer) {
	text := fmt.Sprintf("time: %s\ncomponent: %s\nstatus: %s\n", report.Time, report.Component, report.Status)
	io.WriteString(w, text)

	if len(report.Content) != 0 {
		io.WriteString(w, "content:")

		for i, item := range report.Content {
			text := fmt.Sprintf("\n %d. '%s' has %d bytes", i, item.Name, len(item.Bytes))
			io.WriteString(w, text)
		}

		io.WriteString(w, "\n")

		for _, item := range report.Content {
			w.Write(item.Bytes)
			io.WriteString(w, "\n")
		}
	}
}

// MarshalJSON encodes the id and lenght of the binary data without the actual data.
func (data *Data) MarshalJSON() ([]byte, error) {
	s := struct {
		ID     string `json:"id"`
		Length int    `json:"len"`
	}{
		ID:     data.Name,
		Length: len(data.Bytes),
	}

	return json.Marshal(s)
}

// UnmarshalJSON decodes the id and length of the binary data without the actual data.
func (data *Data) UnmarshalJSON(bytes []byte) (err error) {
	s := new(struct {
		ID     string `json:"id"`
		Length int    `json:"len"`
	})

	if err = json.Unmarshal(bytes, s); err != nil {
		return
	}

	data.Name = s.ID
	data.Bytes = make([]byte, s.Length)
	return
}

// WriteJSON outputs the report as a JSON object followed by the associated binary data separated by new lines.
func (report *Report) WriteJSON(w io.Writer) {
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(report); err != nil {
		panic(err.Error())
	}

	for _, item := range report.Content {
		w.Write(item.Bytes)

		// separate binary data with a new line
		io.WriteString(w, "\n")
	}
}

// ReadJSON creates a report from a JSON object followed by its associated binary data separated by new lines.
func ReadJSON(r io.Reader) (report *Report, err error) {
	decoder := json.NewDecoder(r)

	data := new(Report)
	if err = decoder.Decode(data); err != nil {
		return
	}

	binary := io.MultiReader(decoder.Buffered(), r)

	c := make([]byte, 1)
	if _, err = binary.Read(c); err != nil || c[0] != '\n' {
		err = fmt.Errorf("expected new line character instead of '%v'", c)
		return
	}

	for _, item := range data.Content {
		if _, err = binary.Read(item.Bytes); err != nil {
			return
		}

		if _, err = binary.Read(c); err != nil || c[0] != '\n' {
			err = fmt.Errorf("expected new line character instead of '%v'", c)
			return
		}
	}

	report = data
	return
}
