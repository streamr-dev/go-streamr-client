package streamr

import (
	"net/http"
	"testing"
)

func TestProduceToStream(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/streams/xyz123/data", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
	})

	data := struct {
		Name string
		Age  int
	}{
		Name: "foobar",
		Age:  99,
	}
	_, err := client.Data.ProduceToStream("xyz123", data)
	if err != nil {
		t.Fatalf("produce to stream unexpected error: %v", err)
	}
}
