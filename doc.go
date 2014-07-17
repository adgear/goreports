// Copyright (c) 2014 Datacratic. All rights reserved.

/*
Package report provides extended log support.

First, create a Reporter instance:

	r := Reporter{
		Name: "my reports",
	}

	r.Start()

Then, simply send logs or errors.

	// basic logs
	r.Log("server", "everything is well within parameters")

	// and errors
	r.Error("client", err)

Reports can be build from text or errors. They can also contain any amount of
associated binary data. They are structured in named sequences of bytes supplied
along with the report.

	p, err := http.Post("http://example.com/upload", "text/plain", text)
	if err != nil {
		data := Data{
			Name: "payload",
			Blob: text,
		}

		// with data
		r.Error("client", err, data)
	}

There are some utilities for HTTP requests and HTTP responses that will rebuilt
the original wire representation.

	g, err := http.Get("http://example.com")
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		r.Error(err, ResponseHTTP(g, b))
		return
	}

	if g.StatusCode != http.StatusOk {
		err = fmt.Errorf("html error code: %s", g.Status)
		r.Error(err, ResponseHTTP(g, b))
		return
	}
*/
package report
