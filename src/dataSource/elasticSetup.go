package dataSource

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
)

type elasticIndexInfo struct {
	Name string
	Mapping string
}

const (
	indexPostName = "post"
	indexPostMappings = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"post_id": {"type": "integer"},
			"title": {"type": "text"},
			"description": {"type": "text"},
			"create_date": {"type": "date"}
		}
	}
}
`
)

func ElasticSetup()  {
	elasticIndexInfos := [...]elasticIndexInfo{
		{Name: indexPostName, Mapping: indexPostMappings},
	}
	client := InstanceElasticSearchClient()

	for _, info := range elasticIndexInfos {
		createIndexMappingsIfNotExisting(client, info.Name, info.Mapping)
	}
}

func createIndexMappingsIfNotExisting(client *elastic.Client, name string, mapping string)  {
	exist, err := client.IndexExists(name).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if exist == false {
		index, err := client.CreateIndex(name).Body(mapping).Do(context.TODO())
		if err != nil {
			log.Fatalf("Create Index(%s) fatal.", name)
		}

		if index == nil {
			log.Fatalf("Create Index(%s) fatal, should not nil.", name)
		}
	}
}
