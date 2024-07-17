//go:generate mockery --name=DDBClient --output=./ --outpkg=client --filename=dynamodb_mock.go --inpackage
package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DDBClient interface {
	Scan(context.Context, *dynamodb.ScanInput, ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
	Query(context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

// CreateLocalDDBClient creates a dynamodb client using environment variables
func CreateLocalDDBClient(awsEndpoint string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		// only set the endpoint resolver if an endpoint is provided
		if awsEndpoint != "" {
			o.EndpointResolver = dynamodb.EndpointResolverFromURL(awsEndpoint)
		}
	}), nil
}
