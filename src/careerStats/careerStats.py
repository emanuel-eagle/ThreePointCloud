import boto3
from bs4 import BeautifulSoup
import requests
import json

def handler(event, context):

    urls = json.loads(event)

    print(urls)


    return {
        'statusCode': 200,
        'body': "success"
    }
