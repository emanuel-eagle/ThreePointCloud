import boto3
from bs4 import BeautifulSoup
import requests
import json

def handler(event, context):

    urls = event.get('urls', [])

    for url in urls:
        response = requests.get(url)
        print(f"{url} returned: {response.status_code}")



    return {
        'statusCode': 200,
        'body': "success"
    }