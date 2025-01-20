package tests

import (
	"all-about-the-market/backend/handlers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock the Gin context for testing
func performRequest(router http.Handler, method, path string, body interface{}) *http.Response {
	// Marshal body to JSON
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr.Result()
}

func cleanupUser(t *testing.T, username string, cognitoClient *cognitoidentityprovider.CognitoIdentityProvider) {
	_, err := cognitoClient.AdminDeleteUser(&cognitoidentityprovider.AdminDeleteUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String("us-east-2_z0ClJwBpB"), // Replace with your UserPoolId
	})
	if err != nil {
		t.Errorf("Error cleaning up user: %v", err)
	} else {
		t.Logf("Successfully cleaned up user: %s", username)
	}
}

func TestRegisterUser(t *testing.T) {
	// Initialize a real or mock Cognito Client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)

	username := fmt.Sprintf("testuser_%d@example.com", time.Now().Unix())
	password := "TestPassword123!"

	defer cleanupUser(t, username, cognitoClient)

	// Initialize the Gin router
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})

	// Prepare the request body
	requestBody := map[string]string{
		"username": username,
		"password": password,
	}

	// Perform the request to register the user
	response := performRequest(router, "POST", "/register", requestBody)

	// Check if the status code is OK (200)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Check the response body
	var responseBody map[string]string
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "User registered and logged in successfully", responseBody["message"])
}

func TestRegisterExistingUser(t *testing.T) {
	// Initialize a real or mock Cognito Client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)

	// Initialize the Gin router
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})

	// Prepare the request body
	requestBody := map[string]string{
		"username": "aramkazorian@gmail.com",
		"password": "Kazorian1!",
	}

	// Perform the request to register the user
	response := performRequest(router, "POST", "/register", requestBody)

	// Decode the response body
	var responseBody map[string]string
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Check if the error message matches one of the expected messages
	expectedCustomMessage := "Error registering user"
	awsErrorPrefix := "User already exists" // Customize this based on your AWS error message

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled existing user error")
	} else {
		t.Fatalf("Unexpected error message: %v", responseBody["error"])
	}
}

func TestInvalidEmailFormat(t *testing.T) {

	// Initialize a real or mock Cognito Client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)

	// Initialize the Gin router
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})

	requestBody := map[string]string{
		"username": "invalid-email",
		"password": "Kazorian1!",
	}
	response := performRequest(router, "POST", "/register", requestBody)
	assert.Equal(t, 400, response.StatusCode)
	var responseBody map[string]string
	json.NewDecoder(response.Body).Decode(&responseBody)

	// Check if the error message matches one of the expected messages
	expectedCustomMessage := "Invalid Email."
	awsErrorPrefix := "Username should be an email." // Customize this based on your AWS error message

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled existing user error")
	} else {
		t.Fatalf("Unexpected error message: %v", responseBody["error"])
	}
}

func TestInvalidPassword(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})

	requestBody := map[string]string{
		"username": "testuser@example.com",
		"password": "short", // Invalid password
	}
	response := performRequest(router, "POST", "/register", requestBody)
	assert.Equal(t, 400, response.StatusCode)
	var responseBody map[string]string
	json.NewDecoder(response.Body).Decode(&responseBody)

	expectedCustomMessage := "Password is too weak."
	awsErrorPrefix := "Password did not conform with policy: Password not long enough"

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled invalid password error")
	} else {
		t.Fatalf("Unexpected error message: %v", responseBody["error"])
	}
}

func TestIncorrectLoginCredentials(t *testing.T) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})
	router.POST("/login", func(c *gin.Context) {
		handlers.Login(c, cognitoClient)
	})

	requestBody := map[string]string{
		"username": "aramkazorian@gmail.com",
		"password": "WrongPassword123",
	}

	response := performRequest(router, "POST", "/register", requestBody)
	assert.Equal(t, 400, response.StatusCode)
	var responseBody map[string]string
	json.NewDecoder(response.Body).Decode(&responseBody)

	expectedCustomMessage := "Password does not match"
	awsErrorPrefix := "Password did not conform with policy: Password must have symbol characters"

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled short password error")
	} else {
		t.Fatalf("Unexpected error message: %v", responseBody["error"])
	}

	requestBody = map[string]string{
		"username": "aramkazorian@gmail.com",
		"password": "WrongPassword123!",
	}

	response = performRequest(router, "POST", "/login", requestBody)
	assert.Equal(t, 400, response.StatusCode)
	json.NewDecoder(response.Body).Decode(&responseBody)

	expectedCustomMessage = "Password does not match"
	awsErrorPrefix = "Incorrect username or password."

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled invalid password error")
	} else {
		t.Fatalf("Unexpected error message: %v", responseBody["error"])
	}
}

// func TestMultipleRegistrations(t *testing.T) {
//     // Initialize a real or mock Cognito Client
//     sess := session.Must(session.NewSession(&aws.Config{
// 		Region: aws.String("us-east-2"),
// 	}))
// 	cognitoClient := cognitoidentityprovider.New(sess)
// 	router := gin.Default()
// 	router.POST("/register", func(c *gin.Context) {
// 		handlers.RegisterAndLoginUser(c, cognitoClient)
// 	})

//     // Register users
//     for i := 0; i < 100; i++ {
//         requestBody := map[string]string{
//             "username": fmt.Sprintf("user%d@example.com", i),
//             "password": "Password123!",
//         }

//         response := performRequest(router, "POST", "/register", requestBody)
//         assert.Equal(t, 200, response.StatusCode)
//     }

//     // Clean up users after the test
//     for i := 0; i < 100; i++ {
//         username := fmt.Sprintf("user%d@example.com", i)
//         cleanupUser(t, username, cognitoClient)
//     }
// }

