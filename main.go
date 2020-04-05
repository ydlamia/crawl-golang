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
	c2 := make(chan []extractedJob)
	var jobs []extractedJob
	totalPages := getPages()

	// getPage((totalPages + 1) - (totalPages))

	for i := 0; i < totalPages; i++ {
		go getPage(i, c2)
	}
	for i := 0; i < totalPages; i++ {
		extractedJobs := <-c2
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

func getPage(page int, c2 chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)
	pageURL := baseURL + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting ", pageURL)
	res, err := http.Get(pageURL)
	checkError(err)
	checkCode(res.StatusCode)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	defer res.Body.Close()

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})
	for i := 0; i < searchCards.Length(); i++ {
		job := <-c
		jobs = append(jobs, job)
	}
	c2 <- jobs
}

func extractJob(card *goquery.Selection, c chan<- extractedJob) {
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
	c <- extractedJob{
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
