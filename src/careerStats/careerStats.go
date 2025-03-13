package main

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

type Player struct {
	Player_id                string `dynamodbav:"player-id"`
	Player_database_key      string `dynamodbav:"player-database-key"`
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

func main() {
	c := colly.NewCollector()
	var url = "https://www.basketball-reference.com/players/j/jamesle01.html"

	c.OnHTML("table#per_game_stats tbody tr", func(e *colly.HTMLElement) {
		player := Player{
			Player_id:           fmt.Sprintf("https://www.basketball-reference.com%s", e.ChildAttr("th[data-stat='year_id'] a", "href")),
			Player_database_key: url,
			Games 
		}
		fmt.Println(player)
	})

	c.Visit(url)
}
