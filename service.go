package main

import "github.com/aws/aws-sdk-go/service/dynamodb"

type transactionReaderService struct {
	dynamodb *dynamodb.DynamoDB
}
