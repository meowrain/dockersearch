package models

import "time"

// TagList 代表 Docker Hub 标签列表的完整响应
type TagList struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`     // 可为 null
	Previous *string `json:"previous"` // 可为 null
	Results  []Tag   `json:"results"`
}

// Tag 代表单个标签的元数据
type Tag struct {
	Creator             int        `json:"creator"`
	ID                  int        `json:"id"`
	Images              []Image    `json:"images"` // 可为空数组
	LastUpdated         time.Time  `json:"last_updated"`
	LastUpdater         int        `json:"last_updater"`
	LastUpdaterUsername string     `json:"last_updater_username"`
	Name                string     `json:"name"`
	Repository          int        `json:"repository"`
	FullSize            int64      `json:"full_size"`
	V2                  bool       `json:"v2"`
	TagStatus           string     `json:"tag_status"`
	TagLastPulled       *time.Time `json:"tag_last_pulled"` // 可为 null
	TagLastPushed       *time.Time `json:"tag_last_pushed"` // 可为 null
	MediaType           string     `json:"media_type"`
	ContentType         string     `json:"content_type"`
	Digest              string     `json:"digest"`
}

// Image 代表标签中的单个镜像（按架构区分）
type Image struct {
	Architecture string     `json:"architecture"`
	Features     string     `json:"features"`
	Variant      *string    `json:"variant"` // 可为 null
	Digest       string     `json:"digest"`
	OS           string     `json:"os"`
	OSFeatures   string     `json:"os_features"`
	OSVersion    *string    `json:"os_version"` // 可为 null
	Size         int64      `json:"size"`
	Status       string     `json:"status"`
	LastPulled   *time.Time `json:"last_pulled"` // 可为 null
	LastPushed   *time.Time `json:"last_pushed"` // 可为 null
}
