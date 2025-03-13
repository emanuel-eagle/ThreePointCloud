package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gocolly/colly/v2"
)

type Player struct {
	Player_id string `dynamodbav:"player-id"`
	Name      string `dynamodbav:"name"`
	From      string `dynamodbav:"from_year"`
	To        string `dynamodbav:"to_year"`
	Position  string `dynamodbav:"position"`
	Height    string `dynamodbav:"height"`
	Weight    string `dynamodbav:"weight"`
	DOB_Month string `dynamodbav:"dob_month"`
	DOB_Day   string `dynamodbav:"dob_day"`
	DOB_Year  string `dynamodbav:"dob_year"`
	College   string `dynamodbav:"college"`
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

var dbClient *dynamodb.Client
var playerCount int

// Initialize the DynamoDB client
func init() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	dbClient = dynamodb.NewFromConfig(cfg)
}

func addPlayerToDynamoDB(player Player, ctx context.Context) error {
	// Marshal the player to DynamoDB attribute values
	av, err := attributevalue.MarshalMap(player)
	if err != nil {
		return fmt.Errorf("failed to marshal player: %w", err)
	}

	// Put the item in DynamoDB
	_, err = dbClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("player-database"),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func getPlayers(url string, ctx context.Context, wg *sync.WaitGroup, countMutex *sync.Mutex) {
	defer wg.Done()
	c := colly.NewCollector()
	c.OnHTML("table#players tbody tr", func(e *colly.HTMLElement) {
		player := Player{
			Player_id: fmt.Sprintf("%s_%s", e.ChildText("th[data-stat='player']"), e.ChildText("td[data-stat='year_min']")),
			Name:      e.ChildText("th[data-stat='player']"),
			From:      e.ChildText("td[data-stat='year_min']"),
			To:        e.ChildText("td[data-stat='year_max']"),
			Position:  e.ChildText("td[data-stat='pos']"),
			Height:    e.ChildText("td[data-stat='height']"),
			Weight:    e.ChildText("td[data-stat='weight']"),
			DOB_Month: strings.Split(e.ChildText("td[data-stat='birth_date']"), " ")[0],
			DOB_Day:   strings.Split(strings.Split(e.ChildText("td[data-stat='birth_date']"), " ")[1], ",")[0],
			DOB_Year:  strings.Split(strings.Split(e.ChildText("td[data-stat='birth_date']"), " ")[1], ",")[1],
			College:   e.ChildText("td[data-stat='colleges']"),
		}
		if player.Name != "" {
			// Add player to DynamoDB right away
			err := addPlayerToDynamoDB(player, ctx)
			if err != nil {
				log.Printf("Error adding player %s: %v", player.Name, err)
			} else {
				// Increment the player count safely
				countMutex.Lock()
				playerCount++
				countMutex.Unlock()
			}
		}
	})
	c.Visit(url)
}

func HandleRequest(ctx context.Context, event MyEvent) (MyResponse, error) {

	var urls []string
	playerCount = 0 // Reset count for each invocation

	for r := 'a'; r <= 'z'; r++ {
		urls = append(urls, fmt.Sprintf("https://www.basketball-reference.com/players/%s/", string(r)))
	}

	var wg sync.WaitGroup
	var countMutex sync.Mutex

	wg.Add(len(urls))

	for _, url := range urls {
		go getPlayers(url, ctx, &wg, &countMutex)
	}

	wg.Wait()

	return MyResponse{
		Message: "Successfully scraped basketball players",
		Count:   playerCount,
	}, nil
}

func main() {
	// This is the Lambda entry point
	lambda.Start(HandleRequest)
}
