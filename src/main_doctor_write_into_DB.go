package main

import (
	"DrFinder/src/models"
	"DrFinder/src/service"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func main()  {
	parseCity()
}

func readDoctorCsv()  {
	csvFile, err := os.Open("./src/web/sources/physicians_al.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	service := service.NewDoctorService()

	csvReader := csv.NewReader(csvFile)

	count := 0
	for {
		row, err := csvReader.Read()

		if err == io.EOF {
			break
		}else if err != nil {
			panic(err)
		}

		if count == 0 {
			//
		}else {
			var doctor models.Doctor
			doctor.FirstName = row[4]
			doctor.LastName = row[3]
			doctor.MiddleName = row[5]
			doctor.Gender = row[6]
			doctor.FullName = fmt.Sprintf("%s %s, %s", doctor.FirstName, doctor.MiddleName, doctor.LastName)
			doctor.Credential = strings.ReplaceAll(row[2], ".", "")
			npi, err := strconv.Atoi(row[25])
			if err != nil {
				return
			}
			doctor.Npi = npi
			doctor.Address = row[14]
			doctor.AddressSuit = row[15]
			doctor.City = row[16]
			doctor.State = row[17]
			doctor.Zip = row[18]
			doctor.Phone = row[19]
			doctor.Fax = row[20]

			doctor.Specialty = row[21]
			doctor.CreatedAt = time.Now()
			doctor.UpdatedAt = time.Now()

			service.Add(&doctor)
		}

		count ++
	}

	fmt.Printf("sum count: %d", count)
}

func readSpecialEn()  {
	type void struct {

	}

	var member void
	set := make(map[string] void)


	file, err := os.Open("./src/web/sources/specialty-en.txt")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	br := bufio.NewReader(file)

	for {
		a, _, c := br.ReadLine()

		if c == io.EOF {
			break
		}
		set[string(a)] = member
	}

	var specialty []string
	for k, _ := range set {
		specialty = append(specialty, k)
	}

	sort.Slice(specialty, func(i, j int) bool {
		a := []rune(specialty[i])[0]
		b := []rune(specialty[j])[0]
		return a < b
	})
	
	for i:=0; i < len(specialty); i ++ {
		fmt.Printf("{section: 0, name: \"%s\"}, \n", specialty[i])
	}
}

func parseCity()  {
	service := service.NewDoctorService()

	page := 1
	pageSize := 500

	type void struct {
	}

	var member void

	set := make(map[string]void)

	for {
		doctors := service.GetDoctorByPage(page, pageSize)

		if len(doctors) == 0 {
			var cities []string
			for k, _ := range set {
				//fmt.Printf("\"%s\"  lent: %d\n", k, len(k))
				cities = append(cities, k)
			}
			fmt.Println("finish..")

			sort.Slice(cities, func(i, j int) bool {
				a := []rune(cities[i])[0]
				b := []rune(cities[j])[0]
				return a < b
			})

			for _, v := range cities {
				fmt.Printf("\"%s\",\n", v)
			}

			return
		}

		for i := 0; i < len(doctors); i ++ {
			set[doctors[i].City] = member
		}

		page = page + 1
	}
}