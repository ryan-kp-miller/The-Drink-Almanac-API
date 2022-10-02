#!/bin/bash

echo "################## Creating the-drink-almanac-users table and inserting data ##################"
awslocal dynamodb --endpoint-url=http://localhost:4566 create-table \
    --table-name the-drink-almanac-users\
    --attribute-definitions \
        AttributeName=id,AttributeType=N \
    --key-schema \
        AttributeName=id,KeyType=HASH \
--provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5

awslocal dynamodb put-item --table-name the-drink-almanac-users --item \
    "{
        \"id\":{\"N\":\"0\"},
        \"username\":{\"S\":\"test0\"}, 
        \"password\":{\"S\":\"test0\"}
    }"
awslocal dynamodb put-item --table-name the-drink-almanac-users --item \
    "{
        \"id\":{\"N\":\"1\"},
        \"username\":{\"S\":\"test1\"},
        \"password\":{\"S\":\"test1\"}
    }"
awslocal dynamodb put-item --table-name the-drink-almanac-users --item \
    "{
        \"id\":{\"N\":\"2\"},
        \"username\":{\"S\":\"test2\"},
        \"password\":{\"S\":\"test2\"}
    }"

echo "################## Creating the-drink-almanac-users table and inserting data ##################"
awslocal dynamodb --endpoint-url=http://localhost:4566 create-table \
    --table-name the-drink-almanac-favorites \
    --attribute-definitions \
        AttributeName=id,AttributeType=N \
        AttributeName=user_id,AttributeType=N \
    --key-schema \
        AttributeName=id,KeyType=HASH \
    --global-secondary-indexes \
    "[{\"IndexName\": \"user-index\",\"KeySchema\":[{\"AttributeName\":\"user_id\",\"KeyType\":\"HASH\"}],\"Projection\": {\"ProjectionType\": \"ALL\"},\"ProvisionedThroughput\": {
                    \"WriteCapacityUnits\": 5,
                    \"ReadCapacityUnits\": 10
                }}]" \
    --provisioned-throughput \
            ReadCapacityUnits=10,WriteCapacityUnits=5

awslocal dynamodb put-item --table-name the-drink-almanac-favorites --item \
    "{
        \"id\":{\"N\":\"0\"},
        \"drink_id\":{\"N\":\"0\"},
        \"user_id\":{\"N\":\"0\"}
    }" 
awslocal dynamodb put-item --table-name the-drink-almanac-favorites --item \
    "{
        \"id\":{\"N\":\"1\"},
        \"drink_id\":{\"N\":\"1\"},
        \"user_id\":{\"N\":\"0\"}
    }"
awslocal dynamodb put-item --table-name the-drink-almanac-favorites --item \
    "{
        \"id\":{\"N\":\"2\"},
        \"drink_id\":{\"N\":\"1\"},
        \"user_id\":{\"N\":\"1\"}
    }"
