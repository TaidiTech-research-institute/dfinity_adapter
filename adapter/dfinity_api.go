package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type API struct {
	client *http.Client
	url    string
	port   string
}

func newAPI(url string, port string) *API {
	return &API{
		client: &http.Client{},
		url:    url,
		port:   port,
	}
}

func (api *API) post(pro Proposal) (string, error) {
	bytesData, _ := json.Marshal(pro)
	url := api.url+":"+api.port+"/"
	req, _ := http.NewRequest("POST", url, bytes.NewReader(bytesData))
	resp, _ := api.client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))
	if resp.StatusCode != http.StatusOK {
		return string(body), errors.New("Fail To post the price To dfinity server")
	}

	return string(body), nil
}

