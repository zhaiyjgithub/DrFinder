package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)


func main() {
	type DoxInfo struct {
		Url string
		FirstName string
		MiddleName string
		LastName string
		Credentials string
		SubSpecialty string
		JobTitle string
		OfficePhone string
		Summary string
	}

	var dox DoxInfo
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

	doc.Find(".profile-head").Find("img").Each(func(i int, selection *goquery.Selection) {
		url, _ := selection.Attr("src")
		fmt.Println(url)
		dox.Url = url
	})

	doc.Find(".profile-head").Each(func(i int, selection *goquery.Selection) {
		fmt.Println()
		fmt.Println(selection.Find(".user-name-middle").Text())
		fmt.Println(selection.Find(".user-name-last").Text())
		fmt.Println(selection.Find(".user-name-credentials").Text())
		fmt.Println(selection.Find(".user-subspecialty").Text())
		fmt.Println(selection.Find(".user-job-title").Text())
		dox.FirstName = selection.Find(".user-name-middle").Text()
	})

	doc.Find(".address-info").Each(func(i int, selection *goquery.Selection) {
		phone := selection.Find(".office-info-telephone").Text()
		phone = strings.Replace(phone, "Phone", "", 1)
		phone = strings.Replace(phone, " ", "", -1)
		phone = strings.Replace(phone, "(", "", -1)
		phone = strings.Replace(phone, ")", "", -1)
		phone = strings.Replace(phone, "-", "", -1)
		fmt.Println(phone)
		dox.OfficePhone = phone
	})

	doc.Find(".summary-info").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
		dox.Summary = selection.Text()
	})

	doc.Find(".education-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			fmt.Println(selection.Find("strong").Text())
			fmt.Println(selection.Find("span").Text())
		})
	})

	doc.Find(".certification-info").Each(func(i int, selection *goquery.Selection) {
		selection.Find("ul").Find("li").Each(func(i int, selection *goquery.Selection) {
			fmt.Println(selection.Find("strong").Text())
			fmt.Println(selection.Find("span").Text())
		})
	})

	fmt.Println(dox)
}