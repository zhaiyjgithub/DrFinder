package main

import (
	"DrFinder/src/dao"
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
	"DrFinder/src/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Education struct{
	Name string `bson:"name"`
	Title string	`bson:"title"`
}

type Certification struct{
	Name string `bson:"name"`
	ActiveDate string `bson:"activeDate"`
}

type DoctorDetail struct {
	Npi int `bson:"npi"`
	Specialty string `bson:"specialty"`
	SubSpecialty string `bson:"subSpecialty"`
	Hospitals string	`bson:"hospitals"`
	Education []Education `bson:"educations"`
	Certifications []Certification `bson:"certifications"`
	Awards []string `bson:"awards"`
}

type LocationInfo struct {
	ID string
	Npi int
	Address string
}

type UrlInfo struct {
	ID bson.ObjectId `bson:"_id"`
	Npi int `bson:"npi"`
	Url string `bson:"url"`
	FullName string `bson:"fullname"`
}

const (
	BaseUrl = "https://health.usnews.com"
)

const UsaDoctors = "doctors"
const UsaDoctorsDetail = "doctors_detail"

type DoctorInfoDao struct {
	engine *mongo.Database
}

func NewDoctorInfoDao(engine *mongo.Database) *DoctorInfoDao {
	return &DoctorInfoDao{engine: engine}
}

func (d *DoctorInfoDao) AddDoctorInfo(info *DoctorInfo) error {
	_, err := d.engine.Collection(UsaDoctors).InsertOne(context.TODO(), info)
	return err
}

func (d *DoctorInfoDao) AddDoctorInfoDetail(info *DoctorDetail) error {
	_, err := d.engine.Collection(UsaDoctorsDetail).InsertOne(context.TODO(), info)
	return err
}

func (d *DoctorInfoDao) GetDoctorByPage(page int, pageSize int) []UrlInfo {
	cur, err := d.engine.Collection(UsaDoctors).Aggregate(context.TODO(), mongo.Pipeline{
		{{"$match", bson.M{"address": ""}}},
		{{"$skip", page*pageSize}},
		{{"$limit", pageSize}},
	})

	if err != nil {
		fmt.Println("read err")
		return nil
	}

	var infos []bson.M
	if err = cur.All(context.TODO(), &infos); err != nil {
		fmt.Printf("parse info err: %s \n", err.Error())
		return nil
	}

	var docs []UrlInfo
	for i := 0; i < len(infos); i ++ {
		var doc UrlInfo
		m := infos[i]
		b, err := bson.Marshal(m)
		err = bson.Unmarshal(b, &doc)

		if err == nil {
			docs = append(docs, doc)
		}
	}

	return docs
}

func (d *DoctorInfoDao) GetDoctorInfoByPage(page int, pageSize int) []DoctorInfo {
	cur, err := d.engine.Collection(UsaDoctors).Aggregate(context.TODO(), mongo.Pipeline{
		{{"$skip", page*pageSize}},
		{{"$limit", pageSize}},
	})

	if err != nil {
		fmt.Println("read err")
		return nil
	}

	var infos []bson.M
	if err = cur.All(context.TODO(), &infos); err != nil {
		fmt.Printf("parse info err: %s \n", err.Error())
		return nil
	}

	var docs []DoctorInfo
	for i := 0; i < len(infos); i ++ {
		var doc DoctorInfo
		m := infos[i]
		b, err := bson.Marshal(m)
		err = bson.Unmarshal(b, &doc)

		if err == nil {
			docs = append(docs, doc)
		}
	}

	return docs
}

func (d *DoctorInfoDao) UpdateDoctorAddress(id string, addr string) error {
	docID := id
	objectIdjID, _ := primitive.ObjectIDFromHex(docID)
	filter := bson.M{"_id": objectIdjID}
	update := bson.M{"$set": bson.M{"address": addr}}
	upset := bool(true)
	opt := &options.UpdateOptions{Upsert: &upset}
	_, err := d.engine.Collection(UsaDoctors).UpdateOne(context.TODO(), filter, update, opt)

	return err
}



func (d *DoctorInfoDao) GetDoctorDetailByPage(page int, pageSize int) []DoctorDetail {
	cur , err := d.engine.Collection(UsaDoctorsDetail).Aggregate(context.TODO(), mongo.Pipeline{
		{{"$skip", page*pageSize}},
		{{"$limit", pageSize}},
	})

	if err != nil {
		fmt.Println("read err")
		return nil
	}

	var infos []bson.M
	if err = cur.All(context.TODO(), &infos); err != nil {
		fmt.Printf("parse info err: %s \n", err.Error())
		return nil
	}

	var docs []DoctorDetail
	for i := 0; i < len(infos); i ++ {
		var doc DoctorDetail
		m := infos[i]
		b, err := bson.Marshal(m)
		err = bson.Unmarshal(b, &doc)

		if err == nil {
			docs = append(docs, doc)
		}
	}

	return docs
}


func (d *DoctorInfoDao) GetDoctorInfoByNpi(npi int, page int, pageSize int) []UrlInfo {
	cur, err := d.engine.Collection(UsaDoctors).Aggregate(context.TODO(), mongo.Pipeline{
		{{"$match", bson.M{"npi": npi}}},
		{{"$skip", (page - 1)*pageSize}},
		{{"$limit", pageSize}},
	})

	if err != nil {
		fmt.Println("read err")
		return nil
	}

	var infos []bson.M
	if err = cur.All(context.TODO(), &infos); err != nil {
		fmt.Printf("parse info err: %s \n", err.Error())
		return nil
	}

	var docs []UrlInfo
	for i := 0; i < len(infos); i ++ {
		var doc UrlInfo
		m := infos[i]
		b, err := bson.Marshal(m)
		err = bson.Unmarshal(b, &doc)

		if err == nil {
			docs = append(docs, doc)
		}
	}


	return docs
}

type DoctorInfo struct {
	Subspecialties []string `json:"subspecialties"`
	FullName string `json:"full_name"`
	LastName string `json:"last_name"`
	Title string `json:"title"`
	AppointmentBooking struct{
		Url string `json:"url"`
		SystemName string `json:"system_name"`
	} `json:"appointment_booking"`
	Language []string `json:"language"`
	Location struct{
		City string `json:"city"`
		State string `json:"state"`
		ZipCode string `json:"zip_code"`
	} `json:"location"`
	NamePrefix string `json:"name_prefix"`
	Blurb string `json:"blurb"`
	DoctorType string `json:"doctor_type"`
	Specialty string `json:"specialty"`
	Npi int `json:"npi"`
	Phone string `json:"phone"`
	Url string `json:"url"`
	Gender string `json:"gender"`
	YearsOfExperience []int `json:"years_of_experience"`
	Address string `json:"address"`
}

type Result struct {
	DoctorSearch struct{
		Results struct{
			Doctors struct{
				Matches []DoctorInfo `json:"matches"`
			} `json:"doctors"`
		} `json:"results"`
	} `json:"doctor_search"`
}

type PfileDao struct {
	engine *gorm.DB
}

func NewPfileDao(engine *gorm.DB) *PfileDao {
	return &PfileDao{engine:engine}
}

func (d *PfileDao) AddPfile(p *models.Pfile)  {
	d.engine.Create(p)
}

func (d *PfileDao) FinderDoctorByPage(page int, pageSize int) []models.Pfile {
	var doctors []models.Pfile
	d.engine.Raw("select * from npi limit ? offset ?", pageSize, (page - 1)*pageSize).Scan(&doctors)
	return doctors
}

func main5()  {
	infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())

	//ObjectId("5ee0e2bb56ab36a7db9112c9")
	//_ = infoDao.UpdateDoctorAddress("5ee0e2bb56ab36a7db9112c9", "32st 0000")

	page := 0
	pageSize := 12

	for ;; {
		fmt.Printf("current page:  %d \n", page)
		urlInfos := infoDao.GetDoctorByPage(page, pageSize)

		if len(urlInfos) == 0 {
			fmt.Println("done....")
			fmt.Printf("Current page: %d \n", )
			return
		}

		wg := sync.WaitGroup{}
		cin := make(chan *LocationInfo)
		wg.Add(len(urlInfos))

		for i:=0;i < len(urlInfos); i ++ {
			info := urlInfos[i]

			go func(url UrlInfo) {
				location := collyLocation(info.ID.Hex() ,info.Npi, info.Url)
					cin <- location
				wg.Done()
			}(info)
		}

		go func() {
			wg.Wait()
			close(cin)
		}()

		for v := range cin {
			fmt.Printf("%v \n", v)
			if v != nil {
				infoDao.UpdateDoctorAddress(v.ID, v.Address)
			}
		}

		page = page + 1
		time.Sleep(15*time.Second)
	}
}

func main0()  {
	infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())
	//docDao := dao.NewDoctorDao(dataSource.InstanceMaster())
	hosDao := dao.NewAffiliationDao(dataSource.InstanceMaster())
	eduDao := dao.NewEducationDao(dataSource.InstanceMaster())
	awardDao := dao.NewAwardDao(dataSource.InstanceMaster())

	page := 0
	size := 1000
	for ;; {
		infos := infoDao.GetDoctorDetailByPage(page, size)

		if len(infos) == 0 {
			fmt.Println("Done")
			return
		}

		for i := 0; i < len(infos); i ++ {
			//var hospitals []models.Affiliation
			//var awards []models.Award
			//var edus []models.Education

			info := infos[i]
			npi := info.Npi
			var hospitalNams []string

			fmt.Printf("NPI: %d  \n", npi)
			hospitalNams = strings.Split(info.Hospitals, ",")
			for _, name := range hospitalNams {
				//hospitals = append(hospitals, models.Affiliation{Npi: npi, Name: name})
				_ = hosDao.Add(&models.Affiliation{Npi: npi, Name: name})
			}

			for _, edu := range info.Education {
				//edus = append(edus, models.Education{Npi: npi, Name: edu.Name, Desc: edu.Title})
				_ = eduDao.Add(&models.Education{Npi: npi, Name: edu.Name, Desc: edu.Title})
			}

			for _, ard := range info.Awards {
				//awards = append(awards, models.Award{Npi: npi, Name: ard})
				_ = awardDao.Add(&models.Award{Npi: npi, Name: ard})
			}

			//fmt.Printf("%v - %v - %v", hospitals, edus, awards)

		}

		page = page + 1
		time.Sleep(time.Millisecond*100)
	}
}

func main6()  {
	//infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())
	docDao := dao.NewDoctorDao(dataSource.InstanceMaster())

	cities := docDao.GetCity()

	for _, v := range cities {
		fmt.Printf("\"%s\", \n", v)
	}
	//page := 0
	//pageSize := 1000.
	//for ; ;  {
	//	infos := infoDao.GetDoctorInfoByPage(page, pageSize)
	//
	//	if len(infos) == 0 {
	//		fmt.Println("Done")
	//		fmt.Println(page)
	//		return
	//	}
	//
	//	for i:= 0; i < len(infos); i ++ {
	//		info := infos[i]
	//		var doc models.Doctor
	//		fname := strings.Split(info.FullName, " ")[0]
	//
	//		doc.Npi = info.Npi
	//		doc.LastName = info.LastName
	//		doc.FirstName = fname
	//		doc.FullName = info.FullName
	//		doc.NamePrefix = info.NamePrefix
	//		doc.Credential = info.Title
	//		doc.Gender = info.Gender
	//		doc.Address = info.Address
	//		doc.City = info.Location.City
	//		doc.State = info.Location.State
	//		doc.Zip = info.Location.ZipCode
	//		doc.Phone = info.Phone
	//		doc.Specialty = info.Specialty
	//		doc.SubSpecialty = strings.Join(info.Subspecialties, ", ")
	//		doc.JobTitle = info.DoctorType
	//		doc.Summary = info.Blurb
	//		doc.Lang = strings.Join(info.Language, ", ")
	//
	//		var years []string
	//		for _, y := range info.YearsOfExperience {
	//			years = append(years, fmt.Sprintf("%d", y))
	//		}
	//		doc.YearOfExperience = strings.Join(years, ", ")
	//
	//		docDao.Add(&doc)
	//	}
	//	page = page +1
	//	time.Sleep(time.Millisecond*100)
	//}
}

func main2()  {
	//infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())
	//
	//page := 4820
	//pageSize := 12
	//
	//for ;; {
	//	fmt.Printf("current page:  %d \n", page)
	//	urlInfos := infoDao.GetDoctorByPage(page, pageSize)
	//
	//	if len(urlInfos) == 0 {
	//		fmt.Println("done....")
	//		fmt.Printf("Current page: %d \n", )
	//		return
	//	}
	//
	//	wg := sync.WaitGroup{}
	//	cin := make(chan *DoctorDetail)
	//	wg.Add(len(urlInfos))
	//
	//	for i:=0;i < len(urlInfos); i ++ {
	//		info := urlInfos[i]
	//
	//		go func(url UrlInfo) {
	//			detail := colly(info.Npi, info.Url, info.FullName)
	//			if detail != nil {
	//				cin <- detail
	//			}
	//
	//			wg.Done()
	//		}(info)
	//	}
	//
	//	go func() {
	//		wg.Wait()
	//		close(cin)
	//	}()
	//
	//	for v := range cin {
	//		infoDao.AddDoctorInfoDetail(v)
	//	}
	//
	//	page = page + 1
	//	time.Sleep(15*time.Second)
	//}

}

func main1()  {
	fmt.Println("Hello world.")

	pfileDao := NewPfileDao(dataSource.InstanceMaster())
	infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())

	page := 22984
	pageSize := 20

	for ; ;  {
		fmt.Printf("current page:  %d", page)
		pfileDrs := pfileDao.FinderDoctorByPage(page, pageSize)

		if len(pfileDrs) == 0 {
			fmt.Println("finish....")
			return
		}

		wg := sync.WaitGroup{}
		wg.Add(len(pfileDrs))
		cin := make(chan *DoctorInfo)

		for i:= 0; i < len(pfileDrs); i ++ {
			go func(pfile *models.Pfile) {
				info := queryDr(pfile, page)
				if info != nil {
					cin <- info
				}

				wg.Done()
			}(&pfileDrs[i])
		}

		go func() {
			wg.Wait()

			fmt.Println("close cin ")
			close(cin)
		}()

		for vInfo := range cin {
			fmt.Println("insert info...")
			infoDao.AddDoctorInfo(vInfo)
		}

		page = page + 1
		time.Sleep(10*time.Second)
	}
}

func main()  {
	infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())
	docDao := dao.NewDoctorDao(dataSource.InstanceMaster())
	geoService := service.NewGeoService()

	page := 38
	pageSize := 1
	for ; ;  {
		docs := docDao.GetDoctorsNoAddress(page, pageSize)
		if len(docs) == 0 {
			fmt.Printf("doc page: %d \n", page)
			fmt.Println("finish doctor: ")

			return
		}

		for i:= 0; i < len(docs);i ++ {
			doc := docs[i]

			urlInfos := infoDao.GetDoctorInfoByNpi(doc.Npi, 1, 1)
			if len(urlInfos) == 0 {
				fmt.Printf("doc npi: %d not found \n", doc.Npi)
				continue
			}

			urlInfo := urlInfos[0]
			loc := collyLocation(urlInfo.ID.Hex(), urlInfo.Npi, urlInfo.Url)

			if loc == nil {
				continue
			}

			fmt.Println(loc.Address)

			var newDoc models.Doctor
			newDoc.ID = doc.ID
			newDoc.Npi = doc.Npi
			newDoc.Address = loc.Address

			if len(newDoc.Address) > 0 {
				geo := reverseTheAddressToGeo(newDoc, page)
				fmt.Println(geo)

				geoService.Add(geo)
				docDao.UpdateDoctorAddress(newDoc)

				fmt.Println("insert")
			}

		}

		page = page + 1
		time.Sleep(time.Second*1)
	}
}

func reverseTheAddressToGeo(doctor models.Doctor, page int) *models.Geo {
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
	//https://www.mapquestapi.com/geocoding/v1/address?key=KEY&inFormat=kvp&outFormat=json&location=3916 Prince St,Ste 255, Flushing, NY 11354&thumbMaps=false
	doc := doctor

	if len(doc.Address) == 0 {
		return nil
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

		return nil
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

		//service.Add(geo)
		fmt.Printf("page: %d - addr: %s - %v \n", page, doc.Address, geo)
		fmt.Println(geo)

		return geo
	}

	return nil
}

func queryDr(pfile *models.Pfile, page int) *DoctorInfo {
	baseUrl, err := url.Parse("https://health.usnews.com/health-care/doctors/search-data?")
	if err != nil {
		fmt.Println("Malformed URL: ", err.Error())
		return nil
	}

	if len(pfile.FirstName) == 0 && len(pfile.LastName) == 0 {
		return nil
	}

	// Prepare Query Parameters
	//distance=i   n-state&location=NY&page_num=1&name=Jonathan Leibovitz&gender=male
	params := url.Values{}
	params.Add("distance", "in-state")
	params.Add("location", "NY")

	fName := fmt.Sprintf("%s %s", pfile.FirstName, pfile.LastName)
	//male := "male"
	//if pfile.Gender == "M'" {
	//	male = "female"
	//}
	params.Add("name", fName)
	//params.Add("gender", male)

	// Add Query Parameters to the URL
	baseUrl.RawQuery = params.Encode() // Escape Query Parameters

	fmt.Printf("Encoded URL is %q\n", baseUrl.String())
	path := baseUrl.String()

	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("fatal, page =  %d", page)
		//log.Fatal(err)

		fmt.Println(err.Error())
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	b := []byte(body)
	var result Result
	err = json.Unmarshal(b, &result)

	matches := result.DoctorSearch.Results.Doctors.Matches
	for i := 0; i < len(matches) ; i ++  {
		info := matches[i]
		if pfile.Npi == info.Npi {
			fmt.Printf("%s \n", info.FullName)
			return &info
		}
	}

	return nil
}

func colly(npi int, apiName string, fullName string) *DoctorDetail {
	url := fmt.Sprintf("%s%s", BaseUrl, apiName)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1`)
	res, err := client.Do(req)

	if err != nil {
		fmt.Printf("fatal error, url: %s \n", url)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("request err, url: %s \n", url)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("parse doc from io reader failed, url: %s \n", url)
		return nil
	}

	var doctorInfo DoctorDetail

	const eduCategory = "Medical"
	const cerCategory = "Certifications"
	const awardCategory = "Awards"

	var medicals []Education
	var cers []Certification
	var awards []string

	doc.Find("#experience").Each(func(i int, selection *goquery.Selection) {
		selection.Find(".mb5").Each(func(i int, selection *goquery.Selection) {
			h3 := selection.Find("h3").Text()
			if strings.HasPrefix(h3, eduCategory) {
				selection.Find("div").Each(func(i int, selection *goquery.Selection) {
					firstSel := selection.Find("p").First()
					name := firstSel.Text()
					title := firstSel.Next().Text()

					medicals = append(medicals, Education{Name: name, Title: title})
				})
			}else if strings.HasPrefix(h3, cerCategory) {
				selection.Find("div").Each(func(i int, selection *goquery.Selection) {
					firstSel := selection.Find("p").First()
					name := firstSel.Text()
					date := firstSel.Next().Text()

					cers = append(cers, Certification{Name: name, ActiveDate: date})
				})
			}else if strings.HasPrefix(h3, awardCategory) {
				selection.Find("div").Each(func(i int, selection *goquery.Selection) {
					selection.Find("p").Each(func(i int, selection *goquery.Selection) {
						award := selection.Text()
						awards = append(awards, award)
					})
				})
			}
		})
	})

	doc.Find("div[class$=hxlBWL]").Each(func(i int, selection *goquery.Selection) {
		spec := selection.Find("div[class$=hkrLwy]").Text()
		subSpec := selection.Find("dd").Text()

		doctorInfo.Specialty = spec
		doctorInfo.SubSpecialty = subSpec
	})
	doc.Find("dl[class$=jMcLtN]").Each(func(i int, selection *goquery.Selection) {
		title := selection.Find("dt").Text()
		if title == "Affiliated Hospitals" {
			hospitals := selection.Find("dd").Text()
			doctorInfo.Hospitals = hospitals
		}
	})

	doctorInfo.Education = medicals
	doctorInfo.Certifications = cers
	doctorInfo.Awards = awards
	doctorInfo.Npi = npi

	fmt.Printf("\n full name: %s, npi: %d \n", fullName, npi)

	return  &doctorInfo
}

func collyLocation( objectId string, npi int, apiName string) *LocationInfo {
	url := fmt.Sprintf("%s%s", BaseUrl, apiName)
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("User-Agent", `Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1`)
	res, err := client.Do(req)

	fmt.Printf("request url: %s \n", url)
	if err != nil {
		fmt.Printf("fatal error, url: %s \n", url)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("request err, url: %s \n", url)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Printf("parse doc from io reader failed, url: %s \n", url)
		return nil
	}

	var addr string
	doc.Find("p[class$=cNtUrf]").Each(func(i int, selection *goquery.Selection) {
		addr = selection.Text()
	})

	if len(addr) == 0 {
		doc.Find("p[class$=iyWOoY]").Each(func(i int, selection *goquery.Selection) {
			addr = selection.Text()
		})
	}


	return &LocationInfo{objectId, npi, addr}
}

