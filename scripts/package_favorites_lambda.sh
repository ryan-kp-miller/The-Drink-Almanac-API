#!/bin/bash

ROOT_PATH="./go_api/lambdas/favorites"
APP_NAME="favorites-lambda"
ARTIFACT_NAME="${APP_NAME}"
START_PATH=$(pwd)

# Assure the artifacts directory exists
mkdir -p artifacts/

# Run go build for a linux artifact in the app directory
cd ${ROOT_PATH}
GOOS=linux go build -v -ldflags '-s -w' -a -o ${APP_NAME} ./
zip ${ARTIFACT_NAME}.zip ${APP_NAME}
rm ${APP_NAME}
cd ${START_PATH}

# Move the build artifact to the proper directory
mv ${ROOT_PATH}/${ARTIFACT_NAME}.zip artifacts/${ARTIFACT_NAME}.zip
