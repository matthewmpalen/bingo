package bingo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"
)

type (
	Thumbnail struct {
		ContentURL string `json:"contentUrl"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
	}

	NewsImage struct {
		Thumbnail
	}

	NewsAbout struct {
		ReadLink string `json:"readLink"`
		Name     string `json:"name"`
	}

	NewsProvider struct {
		ProviderType string `json:"_type"`
		Name         string `json:"name"`
	}

	NewsItem struct {
		Name          string         `json:"name"`
		URL           string         `json:"url"`
		Image         NewsImage      `json:"image"`
		Description   string         `json:"description"`
		About         []NewsAbout    `json:"about"`
		Provider      []NewsProvider `json:"provider"`
		DatePublished time.Time      `json:"datePublished"`
		Category      string         `json:"category"`
	}

	NewsResponse struct {
		Type  string     `json:"_type"`
		Value []NewsItem `json:"value"`
	}

	newsAPI struct {
		baseAPI
	}
)

func newNewsAPI(client *http.Client, apiKey string) newsAPI {
	baseURL := bingURL + "news/"
	return newsAPI{
		baseAPI: baseAPI{
			baseURL: baseURL,
			client:  client,
			apiKey:  apiKey,
		},
	}
}

func (a newsAPI) Get(params url.Values) (NewsResponse, error) {
	var newsResp NewsResponse
	var apiURL *url.URL
	apiURL, _ = url.Parse(a.baseURL)
	apiURL.RawQuery = params.Encode()

	body, getErr := a.get(apiURL.String())
	if getErr != nil {
		switch getErr.(type) {
		case statusCodeError:
			log.Println(string(body))
		default:
			return newsResp, getErr
		}
	}

	jsonErr := json.Unmarshal(body, &newsResp)
	return newsResp, jsonErr
}
