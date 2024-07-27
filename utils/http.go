package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HttpGet sends an HTTP GET request and returns the response body and error.
func HttpGet(urlPath string) ([]byte, error) {
	resp, err := http.Get(urlPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// HttpPost sends an HTTP POST request with form or json data and returns the response body and error.
// 'contentType' can be "application/json" or "application/x-www-form-urlencoded".
// 'data' should be either a map (for form data) or a struct (for JSON).
func HttpPost(urlPath string, contentType string, data interface{}) ([]byte, error) {
	var reqBody []byte
	var err error

	if contentType == "application/json" {
		reqBody, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	} else if contentType == "application/x-www-form-urlencoded" {
		formdata, ok := data.(map[string]string)
		if !ok {
			return nil, err
		}
		formData := url.Values{}
		for key, value := range formdata {
			formData.Set(key, value)
		}
		reqBody = []byte(formData.Encode())
	} else {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlPath, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
