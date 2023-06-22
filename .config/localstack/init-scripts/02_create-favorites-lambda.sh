#!/bin/bash

echo "################################## Creating Favorites Lambda ##################################"
FAVORITES_LAMBDA=$(aws lambda create-function \
    --region=us-east-1 \
    --function-name=favorites-lambda \
    --handler=favorites-lambda \
    --runtime=go1.x \
    --zip-file=fileb:///artifacts/favorites-lambda.zip \
    --memory-size=128 \
    --role=arn:aws:iam::123456:role/role-name \
    --environment="Variables={ENVIRONMENT=local,AWS_ENDPOINT=http://localhost:4566}" \
    --endpoint-url=http://localhost:4566 \
    --profile=localstack \
    --output=text \
    --query="FunctionArn")
echo "Favorites Lambda ARN: ${FAVORITES_LAMBDA}"
