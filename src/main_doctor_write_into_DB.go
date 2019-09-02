package main

import (
	"DrFinder/src/models"
	"DrFinder/src/service"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

func main()  {
	csvFile, err := os.Open("./src/web/sources/physicians_al.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

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
		}else if count > 1 {
			return
		}else {
			fmt.Println(len(row))

			var doctor models.Doctor
			doctor.FirstName = row[4]
			doctor.LastName = row[3]
			doctor.MiddleName = row[5]
			doctor.Gender = row[6]
			doctor.Name = fmt.Sprintf("%s %s, %s", doctor.FirstName, doctor.LastName, doctor.MiddleName)
			doctor.Credential = row[2]
			doctor.Npi = row[25]

			doctor.MailingAddress = row[7]
			doctor.MailingAddressDetail = row[8]
			doctor.MailingCity = row[9]
			doctor.MailingState = row[10]
			doctor.MailingZip = row[11]
			doctor.MailingPhone = row[12]
			doctor.MailingFax = row[13]

			doctor.BusinessAddress = row[14]
			doctor.BusinessAddressDetail = row[15]
			doctor.BusinessCity = row[16]
			doctor.BusinessState = row[16]
			doctor.BusinessState = row[17]
			doctor.BusinessZip = row[18]
			doctor.BusinessPhone = row[19]
			doctor.BusinessFax = row[20]
			doctor.Specialty = fmt.Sprintf("%s %s", row[21], row[22])
			doctor.CreatedAt = time.Now()
			doctor.UpdatedAt = time.Now()

			service := service.NewDoctorService()

			service.Add(&doctor)

			dbyte, err := json.Marshal(&doctor)

			if err != nil {
				panic(err)
			}

			str := string(dbyte)

			fmt.Println(str)
		}

		count ++
	}
}
