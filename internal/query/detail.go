package query

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/meowrain/dockersearch/internal/httpclient"
	"github.com/meowrain/dockersearch/internal/models"
	"github.com/sirupsen/logrus"
)

const (
	imageDetailBaseURL = "https://ai-proxy.hidewiki.top/hubdocker/v2/repositories/"
)

// QueryImageDetailInfo retrieves detailed information about a specific image repository
// from a remote registry. It constructs the target URL based on the provided namespace
// and repository name, fetches the data using an HTTP GET request, and unmarshals the
// JSON response into a Repository struct. Any errors during the process are logged.
//
// Parameters:
//   - namespace: The namespace of the image repository.
//   - repository: The name of the image repository.
//
// Returns:
//   - A pointer to a Repository struct containing the image details, or nil if an error occurred.
func QueryImageDetailInfo(namespace string, repository string) *models.Repository {
	log.Println("get Image DetailInfo:", namespace, repository)
	targetURL := fmt.Sprintf("%s/%s/%s", imageDetailBaseURL, namespace, repository)
	data := httpclient.Get(targetURL)
	response := &models.Repository{}
	err := json.Unmarshal(data, response)
	if err != nil {
		logrus.Error("get Image DetailInfo error: ", err)
	}
	return response
}

func QueryImageTagsInfo(namespace string, repository string) *models.TagList {
	log.Println("get Image TagsInfo:", namespace, repository)
	targetURL := fmt.Sprintf("%s/%s/%s/tags", imageDetailBaseURL, namespace, repository)
	data := httpclient.Get(targetURL)
	response := &models.TagList{}
	err := json.Unmarshal(data, response)
	if err != nil {
		logrus.Error("get Image TagsInfo error: ", err)
	}
	return response
}
