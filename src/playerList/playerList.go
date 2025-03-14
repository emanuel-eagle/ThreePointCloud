package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

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
	Active    string `dynamodbav:"active"`
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
		TableName: aws.String("bos-player-database"),
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
	c.OnHTML("table#franchise_register tbody tr", func(e *colly.HTMLElement) {
		active := "False"
		currentYearString := fmt.Sprintf("%d", time.Now().Year())
		if e.ChildText("td[data-stat='year_max']") == currentYearString {
			active = "True"
		}
		player := Player{
			Player_id: fmt.Sprintf("https://www.basketball-reference.com%s", e.ChildAttr("td[data-stat='player'] a", "href")),
			Name:      e.ChildText("td[data-stat='player']"),
			From:      e.ChildText("td[data-stat='year_min']"),
			To:        e.ChildText("td[data-stat='year_max']"),
			Active:    active,
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

	urls = []string{"https://www.basketball-reference.com/teams/BOS/players.html"}

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
