import boto3
import json
import os

TABLE = os.environ["TABLE_NAME"]
HASH_KEY = os.environ["HASH_KEY"]
CAREER_STATS_LAMBDA = os.environ["CAREER_STATS_LAMBDA"]

def split_list_into_chunks(items, target_chunk_size=500):
    # Calculate how many chunks we need
    num_chunks = max(1, (len(items) + target_chunk_size - 1) // target_chunk_size)
    
    # Calculate actual chunk size (distributing items evenly)
    chunk_size = (len(items) + num_chunks - 1) // num_chunks
    
    # Split the list
    result = []
    for i in range(0, len(items), chunk_size):
        result.append(items[i:i + chunk_size])
    
    return result

def handler(event, context):
    dynamodb_client = boto3.client('dynamodb')
    lambda_client = boto3.client('lambda')
    dynamodb_response = dynamodb_client.scan(TableName = TABLE,
                                    ProjectionExpression='#hashKey',
                                    ExpressionAttributeNames={
                                        '#hashKey': HASH_KEY
                                    })
    
    items = dynamodb_response["Items"]

    urls = [item[HASH_KEY]["S"] for item in items]

    chunks = split_list_into_chunks(urls)

    for chunk in chunks:
        payload = {
            'urls' : chunk
        }
        response = lambda_client.invoke(
            FunctionName = CAREER_STATS_LAMBDA,
            InvocationType = "RequestResponse",
            Payload = json.dumps(payload)
        )
        print(response)


    return {
        'statusCode': 200,
        'body': "success"
    }