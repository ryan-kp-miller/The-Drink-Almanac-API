#!/bin/bash

LAMBDA_PREFIX=$1
APP_NAME="${LAMBDA_PREFIX}-lambda"
ARTIFACT_NAME="${APP_NAME}.zip"

aws lambda update-function-code --function-name "ryanm-learnathon-2023-${LAMBDA_PREFIX}-lambda" --zip-file "fileb://./artifacts/${ARTIFACT_NAME}"

