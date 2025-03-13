import boto3
from bs4 import BeautifulSoup
import requests
import json

def handler(event, context):

    urls = json.loads(event)

    for url in urls:
        response = requests.get(url)



    return {
        'statusCode': 200,
        'body': "success"
    }
