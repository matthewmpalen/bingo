package bingo

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//======
// News
//======

type (
	//========
	// Shared
	//========

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

	NewsInstrumentation struct {
		pingURLBase     string `json:"pingUrlBase"`
		pageLoadPingURL string `json:"pageLoadPingUrl"`
	}

	DatePublishedTime struct {
		time.Time
	}

	NewsItem struct {
		Name          string            `json:"name"`
		URL           string            `json:"url"`
		Image         NewsImage         `json:"image"`
		Description   string            `json:"description"`
		About         []NewsAbout       `json:"about"`
		Provider      []NewsProvider    `json:"provider"`
		DatePublished DatePublishedTime `json:"datePublished"`
		Category      string            `json:"category"`
	}

	//==============
	// NewsCategory
	//==============

	NewsCategory struct {
		Type  string     `json:"_type"`
		Value []NewsItem `json:"value"`
	}

	NewsCategoryParams struct {
		Category string
	}

	//============
	// NewsSearch
	//============

	NewsSearchItem struct {
		NewsItem
		URLPingSuffix string `json:"urlPingSuffix"`
	}

	NewsSearch struct {
		Type                  string              `json:"_type"`
		Instrumentation       NewsInstrumentation `json:"instrumentation"`
		ReadLink              string              `json:"readLink"`
		TotalEstimatedMatches int                 `json:"totalEstimatedMatches"`
		Value                 []NewsSearchItem    `json:"value"`
	}

	NewsSearchParams struct {
		Q          string
		Count      int
		Offset     int
		Mkt        string
		SafeSearch string
	}

	//====================
	// NewsTrendingTopics
	//====================

	NewsTrendingTopicImage struct {
		URL      string         `json:"url"`
		Provider []NewsProvider `json:"provider"`
	}

	NewsTrendingTopic struct {
		Name                   string                 `json:"name"`
		Image                  NewsTrendingTopicImage `json:"image"`
		WebSearchURL           string                 `json:"webSearchUrl"`
		WebSearchURLPingSuffix string                 `json:"webSearchUrlPingSuffix"`
		IsBreakingNews         bool                   `json:"isBreakingNews"`
	}

	NewsTrendingTopics struct {
		Type            string              `json:"_type"`
		Instrumentation NewsInstrumentation `json:"instrumentation"`
		Value           []NewsTrendingTopic `json:"value"`
	}

	newsAPI struct {
		baseAPI
	}

	videoAPI struct {
		baseAPI
	}
)

func (d *DatePublishedTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	d.Time, err = time.Parse(time.RFC3339, s+"Z")
	return
}

func (a newsAPI) Category(params NewsCategoryParams) (NewsCategory, error) {
	var newsResp NewsCategory
	v := url.Values{}
	if params.Category != "" {
		v.Set("Category", params.Category)
	}

	apiURL, _ := url.Parse(a.baseURL)
	apiURL.RawQuery = v.Encode()

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

func (a newsAPI) Search(params NewsSearchParams) (NewsSearch, error) {
	var newsResp NewsSearch
	v := valuesFromNewsSearchParams(params)
	apiURL, _ := url.Parse(a.baseURL + "search")
	apiURL.RawQuery = v.Encode()

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

func (a newsAPI) TrendingTopics() (NewsTrendingTopics, error) {
	var newsResp NewsTrendingTopics
	apiURL, _ := url.Parse(a.baseURL + "trendingtopics")

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

func valuesFromNewsSearchParams(params NewsSearchParams) url.Values {
	v := url.Values{}
	v.Set("q", params.Q)

	if params.Count > 0 {
		v.Set("count", strconv.Itoa(params.Count))
	}

	if params.Offset >= 0 {
		v.Set("offset", strconv.Itoa(params.Offset))
	}

	if params.Mkt != "" {
		v.Set("mkt", params.Mkt)
	}

	if params.SafeSearch != "" {
		v.Set("safeSearch", params.SafeSearch)
	}

	return v
}

//=======
// Video
//=======

func NewVideoParams(q string, count, offset int, mkt string) url.Values {
	v := url.Values{}
	v.Set("q", q)
	if count > 0 {
		v.Set("count", strconv.Itoa(count))
	}

	if offset >= 0 {
		v.Set("offset", strconv.Itoa(offset))
	}

	if mkt != "" {
		v.Set("mkt", mkt)
	}

	return v
}

func newVideoAPI(client *http.Client, apiKey string) videoAPI {
	baseURL := bingURL + "videos/search"
	return videoAPI{
		baseAPI: baseAPI{
			baseURL: baseURL,
			client:  client,
			apiKey:  apiKey,
		},
	}
}

//========
// Search
//========
type searchAPI struct {
	News  newsAPI
	Video videoAPI
}

func newSearchAPI(client *http.Client, apiKey string) searchAPI {
	return searchAPI{
		News:  newNewsAPI(client, apiKey),
		Video: newVideoAPI(client, apiKey),
	}
}
