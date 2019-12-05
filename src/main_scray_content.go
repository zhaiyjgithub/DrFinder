package main

import (
	"DrFinder/src/models"
	"DrFinder/src/service"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type DoxInfo struct {
	doctor *models.Doctor
	affiliations []*models.Affiliation
	awards []*models.Award
	certifications []*models.Certification
	clinicals []*models.Clinical
	educations []*models.Education
	memberships []*models.Membership
	langs []*models.Lang
}


func main() {
	doctorService := service.NewDoctorService()
	affiliationService := service.NewAffiliationService()
	awardService := service.NewAwardService()
	cerService := service.NewCertificationService()
	clnService := service.NewClinicalService()
	eduService := service.NewEducationService()
	langService := service.NewLangService()
	memService := service.NewMembershipService()


	page := 1
	pageSize := 5

	start := time.Now().Unix()
	for {
		doctors := doctorService.GetDoctorByPage(page, pageSize)

		if len(doctors) == 0 {
			end := time.Now().Unix()
			fmt.Printf("finish, duration: %d", end - start)
			return
		}

		page = page + 1

		wg := sync.WaitGroup{}
		cin := make(chan *DoxInfo)
		count := len(doctors)

		wg.Add(count)
		for i := 0; i < count; i ++ {
			doctor := doctors[i]
			go func(i int, doctor models.Doctor) {
				fmt.Println(doctor.FirstName, i)
				dox := fetchDoctor(&doctor)
				cin <- dox
				wg.Done()
			}(i, doctor)
		}

		go func() {
			wg.Wait()
			close(cin)
		}()

		for dox := range cin {
			if dox != nil {
				// update doctor info
				info := models.Doctor{
					ID: dox.doctor.ID,
					JobTitle: dox.doctor.JobTitle,
					SubSpecialty: dox.doctor.SubSpecialty,
					Summary: dox.doctor.Summary,
				}

				doctorService.UpdateDoctorInfo(&info)

				for i := 0; i < len(dox.affiliations); i ++ {
					aff := dox.affiliations[i]
					err := affiliationService.Add(aff)

					if err != nil {
						fmt.Println(err)
					}
				}

				for i := 0; i < len(dox.awards); i ++ {
					award := dox.awards[i]

					err := awardService.Add(award)
					if err != nil {
						fmt.Println(err)
					}
				}

				for i := 0; i < len(dox.certifications); i ++ {
					cer := dox.certifications[i]

					err := cerService.Add(cer)
					if err != nil {
						fmt.Println(err)
					}
				}

				for i := 0; i < len(dox.clinicals); i ++ {
					cln := dox.clinicals[i]

					err := clnService.Add(cln)
					if err != nil {
						fmt.Println(err)
					}
				}

				for i := 0; i < len(dox.educations); i ++ {
					edu := dox.educations[i]

					err := eduService.Add(edu)
					if err != nil {
						fmt.Println(err)
					}
				}

				for i := 0; i < len(dox.langs); i ++ {
					lang := dox.langs[i]

					err := langService.Add(lang)
					if err != nil {
						fmt.Println(err)
					}
				}

				for i := 0; i < len(dox.memberships); i ++ {
					mem := dox.memberships[i]

					err := memService.Add(mem)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		}

		time.Sleep(20*time.Second)
	}
}

func fetchDoctor(doctor *models.Doctor) *DoxInfo {
	var dox = DoxInfo{}
	dox.doctor = doctor

	url := fmt.Sprintf("https://www.doximity.com/pub/%s-%s-%s", doctor.FirstName, doctor.LastName, doctor.Credential)
	res, err := http.Get(url)

	if err!= nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Printf("name: %s-%s-%s npi: %d,  200", doctor.FirstName, doctor.LastName, doctor.Credential, doctor.Npi)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("first name: %s, npi: %d,  body", doctor.FirstName, doctor.Npi)
		return nil
	}

	//doc.Find(".profile-head").Find("img").Each(func(i int, selection *goquery.Selection) {
	//	url, _ := selection.Attr("src")
	//	fmt.Println(url)
	//	dox.Url = url
	//})

	isEqual := true
	doc.Find(".address-info").Each(func(i int, selection *goquery.Selection) {
		phone := selection.Find(".office-info-telephone").Text()
		phone = strings.Replace(phone, "Phone", "", 1)
		phone = strings.Replace(phone, " ", "", -1)
		phone = strings.Replace(phone, "(", "", -1)
		phone = strings.Replace(phone, ")", "", -1)
		phone = strings.Replace(phone, "-", "", -1)

		if dox.doctor.Phone != phone {
			fmt.Printf("first name: %s, npi: %d,  phone", doctor.FirstName, doctor.Npi)
			isEqual = false
		}else {
			dox.doctor.Phone = phone
		}
	})

	doc.Find(".profile-head").Each(func(i int, selection *goquery.Selection) {
		dox.doctor.SubSpecialty = selection.Find(".user-subspecialty").Text()
		dox.doctor.JobTitle = selection.Find(".user-job-title").Text()
		dox.doctor.JobTitle = strings.ReplaceAll(dox.doctor.JobTitle, "\n", "")
	})

	doc.Find(".summary-info").Each(func(i int, selection *goquery.Selection) {
		dox.doctor.Summary = selection.Text()
	})

	//var edus []string
	doc.Find(".education-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			var edu models.Education
			edu.Name = selection.Find("strong").Text()
			edu.Desc = selection.Find("span").Text()
			edu.Npi = doctor.Npi
			dox.educations = append(dox.educations, &edu)
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find("span").Text())
		})
	})

	doc.Find(".certification-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			var cer models.Certification
			cer.Name = selection.Find("strong").Text()
			cer.Desc = selection.Find("span").Text()
			cer.Npi = doctor.Npi
			dox.certifications = append(dox.certifications, &cer)
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find("span").Text())
		})
	})

	//trials-info
	doc.Find(".trials-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			var cln  models.Clinical
			cln.Name = selection.Find("strong").Text()
			cln.Desc = selection.Find("span").Text()
			cln.Npi = doctor.Npi
			dox.clinicals = append(dox.clinicals, &cln)
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find("span").Text())
		})
	})

	//membership-info
	doc.Find(".membership-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			var mem models.Membership
			mem.Name = selection.Find("strong").Text()
			mem.Desc = selection.Find(".br").Text()
			mem.Npi = doctor.Npi
			dox.memberships = append(dox.memberships, &mem)
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find(".br").Text())
		})
	})

	//language-info
	doc.Find(".language-info").Each(func(i int, selection *goquery.Selection) {
		//fmt.Println(selection.Find("ul").Find("li").Text())
		var lang models.Lang
		lang.Lang = selection.Find("ul").Find("li").Text()
		lang.Npi = doctor.Npi
		dox.langs = append(dox.langs, &lang)
	})
	//hospital-info
	doc.Find(".hospital-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find(".br").Text())
			var aff models.Affiliation
			aff.Name = selection.Find("strong").Text()
			aff.Desc = selection.Find(".br").Text()
			aff.Npi = doctor.Npi
			dox.affiliations = append(dox.affiliations, &aff)
		})
	})

	//award-info
	doc.Find(".award-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find(".br").Text())
			var award models.Award
			award.Name = selection.Find("strong").Text()
			award.Desc = selection.Find(".br").Text()
			award.Npi = doctor.Npi
			dox.awards = append(dox.awards, &award)
		})
	})

	if !isEqual {
		return nil
	}

	return &dox
}