package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
)

// Combined handler for Register and Login
func RegisterAndLoginUser(c *gin.Context, cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Try to register the user
	signUpInput := &cognitoidentityprovider.SignUpInput{
		Username: aws.String(input.Username),
		Password: aws.String(input.Password),
		ClientId: aws.String("162ibtbklb9e1gmjlch8ope7b4"),
	}
	_, err := cognitoClient.SignUp(signUpInput)
	if err != nil {
		handleAWSError(c, "Error registering user", err)
		return
	}

	// Confirm the user registration
	adminConfirmInput := &cognitoidentityprovider.AdminConfirmSignUpInput{
		Username:   aws.String(input.Username),
		UserPoolId: aws.String("us-east-2_z0ClJwBpB"),
	}
	_, err = cognitoClient.AdminConfirmSignUp(adminConfirmInput)
	if err != nil {
		handleAWSError(c, "Error confirming user", err)
		return
	}

	// Log in the user
	loginInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String("162ibtbklb9e1gmjlch8ope7b4"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(input.Username),
			"PASSWORD": aws.String(input.Password),
		},
	}
	_, err = cognitoClient.InitiateAuth(loginInput)
	if err != nil {
		handleAWSError(c, "Invalid credentials", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered and logged in successfully"})
}

func Login(c *gin.Context, cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Log in the user
	loginInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String("162ibtbklb9e1gmjlch8ope7b4"),
		AuthParameters: map[string]*string{
			"USERNAME": aws.String(input.Username),
			"PASSWORD": aws.String(input.Password),
		},
	}
	_, err := cognitoClient.InitiateAuth(loginInput)
	if err != nil {
		handleAWSError(c, "Invalid credentials", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered and logged in successfully"})

}

func GetUserFromCognito(username string) {
	sess := session.Must(session.NewSession())
	svc := cognitoidentityprovider.New(sess, &aws.Config{
		Region: aws.String("us-east-2"), // Change to your region
	})

	input := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: aws.String("us-east-2_z0ClJwBpB"), // Replace with your Cognito User Pool ID
		Username:   aws.String(username),
	}

	result, err := svc.AdminGetUser(input)
	if err != nil {
		log.Println("Error getting user:", err)
		return
	}

	fmt.Println("User found:", result)
}

func handleAWSError(c *gin.Context, defaultMessage string, err error) {
	if awsErr, ok := err.(awserr.Error); ok {
		// Check for the error code returned by AWS
		awsCode := awsErr.Code()
		statusCode := http.StatusInternalServerError // Default to 500

		// Map AWS error codes to HTTP status codes as needed
		if awsCode == "UsernameExistsException" {
			statusCode = http.StatusConflict // Use 409 for username conflict
		} else if awsCode == "InvalidParameterException" || 
            awsCode == "InvalidPasswordException" || 
            awsCode == "NotAuthorizedException" {
			statusCode = http.StatusBadRequest // Use 400 for invalid input
		}

		c.JSON(statusCode, gin.H{"error": awsErr.Message()})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": defaultMessage})
	}
}

