package bingo

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	bingURL = "https://api.cognitive.microsoft.com/bing/v5.0/"
)

type (
	Bing struct {
		Search searchAPI
	}

	baseAPI struct {
		baseURL string
		client  *http.Client
		apiKey  string
	}

	statusCodeError struct {
		StatusCode int
	}
)

func NewBing(apiKey string) *Bing {
	client := &http.Client{}
	return &Bing{Search: newSearchAPI(client, apiKey)}
}

func (a baseAPI) signRequest(r *http.Request) {
	r.Header.Set("Ocp-Apim-Subscription-Key", a.apiKey)
}

func (a baseAPI) get(apiURL string) ([]byte, error) {
	req, reqErr := http.NewRequest(http.MethodGet, apiURL, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	a.signRequest(req)
	log.Printf("GET: %s\n", apiURL)
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return body, statusCodeError{StatusCode: resp.StatusCode}
	}

	return ioutil.ReadAll(resp.Body)
}

func (e statusCodeError) Error() string {
	return fmt.Sprintf("StatusCode: %d", e.StatusCode)
}
