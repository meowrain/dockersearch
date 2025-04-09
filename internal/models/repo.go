package models

import "time"

// Repository 代表 Docker Hub 仓库的完整响应
type Repository struct {
	User               string      `json:"user"`
	Name               string      `json:"name"`
	Namespace          string      `json:"namespace"`
	RepositoryType     string      `json:"repository_type"`
	Status             int         `json:"status"`
	StatusDescription  string      `json:"status_description"`
	Description        string      `json:"description"`
	IsPrivate          bool        `json:"is_private"`
	IsAutomated        bool        `json:"is_automated"`
	StarCount          int         `json:"star_count"`
	PullCount          int64       `json:"pull_count"`
	LastUpdated        time.Time   `json:"last_updated"`
	LastModified       time.Time   `json:"last_modified"`
	DateRegistered     time.Time   `json:"date_registered"`
	CollaboratorCount  int         `json:"collaborator_count"`
	Affiliation        *string     `json:"affiliation"` // 可为 null
	HubUser            string      `json:"hub_user"`
	HasStarred         bool        `json:"has_starred"`
	FullDescription    string      `json:"full_description"`
	Permissions        Permissions `json:"permissions"`
	MediaTypes         []string    `json:"media_types"`
	ContentTypes       []string    `json:"content_types"`
	Categories         []Category  `json:"categories"`
	ImmutableTags      bool        `json:"immutable_tags"`
	ImmutableTagsRules string      `json:"immutable_tags_rules"`
	StorageSize        int64       `json:"storage_size"`
}

// Permissions 代表权限对象
type Permissions struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
	Admin bool `json:"admin"`
}

// Category 代表分类对象
type Category struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}
