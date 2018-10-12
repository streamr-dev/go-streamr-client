// Package streamr provides a client for Streamr API.
// Usage:
//	import streamr "github.com/streamr-dev/go-streamr-client"
//
// Construct a new Streamr Client, then use the various services on the client to
// access different parts of the Streamr API. For example:
//
//	var (
//		client *streamr.Client
//		err error
//	)
//	client, err = streamr.NewClient("streamr client authentication key")
//
//	// stream data to streamr.com
//	var response *streamr.Response
//	response, err = client.Data.ProduceToStream("stream id", mydata)
package streamr
