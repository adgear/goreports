// Copyright (c) 2014 Datacratic. All rights reserved.

package report

import "crypto/md5"

func ExampleReporter() {
	reporter := Reporter{
		Name: "test",
	}

	// start the reporting background service
	reporter.Start()

	// send data
	reporter.Log("example", "hello world")

	// for more complicated cases
	text := []byte("A day without sunshine is like, you know, night.")
	data := md5.Sum(text)

	// send more data
	reporter.Log("example", "hello world",
		Data{
			Name: "quote",
			Blob: text,
		},
		Data{
			Name: "signature",
			Blob: data[:],
		},
	)
}
