package fofadump

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type FofaClient struct {
	Config              *Config
	FetchResultCallback FetchResultCallback
}

type FetchResultCallback func(*SearchResult)

type SearchResult struct {
	Error   bool          `json:"error,omitempty"`
	Mode    string        `json:"mode,omitempty"`
	Page    int           `json:"page,omitempty"`
	Query   string        `json:"query,omitempty"`
	Results []interface{} `json:"results,omitempty"`
	Size    int           `json:"size,omitempty"`
}

func (fc *FofaClient) DoWork(fofaQuery string, page, size int, fields string, full bool) {
	fr, err := fc.FetchResult(fofaQuery, page, size, fields, full)
	if err == nil && fc.FetchResultCallback != nil {
		fc.FetchResultCallback(fr)
	}
}

func (fc *FofaClient) FetchResult(fofaQuery string, page, size int, fields string, full bool) (*SearchResult, error) {
	b64Query := base64.StdEncoding.EncodeToString([]byte(fofaQuery))
	searchApi := fmt.Sprintf("/api/v1/search/all?email=%s&key=%s&qbase64=%s&page=%d&size=%d&fields=%s&full=%t",
		fc.Config.Email,
		fc.Config.Key,
		b64Query,
		page,
		size,
		fields,
		full,
	)

	searchResult := &SearchResult{}
	resp, err := http.Get(fc.Config.FofaServer + searchApi)
	if err != nil {
		log.Println("Request fofa server failed:", err)
		return nil, err
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response body failed:", err)
		return nil, err
	}
	//log.Println(string(result))
	err = json.Unmarshal(result, searchResult)
	if err != nil {
		log.Println("Unmarshal data failed:", err)
		return nil, err
	}
	return searchResult, nil
}

func (fc *FofaClient) SetFetchResultCallback(fn func(*SearchResult)) {
	fc.FetchResultCallback = fn
}

func NewFofaClient(cfg *Config) *FofaClient {
	fc := &FofaClient{}
	fc.Config = cfg

	fc.SetFetchResultCallback(func(results *SearchResult) {
		log.Printf("%d/%d\n", len(results.Results), results.Size)
		for _, result := range results.Results {
			switch result.(type) {
			case string:
				fmt.Printf("%s\n", result)
			case []interface{}:
				var tmp []string
				for _, r := range result.([]interface{}) {
					tmp = append(tmp, r.(string))
				}
				fmt.Printf("%s\n", strings.Join(tmp, ","))
			}
		}
	})

	return fc
}
