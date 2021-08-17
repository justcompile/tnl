package types

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Headers map[string]string
	Status  int
	Body    []byte
}

func SerializeResponse(r *http.Response) ([]byte, error) {
	headers := make(map[string]string)

	for key := range r.Header {
		headers[key] = r.Header.Get(key)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	resp := &Response{
		Headers: headers,
		Body:    data,
		Status:  r.StatusCode,
	}

	return json.Marshal(resp)
}
