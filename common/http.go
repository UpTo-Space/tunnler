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
