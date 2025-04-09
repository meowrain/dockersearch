package query

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/meowrain/dockersearch/internal/httpclient"
	"github.com/meowrain/dockersearch/internal/models"
)

const baseURL = "https://ai-proxy.hidewiki.top/docker/v1/search?"

func QueryImageInfo(imageName string, page int, limit int) *models.SearchResponse {
	log.Println("search Image:", imageName, "page:", page, "limit:", limit)
	queryParams := url.Values{}
	queryParams.Add("n", fmt.Sprintf("%d", limit))
	queryParams.Add("page", fmt.Sprintf("%d", page))
	queryParams.Add("type", "image")
	queryParams.Add("q", imageName)

	targetURL := baseURL + queryParams.Encode()
	log.Println("query URL: ", targetURL)
	response := httpclient.Get(targetURL)
	var searchResponse models.SearchResponse
	err := json.Unmarshal(response, &searchResponse)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON: %v", err)
	}

	return &searchResponse
}
