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

func TestResponseSerialization(t *testing.T) {
	resp := http.Response{
		StatusCode: 967,
	}

	bytes, err := SerializeResponse(&resp)
	if err != nil {
		t.Error(err)
	}

	tresp, err := DeserializeResponse(bytes)
	if err != nil {
		t.Error(err)
	}

	if tresp.StatusCode != resp.StatusCode {
		t.Errorf("Deserialized Response doesn't match")
	}
}
