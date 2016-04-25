package backend

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/nanopack/slurp/config"
)

type hoarder struct{}

func (self hoarder) initialize() error {
	_, err := self.rest("GET", "ping", nil)
	return err
}

func (self hoarder) readBlob(id string) (io.ReadCloser, error) {
	fmt.Println("Id:", id)
	res, err := self.rest("GET", "blobs/"+id, nil)
	return res.Body, err
}

func (self hoarder) writeBlob(id string, blob io.Reader) error {
	_, err := self.rest("POST", "blobs/"+id, blob)
	return err
}

func (self hoarder) rest(method, path string, body io.Reader) (*http.Response, error) {
	var client *http.Client
	client = http.DefaultClient
	uri := fmt.Sprintf("https://%s/%s", config.StoreAddr, path)

	if config.Insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-AUTH-TOKEN", config.StoreToken)
	res, err := client.Do(req)
	if err != nil {
		// if requesting `https://` failed, server may have been started with `-i`, try `http://`
		uri = fmt.Sprintf("http://%s/%s", config.StoreAddr, path)
		req, er := http.NewRequest(method, uri, body)
		if er != nil {
			panic(er)
		}
		req.Header.Add("X-AUTH-TOKEN", config.StoreToken)
		var err2 error
		res, err2 = client.Do(req)
		if err2 != nil {
			// return original error to client
			return nil, err
		}
	}
	if res.StatusCode == 401 {
		return nil, fmt.Errorf("401 Unauthorized. Please specify backend api token (-T 'backend-token')")
	}
	return res, nil
}
