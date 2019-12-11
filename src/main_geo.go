package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main()  {
	url := "http://open.mapquestapi.com/geocoding/v1/batch?key=R1YKG7Kh243OBaRGeDt4MFwxHD4a47q8"

	type Address struct {
		Street string `json:"street"`
		City string `json:"city"`
		State string `json:"state"`
	}

	type Location struct {
		Locations []Address  `json:"locations"`
	}

	//type Result struct {
	//
	//}

	type Response struct {
		Result struct{
			Locations struct{
				Geo struct{
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"latLng"`
			} `json:"locations"`
		} `json:"results"`
	}

	p := Address{
		Street: "100 PILOT MEDICAL DR",
		City: "BIRMINGHAM",
		State: "AL",
	}

	loc := &Location{}
	loc.Locations = append(loc.Locations, p)

	m, err := json.Marshal(loc)

	fmt.Println(string(m))

	res, err := http.Post(url, "application/json", strings.NewReader(string(m)))

	fmt.Println(err)
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	//ummashal
	response := &Response{}
	err = json.Unmarshal(body, response)
	fmt.Println(string(body))

}
