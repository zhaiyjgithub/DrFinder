package main

import (
	"DrFinder/src/models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)


func main() {
	var doctor models.Doctor
	res, err := http.Get("https://www.doximity.com/pub/krystal-cascetta-md")

	if err!= nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code is not 200")
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//doc.Find(".profile-head").Find("img").Each(func(i int, selection *goquery.Selection) {
	//	url, _ := selection.Attr("src")
	//	fmt.Println(url)
	//	dox.Url = url
	//})

	doc.Find(".address-info").Each(func(i int, selection *goquery.Selection) {
		phone := selection.Find(".office-info-telephone").Text()
		phone = strings.Replace(phone, "Phone", "", 1)
		phone = strings.Replace(phone, " ", "", -1)
		phone = strings.Replace(phone, "(", "", -1)
		phone = strings.Replace(phone, ")", "", -1)
		phone = strings.Replace(phone, "-", "", -1)

		doctor.Phone = phone
	})

	doc.Find(".profile-head").Each(func(i int, selection *goquery.Selection) {
		doctor.SubSpecialty = selection.Find(".user-subspecialty").Text()
		doctor.JobTitle = selection.Find(".user-job-title").Text()
	})

	doc.Find(".summary-info").Each(func(i int, selection *goquery.Selection) {
		doctor.Summary = selection.Text()
	})

	//var edus []string
	doc.Find(".education-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find("span").Text())
		})
	})

	doc.Find(".certification-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find("span").Text())
		})
	})

	//trials-info
	doc.Find(".trials-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find("span").Text())
		})
	})

	//membership-info
	doc.Find(".membership-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			//fmt.Println(selection.Find("strong").Text())
			//fmt.Println(selection.Find(".br").Text())
		})
	})

	//language-info
	doc.Find(".language-info").Each(func(i int, selection *goquery.Selection) {
		//fmt.Println(selection.Find("ul").Find("li").Text())
	})
	//hospital-info
	doc.Find(".hospital-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			fmt.Println(selection.Find("strong").Text())
			fmt.Println(selection.Find(".br").Text())
		})
	})

	//fmt.Println(doctor)
}