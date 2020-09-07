package dao

import (
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
	"context"
	"encoding/json"
	"github.com/olivere/elastic/v7"
)

type DoctorElasticDao struct {
	client *elastic.Client
}

func NewDoctorElasticDao(client *elastic.Client) *DoctorElasticDao {
	return &DoctorElasticDao{client:client}
}

func (d *DoctorElasticDao) AddDoctor(doctor *models.Doctor)  {
	type DoctorES struct {
		Npi int `json:"npi"`
		FullName string `json:"full_name"`
		LastName string `json:"last_name"`
		FirstName string `json:"first_name"`
		Specialty string `json:"specialty"`
		SubSpecialty string `json:"sub_specialty"`
		Address string `json:"address"`
		City string `json:"city"`
		State string `json:"state"`
		ZipCode int `json:"zip_code"`
		Gender int `json:"gender"`
		Location string `json:"location"`
	}


}

func (d *DoctorElasticDao)QueryDoctor(doctorName string,
	specialty string,
	gender int,
	state string,
	city string,
	address string,
	zipCode int,
	page int,
	pageSize int,
	) []int {

	npiList := make([]int, 0)

	q := elastic.NewBoolQuery()

	if len(doctorName) > 0 {
		q = q.Must(elastic.NewMatchQuery("full_name", doctorName))
	}

	if gender > 0 {
		q = q.Must(elastic.NewTermQuery("gender", gender))
	}

	if len(state) > 0 {
		q = q.Must(elastic.NewTermQuery("state", state))
	}

	if len(city) > 0 {
		q = q.Must(elastic.NewTermQuery("city", city))
	}

	if len(address) > 0 {
		q = q.Must(elastic.NewMatchQuery("address", address))
	}

	if zipCode > 0 {
		q = q.Must(elastic.NewTermQuery("zip_code", zipCode))
	}

	result, err := d.client.Search().Index(dataSource.IndexDoctorName).
		Size(pageSize).
		From((page - 1)*pageSize).
		Query(q).Pretty(true).Do(context.Background())

	type Npi struct {
		Npi int `json:"npi"`
	}

	for _, hit := range result.Hits.Hits {
		var npi Npi
		err = json.Unmarshal(hit.Source, &npi)
		if err != nil {
			continue
		}

		npiList = append(npiList, npi.Npi)
	}

	return npiList
}

func (d *DoctorElasticDao) QueryNewByDoctor(lat float64, 
	lon float64,
	distance string,
	page int,
	pageSize int,
	) []int {
	npiList := make([]int, 0)

	q := elastic.NewGeoDistanceQuery("doctor.location")
	q = q.Lat(lat)
	q = q.Lon(lon)
	q = q.Distance(distance)
	
	type Npi struct {
		Npi int `json:"npi"`
	}
	
	result, err := d.client.Search().
		Index(dataSource.IndexDoctorName).
		Size(pageSize).
		From((page - 1)* pageSize).
		Query(q).Pretty(true).
		Do(context.Background())

	for _, hit := range result.Hits.Hits {
		var npi Npi
		err = json.Unmarshal(hit.Source, &npi)
		if err != nil {
			continue
		}
		
		npiList = append(npiList, npi.Npi)
	}

	return npiList
}