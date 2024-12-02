package utils

import (
    "github.com/dgrijalva/jwt-go"
    "time"
    "os"
)

var secretKey = []byte("your_secret_key") // Use environment variables for sensitive info

func GenerateJWT(userId string) (string, error) {
    claims := jwt.MapClaims{
        "userId": userId,
        "exp":    time.Now().Add(time.Hour * 72).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(secretKey)
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return secretKey, nil
    })
}
