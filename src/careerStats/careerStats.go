package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gocolly/colly/v2"
)

type Player struct {
	Player_id                string `dynamodbav:"player-id"`
	Player_database_key      string `dynamodbav:"player-database-key"`
	Age                      string `dynamodbav:"age"`
	Team                     string `dynamodbav:"team"`
	Games                    string `dynamodbav:"games"`
	Games_started            string `dynamodbav:"games-started"`
	Minutes_played           string `dynamodbav:"minutes"`
	Field_goals              string `dynamodbav:"fg"`
	Field_goals_attempted    string `dynamodbav:"fga"`
	Three_pointers           string `dynamodbav:"3p"`
	Three_pointers_attempted string `dynamodbav:"3pa"`
	Free_throws              string `dynamodbav:"ft"`
	Free_throws_attempted    string `dynamodbav:"fta"`
	Offensive_rebounds       string `dynamodbav:"oreb"`
	Defensive_rebounds       string `dynamodbav:"dreb"`
	Assists                  string `dynamodbav:"ast"`
	Steals                   string `dynamodbav:"stl"`
	Blocks                   string `dynamodbav:"blk"`
	Turnovers                string `dynamodbav:"to"`
	Fouls                    string `dynamodbav:"fls"`
	Points                   string `dynamodbav:"pts"`
	Awards                   string `dynamodbav:"awards"`
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
		TableName: aws.String("career-stats-database"),
		Item:      av,
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}

	return nil
}

func getPlayerData(url string, ctx context.Context, wg *sync.WaitGroup, countMutex *sync.Mutex) {
	c := colly.NewCollector()

	c.OnHTML("table#per_game_stats tbody tr", func(e *colly.HTMLElement) {
		player := Player{
			Player_id:                fmt.Sprintf("https://www.basketball-reference.com%s", e.ChildAttr("th[data-stat='year_id'] a", "href")),
			Player_database_key:      url,
			Age:                      e.ChildText("td[data-stat='age']"),
			Team:                     e.ChildText("td[data-stat='team_name_abbr']"),
			Games:                    e.ChildText("td[data-stat='games']"),
			Games_started:            e.ChildText("td[data-stat='games_started']"),
			Minutes_played:           e.ChildText("td[data-stat='mp_per_g']"),
			Field_goals:              e.ChildText("td[data-stat='fg_per_g']"),
			Field_goals_attempted:    e.ChildText("td[data-stat='fga_per_g']"),
			Three_pointers:           e.ChildText("td[data-stat='fg3_per_g']"),
			Three_pointers_attempted: e.ChildText("td[data-stat='fg3a_per_g']"),
			Free_throws:              e.ChildText("td[data-stat='ft_per_g']"),
			Free_throws_attempted:    e.ChildText("td[data-stat='fta_per_g']"),
			Offensive_rebounds:       e.ChildText("td[data-stat='orb_per_g']"),
			Defensive_rebounds:       e.ChildText("td[data-stat='drb_per_g']"),
			Assists:                  e.ChildText("td[data-stat='ast_per_g']"),
			Steals:                   e.ChildText("td[data-stat='stl_per_g']"),
			Blocks:                   e.ChildText("td[data-stat='blk_per_g']"),
			Turnovers:                e.ChildText("td[data-stat='tov_per_g']"),
			Fouls:                    e.ChildText("td[data-stat='pf_per_g']"),
			Points:                   e.ChildText("td[data-stat='pts_per_g']"),
			Awards:                   e.ChildText("td[data-stat='awards']"),
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
		fmt.Println(player)
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
		go getPlayerData(url, ctx, &wg, &countMutex)
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
