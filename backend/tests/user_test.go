package tests

import (
	"all-about-the-market/backend/handlers"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
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
	userPoolID := os.Getenv("COGNITO_USER_POOL_ID") // Retrieve User Pool ID from environment variables
	if userPoolID == "" {
		t.Fatal("COGNITO_USER_POOL_ID is not set")
	}

	_, err := cognitoClient.AdminDeleteUser(&cognitoidentityprovider.AdminDeleteUserInput{
		Username:   aws.String(username),
		UserPoolId: aws.String(userPoolID),
	})
	if err != nil {
		t.Errorf("Error cleaning up user: %v", err)
	} else {
		t.Logf("Successfully cleaned up user: %s", username)
	}
}

func TestRegisterUser(t *testing.T) {
	region := os.Getenv("COGNITO_REGION")
	if region == "" {
		t.Fatal("COGNITO_REGION is not set")
	}

	// Initialize a Cognito Client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)

	testEmailDomain := os.Getenv("TEST_EMAIL_DOMAIN")
	testValidPassword := os.Getenv("TEST_VALID_PASSWORD")

	username := fmt.Sprintf("testuser_%d@%s", time.Now().Unix(), testEmailDomain)
	password := testValidPassword

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
	fmt.Print("hello")
	region := os.Getenv("COGNITO_REGION")
	if region == "" {
		t.Fatal("COGNITO_REGION is not set")
	}

	// Initialize a Cognito Client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)

	// Initialize the Gin router
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})

	// Prepare the request body
	requestBody := map[string]string{
		"username": os.Getenv("TEST_EXISTING_USER_EMAIL"),
		"password": os.Getenv("TEST_EXISTING_USER_PASSWORD"),
	}

	// Perform the request to register the user
	response := performRequest(router, "POST", "/register", requestBody)

	// Decode the response body
	var responseBody map[string]string
	err := json.NewDecoder(response.Body).Decode(&responseBody)
	assert.NoError(t, err)

	// Check if the error message matches one of the expected messages
	expectedCustomMessage := "Error registering user"
	awsErrorPrefix := "User already exists"

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled existing user error")
	} else {
		t.Fatalf("Unexpected error message: %v", responseBody["error"])
	}
}

func TestInvalidEmailFormat(t *testing.T) {
	region := os.Getenv("COGNITO_REGION")
	if region == "" {
		t.Fatal("COGNITO_REGION is not set")
	}

	// Initialize a Cognito Client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	cognitoClient := cognitoidentityprovider.New(sess)

	// Initialize the Gin router
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		handlers.RegisterAndLoginUser(c, cognitoClient)
	})

	requestBody := map[string]string{
		"username": "invalid-email",
		"password": os.Getenv("TEST_VALID_PASSWORD"),
	}
	response := performRequest(router, "POST", "/register", requestBody)
	assert.Equal(t, 400, response.StatusCode)
	var responseBody map[string]string
	json.NewDecoder(response.Body).Decode(&responseBody)

	expectedCustomMessage := "Invalid Email."
	awsErrorPrefix := "Username should be an email."

	if responseBody["error"] == expectedCustomMessage ||
		strings.HasPrefix(responseBody["error"], awsErrorPrefix) {
		t.Log("Correctly handled invalid email error")
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

