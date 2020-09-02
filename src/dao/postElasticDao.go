package dao

import (
	"DrFinder/src/models"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"time"
)

type PostElasticDao struct {
	client *elastic.Client
}

const (
	IndexPostName = "post"
	IndexPostMappings = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"title":{
				"type":"text"
			},
			"description":{
				"type":"text",
			},
			"post_id":{
				"type": int
			},
			"create_date":{
				"type":"geo_point"
			},
		}
	}
}
`
)

func NewElasticDao(client *elastic.Client) *PostElasticDao {
	return &PostElasticDao{client:client}
}

func (d *PostElasticDao) AddOnePost(post *models.Post) error {
	type PostType struct {
		PostId int `json:"post_id"`
		Title string `json:"title"`
		Description string `json:"description"`
		CreateDate string `json:"create_date"`
	}

	date := post.CreatedAt.UTC().Format(time.RFC3339)
	postType := PostType{PostId: post.ID, Title:post.Title, Description: post.Description, CreateDate: date}
	_, err := d.client.Index().Index(IndexPostName).BodyJson(postType).Do(context.Background())

	return err
}

func (d *PostElasticDao) QueryPost(content string, page int, pageSize int) []int {
	q := elastic.NewMatchAllQuery()
	result, err := d.client.Search().Index(IndexPostName).
		Size(pageSize).
		From((page - 1)*pageSize).
		Query(q).
		Pretty(true).
		Do(context.Background())

	type PostType struct {
		PostId int `json:"post_id"`
	}

	var postIds []int
	for _, hit := range result.Hits.Hits {
		var postType PostType
		err = json.Unmarshal(hit.Source, &postType)

		if err != nil {
			continue
		}

		postIds = append(postIds, postType.PostId)
	}

	return postIds
}
