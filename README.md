# bingo
Golang API client for Bing Search

## Client Usage

```go
bing := bingo.NewBing(apiKey)
```

### Search

#### News

Category

```go
params := bing.NewsCategoryParams{}
resp, _ := bing.Search.News.Category(params)
```

Search

```go
params := bing.NewsSearchParams{
  Q:          "something",
  Count:      10,          // Optional
  Offset:     0,           // Optional
  Mkt:        "en-US",     // Optional
  SafeSearch: "Moderate",  // Optional
}

resp, _ := bing.Search.News.Search(params)
```

TrendingTopics

```go
resp, _ := bing.Search.News.TrendingTopics()
```
