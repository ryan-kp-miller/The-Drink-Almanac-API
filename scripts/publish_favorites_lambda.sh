#!/bin/bash

APP_NAME="favorites-lambda"
ARTIFACT_NAME="${APP_NAME}.zip"

aws lambda update-function-code --function-name ryanm-learnathon-2023-favorites-lambda --zip-file "fileb://./artifacts/${ARTIFACT_NAME}"

