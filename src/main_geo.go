package main

import (
	"DrFinder/src/models"
	"DrFinder/src/service"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
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

	page := 66489
	pageSize := 1

	for {
		docs := doctorService.GetDoctorByPage(page, pageSize)

		if len(docs) == 0 {
			fmt.Printf("finis request:  page Index: %d \n", page, )
			return
		}
		doc := docs[0]
		reverseAddressToGeo(doc, geoService, page)
		time.Sleep(time.Millisecond*200)

		page = page + 1
	}
}

func reverseAddressToGeo(doctor models.Doctor, service service.GeoService, page int) {
//https://www.mapquestapi.com/geocoding/v1/address?key=KEY&inFormat=kvp&outFormat=json&location=1+dakota+dr%2C+Lake+Success%2C+NY%2C+11042&thumbMaps=false
    doc := doctor

    if len(doc.Address) == 0 {
		return
	}

	baseUrl, err := url.Parse("https://www.mapquestapi.com/geocoding/v1/address?")

	params := url.Values{}
	params.Add("key", "lYrP4vF3Uk5zgTiGGuEzQGwGIVDGuy24")
	params.Add("inFormat", "kvp")
	params.Add("outFormat", "json")
	params.Add("location",doc.Address)

	//params.Add("json", locstr)
	//params.Add("postalCode", "11042")
	// Add Query Parameters to the URL
	baseUrl.RawQuery = params.Encode() // Escape Query Parameters

	fmt.Printf("Encoded URL is %q\n", baseUrl.String())
	path := baseUrl.String()

	res, err := http.Get(path)
	if err != nil {
		fmt.Printf("page index: %d \n", page)
		log.Fatal(err.Error())
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

	if len(response.Result) > 0 {
		npi := doc.Npi

		geo := &models.Geo{}
		geo.Npi = npi
		geo.Lat = response.Result[0].Locations[0].Geo.Lat
		geo.Lng = response.Result[0].Locations[0].Geo.Lng

		service.Add(geo)
		fmt.Printf("page: %d - addr: %s - %v \n", page, doc.Address, geo)
		fmt.Println(geo)
	}
}
