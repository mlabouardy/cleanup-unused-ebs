#!/bin/bash

echo "Building binary"
GOOS=linux GOARCH=amd64 go build -o main main.go

echo "Preparing deployment package"
zip deployment.zip main

echo "Deploying"
aws lambda create-function --function-name RemoveUnusedEBS \
     --runtime go1.x --handler main \
     --role arn:aws:iam::ACCOUNT_ID:role/RemoveUnusedEBS \
     --zip-file fileb://./deployment.zip

echo "Cleaning up "
rm deployment.zip main