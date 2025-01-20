package utils

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	// CachedJWKS is the exported variable for the JWKS cache
	CachedJWKS   map[string]*rsa.PublicKey
	// CachedJWKSmu is the exported mutex for the JWKS cache
	CachedJWKSmu sync.RWMutex

	// FetchJWKSFunc is the exported function used to fetch the JWKS (overridable for testing)
	FetchJWKSFunc = FetchJWKS
)

// FetchJWKS fetches the JWKS from Cognito and caches the public keys
func FetchJWKS() error {
	jwksURL := os.Getenv("JWKS_URL")
	if jwksURL == "" {
		return errors.New("JWKS_URL is empty. Make sure it is set in the environment variables")
	}

	log.Printf("Fetching JWKS from URL: %s", jwksURL)

	resp, err := http.Get(jwksURL)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch JWKS: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read JWKS response body: %v", err)
	}

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return fmt.Errorf("failed to unmarshal JWKS JSON: %v", err)
	}

	CachedJWKSmu.Lock()
	defer CachedJWKSmu.Unlock()

	CachedJWKS = make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		nBytes, err := jwt.DecodeSegment(key.N)
		if err != nil {
			return fmt.Errorf("failed to decode modulus (N): %v", err)
		}
		eBytes, err := jwt.DecodeSegment(key.E)
		if err != nil {
			return fmt.Errorf("failed to decode exponent (E): %v", err)
		}

		pubKey := &rsa.PublicKey{
			N: new(big.Int).SetBytes(nBytes),
			E: int(new(big.Int).SetBytes(eBytes).Int64()),
		}
		CachedJWKS[key.Kid] = pubKey
	}

	log.Printf("Successfully fetched and cached %d JWKS keys", len(jwks.Keys))
	return nil
}

// GetPublicKey retrieves the public key for a given kid from the cached JWKS
func GetPublicKey(kid string) (*rsa.PublicKey, error) {
	CachedJWKSmu.RLock()
	defer CachedJWKSmu.RUnlock()

	fmt.Printf("Looking for kid: %s in JWKS\n", kid) // Debugging: Print the kid
	if key, ok := CachedJWKS[kid]; ok {
		fmt.Printf("Public key found for kid: %s\n", kid)
		return key, nil
	}
	fmt.Printf("Public key not found for kid: %s\n", kid)
	return nil, fmt.Errorf("public key not found for kid: %s", kid)
}

// ValidateToken validates a JWT using the cached JWKS keys
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Explicitly call FetchJWKSFunc to fetch JWKS
	log.Println("Calling FetchJWKSFunc...")
	if err := FetchJWKSFunc(); err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %v", err)
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		fmt.Printf("Token Header: %+v\n", token.Header) // Debugging: Print the token header

		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header missing in token")
		}

		fmt.Printf("Looking up public key for kid: %s\n", kid) // Debugging: Print the kid
		return GetPublicKey(kid)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expectedAud := os.Getenv("COGNITO_CLIENT_ID")
		log.Printf("Expected Audience: %s", expectedAud)
		log.Printf("Token Claims: %+v", claims)
	
		if claims["aud"] != expectedAud {
			return nil, fmt.Errorf("invalid audience: expected %s, got %s", expectedAud, claims["aud"])
		}
	
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return nil, fmt.Errorf("token is expired")
			}
		}
	
		log.Printf("Validated token claims: %v", claims)
	} else {
		return nil, errors.New("failed to parse token claims")
	}

	return token, nil
}
