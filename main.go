package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

type extractedJob struct {
	id       string
	title    string
	company  string
	location string
	salary   string
}

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	// getPage((totalPages + 1) - (totalPages))

	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	// for _, value := range jobs {
	// 	fmt.Println("title: ", value.title)
	// 	fmt.Println("company: ", value.company)
	// 	fmt.Println("location: ", value.location)
	// 	fmt.Println("salary: ", value.salary)
	// 	fmt.Println("========================")
	// }
	writeJobs(jobs)
	fmt.Println("Done, extracted ", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkError(err)
	w := csv.NewWriter(file)
	defer w.Flush()
	headers := []string{"id", "title", "company", "location", "salary"}
	wErr := w.Write(headers)
	checkError(wErr)
	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.company, job.location, job.salary}
		jwErr := w.Write(jobSlice)
		// a, _ := utf8.DecodeRuneInString(jobSlice)
		// jwErr := w.Write(a)
		checkError(jwErr)
	}
}

func getPage(page int) (jobs []extractedJob) {
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting ", pageURL)
	res, err := http.Get(pageURL)
	checkError(err)
	checkCode(res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, s *goquery.Selection) {
		jobs = append(jobs, extractJob(s))
	})
	return
}

func extractJob(card *goquery.Selection) extractedJob {
	// title := card.Find(".title").Find("a").Text()
	id, _ := card.Attr("data-jk")
	title := cleanString(card.Find(".title>a").Text())
	company := cleanString(card.Find(".sjcl").Find(".company").Text())
	location := cleanString(card.Find(".sjcl").Find(".location").Text())
	// salary := cleanString(card.Find(".salarySnippet").Find(".salaryText").Text())
	salary := cleanString(card.Find(".salaryText").Text())

	// // fmt.Println("id: ", id)
	// fmt.Println("title: ", title)
	// fmt.Println("company: ", company)
	// fmt.Println("location: ", location)
	// fmt.Println("salary: ", salary)
	// fmt.Println("=============")
	return extractedJob{
		id:       id,
		title:    title,
		company:  company,
		location: location,
		salary:   salary,
	}
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkError(err)
	checkCode(res.StatusCode)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	defer res.Body.Close()
	return pages
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(statusCode int) {
	if statusCode != 200 {
		log.Fatalln("Request failed with Status: ", statusCode)
	}
}

func cleanString(origin string) string {
	// return strings.TrimSpace(origin)
	return strings.Join(strings.Fields(strings.TrimSpace(origin)), " ")

}
