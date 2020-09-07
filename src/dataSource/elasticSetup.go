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
				"type":"text"
			},
			"post_id":{
				"type": int
			},
			"create_date":{
				"type":"date"
			}
		}
	}
}
`
	IndexDoctorName = "doctor"
	IndexDoctorMappings = `
{
	"settings":{
		"number_of_shards":1,
		"number_of_replicas":0
	},
	"mappings":{
		"properties":{
			"npi":{
				"type":"integer"
			},
			"full_name":{
				"type":"text"
			},
			"last_name":{
				"type":"text"
			},
			"first_name":{
				"type":"text"
			},
			"specialty":{
				"type":"text"
			},
			"sub_specialty":{
				"type":"text"
			},	
			"address":{
				"type":"text"
			},
			"city":{
				"type":"keyword"
			},
			"state":{
				"type":"keyword"
			},
			"zip_code":{
				"type":"integer"
			},
			"gender":{
				"type":"integer"
			},
			"location": {
            	"type": "geo_point"
          	},
			"create_date":{
				"type":"date"
			}
		}
	}
}
`
)

func ElasticSetup()  {
	elasticIndexInfos := [...]elasticIndexInfo{
		{Name: IndexPostName, Mapping: IndexPostMappings},
		{Name: IndexDoctorName, Mapping: IndexDoctorMappings},
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
