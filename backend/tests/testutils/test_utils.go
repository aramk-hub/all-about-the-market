package testutils

import (
	"crypto/rand"
	"crypto/rsa"
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
	"all-about-the-market/backend/utils"
)

var mockPrivateKey *rsa.PrivateKey

func init() {
	// Generate a mock RSA private key for signing tokens
	var err error
	mockPrivateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("failed to generate mock RSA private key")
	}
}

func GenerateExpiredToken() (string, error) {
	// Use your RSA private key here for signing
	signingKey := []byte("your-signing-key")

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "162ibtbklb9e1gmjlch8ope7b4",
		"exp": time.Now().Add(-time.Hour).Unix(), // Set expiration in the past
	})

	// Sign the token
	return token.SignedString(signingKey)
}

func GenerateTokenWithInvalidAudience() (string, error) {
	// Use your RSA private key here for signing
	signingKey := []byte("your-signing-key")

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "invalid-audience", // Invalid audience
		"exp": time.Now().Add(time.Hour).Unix(),
	})

	// Sign the token
	return token.SignedString(signingKey)
}

// GenerateMockPublicKey generates a public key from the mock private key
func GenerateMockPublicKey() *rsa.PublicKey {
	return &mockPrivateKey.PublicKey
}

// MockJWKS populates the JWKS cache with the mock public key
func MockJWKS() {
	utils.CachedJWKSmu.Lock()
	defer utils.CachedJWKSmu.Unlock()

	utils.CachedJWKS = map[string]*rsa.PublicKey{
		"mockKid": GenerateMockPublicKey(), // Use "mockKid" as the key ID
	}

	// Debugging: Print the cached JWKS
	fmt.Printf("MockJWKS: %+v\n", utils.CachedJWKS)
}

// GenerateValidToken generates a valid JWT signed with the mock private key
func GenerateValidToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": "1234567890",
		"aud": "162ibtbklb9e1gmjlch8ope7b4", // Replace with your Cognito Client ID
		"exp": time.Now().Add(time.Hour).Unix(), // Token expires in 1 hour
		"iat": time.Now().Unix(),               // Issued at
		"name": "John Doe",
	})

	// Set the "kid" header to match the mocked JWKS
	token.Header["kid"] = "mockKid"

	// Debugging: Print the token's header
	fmt.Printf("Generated Token Header: %+v\n", token.Header)

	// Sign the token with the mock private key
	return token.SignedString(mockPrivateKey)
}

// // MockFetchJWKS overrides the FetchJWKS function for testing
func MockFetchJWKS() {
	utils.FetchJWKSFunc = func() error {
		MockJWKS()
		fmt.Println("MockFetchJWKS called")
		return nil
	}
}
