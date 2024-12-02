package handlers

import (
    "github.com/gin-gonic/gin"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
    "net/http"
)

var cognitoClient = cognitoidentityprovider.New(session.Must(session.NewSession()), &aws.Config{
    Region: aws.String("us-west-2"),
})

// Combined handler for Register and Login
func RegisterAndLoginUser(c *gin.Context) {
    var input struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Try to register user (using Cognito)
    signUpInput := &cognitoidentityprovider.SignUpInput{
        Username: aws.String(input.Username),
        Password: aws.String(input.Password),
        ClientId: aws.String("your_cognito_app_client_id"), // Use your Cognito Client ID
    }
    _, err := cognitoClient.SignUp(signUpInput)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error registering user"})
        return
    }

    // After successful registration, initiate the login process
    loginInput := &cognitoidentityprovider.InitiateAuthInput{
        AuthFlow: aws.String("USER_PASSWORD_AUTH"),
        ClientId: aws.String("your_cognito_app_client_id"),
        AuthParameters: map[string]*string{
            "USERNAME": aws.String(input.Username),
            "PASSWORD": aws.String(input.Password),
        },
    }
    _, err = cognitoClient.InitiateAuth(loginInput)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User registered and logged in successfully"})
}
