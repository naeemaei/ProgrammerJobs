// find_in_page
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	baseURL string    = "https://jobinja.ir/"
	jobPage string    = "jobs?"
	params  [4]string = [4]string{
		"filters[job_categories][0]=وب،‌ برنامه‌نویسی و نرم‌افزار",
		"filters[locations][0]=تهران",
		"sort_by=relevance_desc",
		"page=",
	}
)

func main() {

	TestConnection()
	var document = makeHTTPRequest(baseURL+jobPage+strings.Join(params[:], "&"), 1)
	// Get page count
	var allLIs = document.Find(".paginator").Find("ul li")
	var lastPage = document.Find(".paginator").Find("ul li").Eq(allLIs.Length() - 2)
	pageCount, _ := toLatinDigits(lastPage.Find("a").Text())
	var nextPage int64 = 2
	//pageCount = 2
	for nextPage <= pageCount {

		// Find and print image URLs
		document.Find(".o-listView__itemInfo").Each(func(index int, element *goquery.Selection) {

			fmt.Println(baseURL+jobPage+strings.Join(params[:], "&"),nextPage-1)
			// Get company name
			var company = element.Find(".c-icon--construction").Parent().Find("span").Text()
			fmt.Println(company)

			// Get city name
			var place = element.Find(".c-icon--place").Parent().Find("span").Text()
			fmt.Println(place)

			// Get job title name
			var jobTitle = element.Find(".c-jobListView__titleLink").Text()
			var jobLink, _ = element.Find(".c-jobListView__titleLink").Attr("href")
			// fmt.Println(jobLink)

			// Go to detail page
			// Save in db
			condb := GetConnection()

			var jobId, _ =  CreateMasterRecord(condb, jobTitle, company, place)

			var newDocument = makeHTTPRequest(jobLink, 0)
			newDocument.Find(".c-infoBox__item").Each(func(index int, element *goquery.Selection) {

				 go CreateDetailRecord(condb, int(jobId), element.Find("h4").Text(), element.Find("span").Text())
			})

			go CreateDetailRecord(condb, int(jobId), "Description", newDocument.Find(".s-jobDesc ").Text())

		})
		nextPage++
		document = makeHTTPRequest(baseURL+jobPage+strings.Join(params[:], "&"), nextPage)
	}

	fmt.Println(pageCount)

}

func makeHTTPRequest(pageAddress string, pageNumber int64) *goquery.Document {
	// Make HTTP request

	if pageNumber > 0 {
		pageAddress += strconv.FormatInt(pageNumber, 10)
	}
	req, _ := http.NewRequest("GET", pageAddress, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("cookie", "__cfduid=d7bc9f34afa3e43dcb38b6d3606b6d97a1572955671; _ga=GA1.2.1900100706.1572956281; remember_82e5d2c56bdd0811318f0cf078b78bfc=eyJpdiI6ImxaYjVTV25ld3FTa01OUmVDbXBUU3c9PSIsInZhbHVlIjoiUEpWVDI3ZHJhcUw4dENsaGRXSm5GVlpcL0FIc2QwYUZjNnFOTERxZ01nVXZHd3h3QVJ1RVB2XC84YUlXY25KNEFUZVg0STR6UkFMb1RFc0R2ZVl2bFhYblc3M250dmRWdHpla3dcL3B6aEJxaTA9IiwibWFjIjoiOTdkMmJlZTc3MTJlYjYzMjVjYWY0ZDI4Yzk2Yjc5ZTRjOWRkZGU3OWY1YmE5MmU1ZDRhY2YyOTkwNGExYTc2MSJ9; device_id=eyJpdiI6IjVYVXVtb2NudXRvd3JpeUo0ZDN2V0E9PSIsInZhbHVlIjoiMkduY09cLzlCZEdiV2pIQVpWdHU2OFE9PSIsIm1hYyI6ImY5ZTlmNWM5ZTE5OWQ4YThmYWFiYzk5MWI0MTMwZDQ0YzRkZDQ2NjViZWY0ZTFmMDBmMmI3NmE0YTFjNDM0NjEifQ%3D%3D; user_mode=eyJpdiI6ImZ5eVBubmFTaFh2TXV5UnVRNmRDUGc9PSIsInZhbHVlIjoiUWtzVjYrc0ZwTWhyWUpYR1R5SUJGUT09IiwibWFjIjoiYWI0Y2E1NzEzNDVlNjVjZjA5NTNjOWMzNWM3ZWM3NWQ2Y2UyMGI2NWM2YzQ1ZWFjOTUwZDFmMWQxYWMwMDM3MCJ9; logglytrackingsession=e39ce616-7c0a-44a5-82da-37207443fe15; _gid=GA1.2.547398638.1582531405; _gat_gtag_UA_129125700_1=1; XSRF-TOKEN=eyJpdiI6ImkwQ1czOVBOTU1EZUJWNUIySFdOQVE9PSIsInZhbHVlIjoiK2tnY0E0UkhSU3ZvKzBiT1lmZ3ZXSEJQTW5GNjkrQzRLVDFvUks0Zit0aERFcForU2E1cjNOcSszbERjOWxSSmlERXZzWk9JejZnNmVJcU5mVEs0akE9PSIsIm1hYyI6ImY3NjkzOTUxMmZmM2ZjMWUyZDRkMjI2Njk5MzEyNDk2NWZhM2FlNzJlMjMxMWQ2OGFiZWQxOWYzMzYxNTljMzUifQ%3D%3D; JSESSID=eyJpdiI6InFOQVNJbDl0RVoyZ1dDdEZDNll1MHc9PSIsInZhbHVlIjoiSnBoSngzYmY1VjdxbldkS1lCNjBxbUpcL0lpZ3FtNTNVaWtlU0JPYjRqTVdtMWQwN1JFU1NrMUtRcHR1Yk5NYlEyNHR1NnZNY2pFNFliN0lKWG5pQ2JBPT0iLCJtYWMiOiIxY2MxODYzZjc0MWJkNzg0MDUzZTljZTYyOWM0NzViYWFiNTI2MDYwOWFiOGQwZjFkMWViZTM5OWMyM2E2NmY1In0%3D")
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, _ := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	return document
}

func toLatinDigits(persianNumber string) (i int64, err error) {
	var LatinDigits = strings.ReplaceAll(persianNumber, "۰", "0")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۷", "7")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۱", "1")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۲", "2")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۳", "3")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۴", "4")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۵", "5")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۶", "6")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۸", "8")
	LatinDigits = strings.ReplaceAll(LatinDigits, "۹", "9")

	return strconv.ParseInt(LatinDigits, 10, 64)

}
