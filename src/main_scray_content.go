package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func main() {
	res, err := http.Get("https://www.doximity.com/pub/richard-xu-md")

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
	})

	doc.Find(".profile-head").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Find(".user-name-first").Text())
		fmt.Println(selection.Find(".user-name-middle").Text())
		fmt.Println(selection.Find(".user-name-last").Text())
		fmt.Println(selection.Find(".user-name-credentials").Text())
		fmt.Println(selection.Find(".user-subspecialty").Text())
		fmt.Println(selection.Find(".user-job-title").Text())
	})

	doc.Find(".address-info").Each(func(i int, selection *goquery.Selection) {
		phone := selection.Find(".office-info-telephone").Text()
		phone = strings.Replace(phone, "Phone", "", 1)
		phone = strings.Replace(phone, " ", "", -1)
		phone = strings.Replace(phone, "(", "", -1)
		phone = strings.Replace(phone, ")", "", -1)
		phone = strings.Replace(phone, "-", "", -1)
		fmt.Println(phone)

	})

	doc.Find(".summary-info").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
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
}