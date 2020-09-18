package dao

import (
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
)

type DoctorElasticDao struct {
	client *elastic.Client
}

func NewDoctorElasticDao(client *elastic.Client) *DoctorElasticDao {
	return &DoctorElasticDao{client:client}
}

func (d *DoctorElasticDao) AddDoctor(doctor *models.Doctor, lat float64, lon float64) error {
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
		Pin struct{
			Location struct{
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			} `json:"location"`
		} `json:"pin"`
	}

	zip, _ := strconv.Atoi(doctor.Zip)
	gender := 0
	if doctor.Gender == "M" {
		gender = 1
	}else {
		gender = 2
	}
	doctorEs := &DoctorES{
		Npi: doctor.Npi,
		FullName: doctor.FullName,
		LastName: doctor.LastName,
		FirstName: doctor.FirstName,
		Specialty: doctor.Specialty,
		SubSpecialty: doctor.SubSpecialty,
		Address: doctor.Address,
		City: doctor.City,
		State: doctor.State,
		ZipCode: zip,
		Gender: gender,
		}

	doctorEs.Pin.Location.Lat = lat
	doctorEs.Pin.Location.Lon = lon

	_, err := d.client.Index().Index(dataSource.IndexDoctorName).BodyJson(doctorEs).Do(context.Background())

	return err
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

	if len(doctorName) > 0 || len(address) > 0 { //客户端输入框的地址和医生名字同事使用
		q = q.Must(elastic.NewMultiMatchQuery(doctorName, "full_name^3", "address" )) //map type = text
	}

	if len(specialty) > 0 {
		q = q.Filter(elastic.NewMatchQuery("specialty", specialty)) // map type = text
	}

	if gender > 0 {
		q = q.Must(elastic.NewTermQuery("gender", gender)) //int
	}

	if len(state) > 0 {
		q = q.Must(elastic.NewTermQuery("state", state)) //keyword
	}

	if len(city) > 0 {
		q = q.Must(elastic.NewTermQuery("city", city))
	}

	if zipCode > 0 {
		q = q.Must(elastic.NewTermQuery("zip_code", zipCode))
	}

	q = q.QueryName("Test")
	src, err := q.Source()
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(src)
	if err != nil {
		log.Fatalf("marshaling to JSON failed: %v", err)
	}
	got := string(data)

	fmt.Println(got)

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

	q := elastic.NewGeoDistanceQuery("pin.location")
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