package dao

import (
	"DrFinder/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
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

func CreateIndexMappingsIfNotExisting(indexName string)  {

}

func NewElasticDao(client *elastic.Client) *PostElasticDao {
	return &PostElasticDao{client:client}
}

func (d *PostElasticDao) AddOnePost(post *models.Post) error {
	type PostType struct {
		PostId int `json:"post_id"`
		Title string `json:"title"`
		Description string `json:"description"`
	}

	postType := PostType{PostId: post.ID, Title:post.Title, Description: post.Description}
	_, err := d.client.Index().Index(IndexPostName).BodyJson(postType).Do(context.Background())

	return err
}

func (d *PostElasticDao) QueryPost(content string)  {
	q := elastic.NewMatchAllQuery()
	result, err := d.client.Search().Index(IndexPostName).Query(q).Pretty(true).Do(context.Background())

	type PostType struct {
		PostId int
		Title string
		Description string
	}

	var postTypes []PostType
	for _, hit := range result.Hits.Hits {
		var postType PostType
		err = json.Unmarshal(hit.Source, &postType)

		if err != nil {
			continue
		}

		postTypes = append(postTypes, postType)
	}

	fmt.Print(postTypes)
}
