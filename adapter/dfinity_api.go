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
}

func newAPI(url string) *API {
	return &API{
		client: &http.Client{},
		url:    url,
	}
}

func (api *API) post(pro Proposal) (string, error) {
	bytesData, _ := json.Marshal(pro)
	req, _ := http.NewRequest("POST", api.url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	resp, err := api.client.Do(req)
	if err != nil {
		fmt.Println("Can not propose request, error:",err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("resp body:",string(body))

	if resp.StatusCode != http.StatusOK {
		requestBody,_ := ioutil.ReadAll(resp.Request.Body)
		fmt.Println(resp.StatusCode,resp.Status,resp.Request.URL,requestBody)
		return string(body), errors.New("Fail To post the price To dfinity server")
	}
	return string(body), nil
}
