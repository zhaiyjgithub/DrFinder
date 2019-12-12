package main

import (
	"DrFinder/src/models"
	"DrFinder/src/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Address struct {
	Street string `json:"street"`
	City string `json:"city"`
	State string `json:"state"`
}

type Location struct {
	Locations []Address  `json:"locations"`
}

type Response struct {
	Result []struct{
		Locations []struct{
			Geo struct{
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"latLng"`
		} `json:"locations"`
	} `json:"results"`
}



func main()  {
	doctorService := service.NewDoctorService()
	geoService := service.NewGeoService()

	page := 1
	pageSize := 90

	for {
		doctors := doctorService.GetDoctorByPage(page, pageSize)
		page = page + 1

		if len(doctors) == 0 {
			fmt.Printf("finis request:  page Index: %d \n", page, )

			return
		}

		reverseAddressToGeo(doctors, geoService)
	}
}

func reverseAddressToGeo(doctors []models.Doctor, service service.GeoService) {
	url := "http://open.mapquestapi.com/geocoding/v1/batch?key=R1YKG7Kh243OBaRGeDt4MFwxHD4a47q8"

	loc := &Location{}
	for i := 0; i < len(doctors); i ++ {
		d := doctors[i]
		p := Address{
			Street: d.Address,
			City:   d.City,
			State:  d.State,
		}

		loc.Locations = append(loc.Locations, p)
	}

	m, _ := json.Marshal(loc)

	fmt.Println(string(m))

	res, err := http.Post(url, "application/json", strings.NewReader(string(m)))
	if err != nil {
		fmt.Println("-----------------")
		fmt.Printf("request error: %v \n", err)
		fmt.Println("-----------------")
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("-----------------")
		fmt.Printf("read body: %v \n", err)
		fmt.Println("-----------------")
	}
	//ummashal
	var response Response
	json.Unmarshal(body, &response)

	for i := 0; i < len(response.Result); i ++ {
		npi := doctors[i].Npi
		fmt.Printf("npi: %d {%f, %f} \n", npi, response.Result[i].Locations[0].Geo.Lat, response.Result[i].Locations[0].Geo.Lng)

		geo := &models.Geo{}
		geo.Npi = npi
		geo.Lat = response.Result[i].Locations[0].Geo.Lat
		geo.Lng = response.Result[i].Locations[0].Geo.Lng

		service.Add(geo)
	}
}
