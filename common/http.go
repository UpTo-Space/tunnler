package common

import (
	"bufio"
	"bytes"
	"net/http"
)

func SerializeRequest(req *http.Request) ([]byte, error) {
	var b = &bytes.Buffer{}
	if err := req.Write(b); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func DeserializeRequest(b []byte) (*http.Request, error) {
	bytes := bytes.NewReader(b)
	r := bufio.NewReader(bytes)

	req, err := http.ReadRequest(r)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func SerializeResponse(resp *http.Response) ([]byte, error) {
	var b = &bytes.Buffer{}
	if err := resp.Write(b); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func DeserializeResponse(b []byte) (*http.Response, error) {
	bytes := bytes.NewReader(b)
	r := bufio.NewReader(bytes)

	resp, err := http.ReadResponse(r, nil)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
