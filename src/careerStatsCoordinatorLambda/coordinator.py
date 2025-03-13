import boto3
import os

TABLE = os.environ["TABLE_NAME"]
HASH_KEY = os.environ["HASH_KEY"]

def handler(event, context):
    client = boto3.client('dynamodb')
    print(TABLE)
    dynamodb_response = client.scan(TableName = TABLE,
                                    ProjectionExpression='#hashKey',
                                    ExpressionAttributeNames={
                                        '#hashKey': HASH_KEY
                                    })
    
    items = dynamodb_response["Items"]

    for item in items:
        url = item[HASH_KEY]["S"]
        print(url)

    return {
        'statusCode': 200,
        'body': "success"
    }