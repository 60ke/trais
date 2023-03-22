package tools

import (
	"io"
	"net/http"
	"strings"
	"time"
)

func Get(url string, timeout int) ([]byte, error) {
	// url := "http://106.3.133.179:46657/tri_block_info?height=104360"

	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(timeout)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Post(url string, timeout int, payload *strings.Reader) ([]byte, error) {

	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(timeout)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
