package main

import (
	"DrFinder/src/models"
	"DrFinder/src/service"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

func main()  {
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
