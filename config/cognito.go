package config

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
)

var CognitoClient *cognitoidentityprovider.CognitoIdentityProvider

func InitCognito() {
    CognitoClient = cognitoidentityprovider.New(session.Must(session.NewSession()), &aws.Config{
        Region: aws.String("us-west-2"), // Set the region
    })
}
