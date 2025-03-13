package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gocolly/colly/v2"
)

type Player struct {
	Name     string
	From     string
	To       string
	Position string
	Height   string
	Weight   string
	DOB      string
	College  string
}

// Lambda event structure
type MyEvent struct {
	// Can be empty or customize as needed
}

// Response from Lambda
type MyResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

func getPlayers(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	c := colly.NewCollector()
	c.OnHTML("table#players tbody tr", func(e *colly.HTMLElement) {
		player := Player{
			Name:     e.ChildText("th[data-stat='player']"),
			From:     e.ChildText("td[data-stat='year_min']"),
			To:       e.ChildText("td[data-stat='year_max']"),
			Position: e.ChildText("td[data-stat='pos']"),
			Height:   e.ChildText("td[data-stat='height']"),
			Weight:   e.ChildText("td[data-stat='weight']"),
			DOB:      e.ChildText("td[data-stat='birth_date']"),
			College:  e.ChildText("td[data-stat='colleges']"),
		}
		fmt.Println(player) //test
	})
	c.Visit(url)
}

func HandleRequest(ctx context.Context, event MyEvent) (MyResponse, error) {

	var urls []string

	for r := 'a'; r <= 'z'; r++ {
		urls = append(urls, fmt.Sprintf("https://www.basketball-reference.com/players/%s/", string(r)))
	}

	var wg sync.WaitGroup

	wg.Add(len(urls))

	for _, url := range urls {
		go getPlayers(url, &wg)
	}

	wg.Wait()

	return MyResponse{
		Message: "Successfully scraped basketball players",
		Count:   1,
	}, nil
}

func main() {
	// This is the Lambda entry point
	lambda.Start(HandleRequest)
}
