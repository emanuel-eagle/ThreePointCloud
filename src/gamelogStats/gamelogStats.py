import boto3
from boto3.dynamodb.types import TypeSerializer
from bs4 import BeautifulSoup
import requests
import json
import random
import time
import os

TABLE = os.environ["TABLE_NAME"]

user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
    "Mozilla/5.0 (X11; CrOS x86_64 8172.45.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.64 Safari/537.36",
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_2) AppleWebKit/601.3.9 (KHTML, like Gecko) Version/9.0.2 Safari/601.3.9",
    "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:15.0) Gecko/20100101 Firefox/15.0.1"
]

def handler(event, context):

    dynamodb_client = boto3.client("dynamodb")

    urls = event['urls']

    serializer = TypeSerializer()

    for url in urls:
        time.sleep(random.random())
        response = requests.get(url, headers = {'User-agent': user_agents[random.randint(0, len(user_agents)-1)]})
        status_code = response.status_code
        while status_code == 429:
            time.sleep(20)
            print("Retrying")
            response = requests.get(url, headers = {'User-agent': user_agents[random.randint(0, len(user_agents)-1)]})
            status_code = response.status_code
        soup = BeautifulSoup(response.text, "html.parser")
        table = soup.find(id = "pgl_basic")
        tbody = table.find("tbody")
        for tr in tbody.find_all("tr"):
            td_list = tr.find_all("td")
            game_stats = {}
            for td in td_list:
                game_stats[td['data-stat']] = td.text
            if game_stats:
                game_stats["player-database-key"] = url
                game_stats["game-id"] = f"{url}-{game_stats['date_game']}-{game_stats['age']}"
                dynamodb_item = {k: serializer.serialize(v) for k, v in game_stats.items()}
                response = dynamodb_client.put_item(
                    TableName=TABLE,
                    Item=dynamodb_item
                )

    return {
        'statusCode': 200,
        'body': "success"
    }

handler('', '')