package common

import (
	"net/http"
	"testing"
)

func TestRequestSerialization(t *testing.T) {
	req, err := http.NewRequest("GET", "http://google.com", nil)
	if err != nil {
		t.Error(err)
	}

	bytes, err := SerializeRequest(req)
	if err != nil {
		t.Error(err)
	}

	treq, err := DeserializeRequest(bytes)
	if err != nil {
		t.Error(err)
	}

	if treq.Method != req.Method || treq.Host != req.Host {
		t.Errorf("Deserialized Request doesn't match")
	}
}
