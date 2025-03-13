import boto3
import os

TABLE = os.environ["TABLE_NAME"]
HASH_KEY = os.environ["HASH_KEY"]

def handler(event, context):
    client = boto3.client('dynamodb')
    print(TABLE)
    dynamodb_response = client.scan(TableName = TABLE,
                                    ProjectionExpression = HASH_KEY)

    return {
        'statusCode': 200,
        'body': dynamodb_response
    }