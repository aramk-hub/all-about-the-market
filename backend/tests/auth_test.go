package tests

import (
	"all-about-the-market/backend/tests/testutils"
	"all-about-the-market/backend/utils"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Load environment variables from .env
	err := godotenv.Load("../cmd/.env") // Adjust path as needed
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	log.Println("Environment variables loaded successfully")

	// Run tests
	os.Exit(m.Run())
}

func TestMockJWKSSetup(t *testing.T) {
	testutils.MockJWKS()

	utils.CachedJWKSmu.RLock()
	defer utils.CachedJWKSmu.RUnlock()

	fmt.Printf("Cached JWKS: %+v\n", utils.CachedJWKS)
	assert.NotNil(t, utils.CachedJWKS["mockKid"], "Public key for mockKid should be present in JWKS cache")
}

func TestValidTokenWithMockJWKS(t *testing.T) {
	// Mock the JWKS
	testutils.MockFetchJWKS()

	// Debug: Check the cached JWKS after mocking
	utils.CachedJWKSmu.RLock()
	fmt.Printf("Cached JWKS after mocking: %+v\n", utils.CachedJWKS)
	utils.CachedJWKSmu.RUnlock()

	// Generate a valid token signed with the mock private key
	validToken, err := testutils.GenerateValidToken()
	assert.NoError(t, err, "Failed to generate a valid token")

	// Debug: Print the token's header and claims
	token, _ := jwt.Parse(validToken, nil)
	fmt.Printf("Generated Token Header: %+v\n", token.Header)

	// Validate the token
	token, err = utils.ValidateToken(validToken)
	if err != nil {
		log.Printf("Validation error: %v", err)
	}
	assert.NoError(t, err, "Valid token with mocked JWKS should not return an error")
	assert.NotNil(t, token, "Valid token should return a parsed token")
}

func TestInvalidToken(t *testing.T) {
	// Mock an invalid token
	invalidToken := "invalid.token.string"

	token, err := utils.ValidateToken(invalidToken)
	assert.Error(t, err, "Invalid token should return an error")
	assert.Nil(t, token, "Invalid token should not return a valid parsed token")
}

func TestExpiredToken(t *testing.T) {
	expiredToken, _ := testutils.GenerateExpiredToken()
	token, err := utils.ValidateToken(expiredToken)
	assert.Error(t, err, "Expired token should return an error")
	assert.Nil(t, token, "Expired token should not return a valid parsed token")
}

func TestInvalidAudience(t *testing.T) {
	tokenWithInvalidAudience, _ := testutils.GenerateTokenWithInvalidAudience()
	token, err := utils.ValidateToken(tokenWithInvalidAudience)
	assert.Error(t, err, "Token with invalid audience should return an error")
	assert.Nil(t, token, "Token with invalid audience should not return a valid parsed token")
}
