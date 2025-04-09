package query

import (
	"fmt"
	"testing"
)

func TestQueryImageDetailInfo(t *testing.T) {

	// Call the function
	result := QueryImageDetailInfo("bitnami", "mysql")

	t.Log("result: ", result.FullDescription)
	fmt.Printf("result: %+v\n", result)
	fmt.Printf("result.User: %s\n", result.User)
	fmt.Printf("result.Name: %s\n", result.Name)
	fmt.Printf("result.Namespace: %s\n", result.Namespace)
	fmt.Printf("result.RepositoryType: %s\n", result.RepositoryType)
	fmt.Printf("result.Status: %d\n", result.Status)
	fmt.Printf("result.StatusDescription: %s\n", result.StatusDescription)
	fmt.Printf("result.Description: %s\n", result.Description)
	fmt.Printf("result.IsPrivate: %t\n", result.IsPrivate)
	fmt.Printf("result.IsAutomated: %t\n", result.IsAutomated)
	fmt.Printf("result.StarCount: %d\n", result.StarCount)
	fmt.Printf("result.PullCount: %d\n", result.PullCount)
	fmt.Printf("result.LastUpdated: %s\n", result.LastUpdated)
	fmt.Printf("result.LastModified: %s\n", result.LastModified)
	fmt.Printf("result.DateRegistered: %s\n", result.DateRegistered)
	fmt.Printf("result.CollaboratorCount: %d\n", result.CollaboratorCount)
	if result.Affiliation != nil {
		fmt.Printf("result.Affiliation: %s\n", *result.Affiliation)
	} else {
		fmt.Printf("result.Affiliation: <nil>\n")
	}
	fmt.Printf("result.HubUser: %s\n", result.HubUser)
	fmt.Printf("result.HasStarred: %t\n", result.HasStarred)
	fmt.Printf("result.FullDescription: %s\n", result.FullDescription)
	fmt.Printf("result.Permissions: %+v\n", result.Permissions)
	fmt.Printf("result.MediaTypes: %v\n", result.MediaTypes)
	fmt.Printf("result.ContentTypes: %v\n", result.ContentTypes)
	fmt.Printf("result.Categories: %v\n", result.Categories)
	fmt.Printf("result.ImmutableTags: %t\n", result.ImmutableTags)
	fmt.Printf("result.ImmutableTagsRules: %s\n", result.ImmutableTagsRules)
	fmt.Printf("result.StorageSize: %d\n", result.StorageSize)
}
