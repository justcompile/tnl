package types

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Request struct {
	Headers map[string]string
	Method  string
	URL     string
	Body    []byte
}

func SerializeRequest(req *http.Request) ([]byte, error) {
	headers := make(map[string]string)

	for key := range req.Header {
		headers[key] = req.Header.Get(key)
	}

	data, _ := ioutil.ReadAll(req.Body)

	sReq := &Request{
		Headers: headers,
		Method:  req.Method,
		URL:     req.RequestURI,
		Body:    data,
	}

	return json.Marshal(sReq)
}
