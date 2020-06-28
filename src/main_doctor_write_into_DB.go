package main

import (
	"DrFinder/src/dataSource"
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
	readNPICsv()
}

//type PfileDao struct {
//	engine *gorm.DB
//}
//
//func NewPfileDao(engine *gorm.DB) *PfileDao {
//	return &PfileDao{engine:engine}
//}
//
//func (d *PfileDao) AddPfile(p *models.Pfile)  {
//	d.engine.Create(p)
//}

func getSpecialty()  {
	spMap := make(map[string]string)
	drService := service.NewDoctorService()

	sps := drService.GetSpecialty()
	for _, sp := range sps {
		if len(spMap[sp]) == 0{
			spMap[sp] = sp
		}
	}

	for k, _ := range spMap {
		fmt.Println(k)
	}
}

func readNPICsv()  {
	dao := NewPfileDao(dataSource.InstanceMaster())

	csvFile, err := os.Open("./src/web/sources/npidata_pfile_20050523-20200510.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	csvReader := csv.NewReader(csvFile)

	count := 0
	fmt.Println(time.Now())
	for {
		row, err := csvReader.Read()

		if err == io.EOF {
			fmt.Println("EOF time: ")
			fmt.Println(time.Now())
			break
		} else if err != nil {
			fmt.Println("pannic error: ")
			fmt.Println(time.Now())
			panic(err)
		}

		state := row[23]
		if state == "NY" {
			count = count + 1
			npi := row[0]
			firstName := row[6]
			lastName := row[5]
			midName := row[7]

			firstAddr := row[20]
			secondAddr := row[21]
			city := row[22]
			postalCode := row[24]
			phone := row[26]
			fax := row[27]
			gender := row[42]

			var pfile models.Pfile
			pfile.Npi, _ = strconv.Atoi(npi)
			pfile.FirstName = firstName
			pfile.LastName = lastName
			pfile.MidName = midName

			pfile.FirstAddress = firstAddr
			pfile.SecondAddress = secondAddr
			pfile.City = city
			pfile.State = state
			pfile.PostalCode = postalCode
			pfile.Phone = phone
			pfile.Fax = fax
			pfile.Gender = gender

			if count > 20000 {
				// inser many
			}
			dao.AddPfile(&pfile)
		}
	}
}

func readDoctorCsv()  {
	csvFile, err := os.Open("./src/web/sources/npidata_pfile_20050523-20200510.csv")
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

			//sort.Slice(cities, func(i, j int) bool {
			//	a := []rune(cities[i])[0]
			//	b := []rune(cities[j])[0]
			//	return a < b
			//})

			sort.Sort(sort.StringSlice(cities))

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

func parseState()  {
	fs, err := os.Open("./src/web/sources/specialty-en.txt")
	if err != nil {
		panic(err)
	}

	defer fs.Close()

	br := bufio.NewReader(fs)

	for {
		a, _, c := br.ReadLine()

		if c == io.EOF {
			break
		}

		fullName := string(a)
		code := string([]rune(string(fullName))[0:2])
		name := string([]rune(string(fullName))[3: len(fullName) - 1])

		fmt.Printf("{code: \"%s\", name: \"%s\"},\n", code, name)
	}
}