package main

import (
	"DrFinder/src/dataSource"
	"DrFinder/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jinzhu/gorm"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
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


type UrlInfo struct {
	Npi int `bson:"npi"`
	Url string `bson:"url"`
	FullName string `bson:"fullname"`
}

const (
	BaseUrl = "https://health.usnews.com"
)

const UsaDoctors = "usadoctors"
const UsaDoctorsDetail = "usadoctorsdetail"

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
		{{"$skip", page*pageSize}},
		{{"$limit", pageSize}},
	})

	if err != nil {
		fmt.Println("read err")
		return nil
	}

	var info []bson.M
	if err = cur.All(context.TODO(), &info); err != nil {
		fmt.Printf("parse info err: %s \n", err.Error())
		return nil
	}

	

	var urlInfos []UrlInfo
	for i := 0; i < len(info); i ++ {
		var urlInfo UrlInfo
		m := info[i]
		b, err := bson.Marshal(m)
		err = bson.Unmarshal(b, &urlInfo)

		if err == nil {
			urlInfos = append(urlInfos, urlInfo)
		}
	}

	return urlInfos
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

func main()  {
	infoDao := NewDoctorInfoDao(dataSource.InstanceMongoDB())

	page := 2016 //198~
	pageSize := 10

	for ;; {
		fmt.Printf("current page:  %d", page)
		urlInfos := infoDao.GetDoctorByPage(page, pageSize)

		if len(urlInfos) == 0 {
			fmt.Println("done....")
			fmt.Printf("Current page: %d \n", )
			return
		}

		wg := sync.WaitGroup{}
		cin := make(chan *DoctorDetail)
		wg.Add(len(urlInfos))

		for i:=0;i < len(urlInfos); i ++ {
			info := urlInfos[i]

			go func(url UrlInfo) {
				detail := colly(info.Npi, info.Url, info.FullName)
				if detail != nil {
					cin <- detail
				}

				wg.Done()
			}(info)
		}

		go func() {
			wg.Wait()
			close(cin)
		}()

		for v := range cin {
			infoDao.AddDoctorInfoDetail(v)
		}

		page = page + 1
		time.Sleep(20*time.Second)
	}

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

