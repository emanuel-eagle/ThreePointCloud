import boto3
import os

TABLE = os.environ["TABLE_NAME"]

def handler(event, context):
    client = boto3.client('dynamodb')
    print(TABLE)
    dynamodb_response = client.scan(TableName = TABLE)

    return {
        'statusCode': 200,
        'body': dynamodb_response
    }