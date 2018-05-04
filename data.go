package streamr

import (
	"fmt"
	"net/http"
)

// DataService handles pushing data to streamr.
type DataService service

// Data is a Streamr Data entity.
type Data interface{}

// ProduceToStream sends data to Streamr.
func (d DataService) ProduceToStream(streamID string, data Data) (*Response, error) {
	url := fmt.Sprintf("streams/%v/data", streamID)
	req, err := d.client.NewRequest(http.MethodPost, url, data)
	if err != nil {
		return nil, fmt.Errorf("http request to streamr.com failed: %v", err)
	}
	res, err := d.client.Do(req, nil)
	if err != nil {
		return nil, fmt.Errorf("http client error: %v", err)
	}
	return res, nil
}
