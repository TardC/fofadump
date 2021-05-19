package fofadump

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tardc/fofadump/config"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type FofaClient struct {
	Config              *config.Config
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

func (fc *FofaClient) DoWork(fofaQuery string, size int, fields string) {
	fr := fc.FetchResult(fofaQuery, size, fields)
	fc.FetchResultCallback(fr)
}

func (fc *FofaClient) FetchResult(fofaQuery string, size int, fields string) *SearchResult {
	b64Query := base64.StdEncoding.EncodeToString([]byte(fofaQuery))
	searchApi := fmt.Sprintf("/api/v1/search/all?email=%s&key=%s&qbase64=%s&size=%d&fields=%s",
		fc.Config.Email,
		fc.Config.Key,
		b64Query,
		size,
		fields,
	)

	searchResult := &SearchResult{}
	resp, err := http.Get(fc.Config.FofaServer + searchApi)
	if err != nil {
		log.Fatalln("Request fofa server failed:", err)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("Read response body failed:", err)
	}
	//log.Println(string(result))
	err = json.Unmarshal(result, searchResult)
	if err != nil {
		log.Fatalln("Unmarshal data failed:", err)
	}
	return searchResult
}

func (fc *FofaClient) SetFetchResultCallback(fn func(*SearchResult)) {
	fc.FetchResultCallback = fn
}

func NewFofaClient(cfg *config.Config) *FofaClient {
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
