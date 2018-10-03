package streamr

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"time"
)

// StreamService handles operations related to Streams.
type StreamService struct {
	service
	subByID     map[string]*Subscription
	subByStream map[string][]*Subscription
}

// Stream models the Streamr data stream.
type Stream struct {
	// ID is the unique identifier of this Stream.
	ID string `json:"id"`

	// Name is the name of this Stream.
	Name string `json:"name"`

	// Description is describes this Stream.
	Description string `json:"description"`

	// DateCreated represents the date when this Stream was created.
	DateCreated time.Time `json:"dateCreated"`

	// LastUpdated represents the date when this Stream was updated.
	LastUpdated time.Time `json:"lastUpdated"`
}

// Subscription is a subscription to a stream.
type Subscription struct {
	// ID is an unique identifier for a subscription.
	ID string `json:"id"`

	// StreamID is a unique identifier of the stream.
	StreamID string `json:"streamId"`

	// StreamPartition
	StreamPartition string `json:"streamPartition"`

	//APIKey is the API key provided by Streamr.
	APIKey string `json:"apiKey"`

	// Callback is a function that gets called when an event occurs.
	Callback func(message string) `json:"-"`
}

// NewStream creates a new Stream.
func NewStream(id, name string) *Stream {
	return &Stream{
		ID:   id,
		Name: name,
	}
}

// NewSubscription creates a new Subscription with given arguments.
func NewSubscription(streamID, streamPartition, apiKey string, callback func(string)) (*Subscription, error) {
	id, err := generateSubscriptionID()
	if err != nil {
		return nil, err
	}
	return &Subscription{
		ID:              id,
		StreamID:        streamID,
		StreamPartition: streamPartition,
		APIKey:          apiKey,
		Callback:        callback,
	}, nil
}

// generateSubscriptionID generates an unique subscription ID.
func generateSubscriptionID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	uuid := fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid, nil
}

// GetStream returns the stream specified by streamID.
func (s *StreamService) GetStream(streamID string) (*Stream, error) {
	url := fmt.Sprintf("streams/%v", streamID)
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request to streamr.com failed: %v", err)
	}
	var stream Stream
	res, err := s.client.Do(req, &stream)
	if err != nil {
		return nil, fmt.Errorf("http client error: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		return &stream, err
	}
	return &stream, nil
}

// ListStreams lists all streams under API key.
func (s *StreamService) ListStreams(query string) ([]*Stream, error) {
	// TODO(kkn): Add query parameters
	url := "streams"
	req, err := s.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("http request to streamr.com failed: %v", err)
	}
	var streams []*Stream
	res, err := s.client.Do(req, streams)
	if err != nil {
		return nil, fmt.Errorf("http client error: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		return streams, err
	}
	return streams, nil
}

// GetStreamByName returns a stream that matches given name.
func (s *StreamService) GetStreamByName(name string) (*Stream, error) {
	return nil, nil
}

// CreateStream creates a new stream.
func (s *StreamService) CreateStream(name string) (*Stream, error) {
	// TODO(kkn): Figure out what parameters are required
	req, err := s.client.NewRequest(http.MethodPost, "streams", struct {
		Name string `json:"name"`
	}{
		Name: name,
	})
	if err != nil {
		return nil, fmt.Errorf("http request to streamr.com failed: %v", err)
	}
	var stream Stream
	res, err := s.client.Do(req, &stream)
	if err != nil {
		return nil, fmt.Errorf("http client error: %v", err)
	}
	if err := res.Body.Close(); err != nil {
		return &stream, err
	}
	return &stream, nil
}

// GetOrCreateStream gets an existing Stream or creates a new one.
func (s *StreamService) GetOrCreateStream(name string) (*Stream, error) {
	return nil, nil
}

// Subscribe subscribes to real time events of a stream.
func (s *StreamService) Subscribe(streamID string, callback func(message string)) (*Subscription, error) {
	sub, err := NewSubscription(streamID, "", "", callback)
	if err != nil {
		return nil, err
	}

	// add subscription
	s.subByID[sub.ID] = sub
	if _, exists := s.subByStream[sub.StreamID]; !exists {
		s.subByStream[sub.StreamID] = make([]*Subscription, 0)
	}
	s.subByStream[sub.StreamID] = append(s.subByStream[sub.StreamID], sub)

	return sub, nil
}
