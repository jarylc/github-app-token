package main

import (
	"encoding/base64"
	"fmt"
	"github-app-token/http/github"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"time"
)

var githubIssuer = "121562"

func main() {
	if len(os.Args) == 1 {
		errored("format: {base64 private key} {name of org (optional, defaults to current)}", nil)
	}
	key := os.Args[1]
	owner := ""
	if len(os.Args) >= 3 {
		owner = os.Args[2]
	}
	if owner == "" {
		owner = os.Getenv("GITHUB_REPOSITORY_OWNER")
		if owner == "" {
			errored("owner missing", nil)
		}
	}

	pem, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		errored("decoding base64 key", err)
	}
	jwt, err := generateJWT(pem)
	if err != nil {
		errored("generating JWT", err)
	}
	integrations, err := github.GetIntegrations(jwt)
	if err != nil {
		errored("getting integrations from API", err)
	}
	integration, ok := integrations[owner]
	if !ok {
		errored("organization integration not found for app", nil)
	}
	accessToken, err := github.GetAccessToken(jwt, integration)
	if err != nil {
		errored("generating access token from API", err)
	}

	fmt.Printf("::add-mask::%s\n", accessToken)
	fmt.Printf("::set-output name=token::%s", accessToken)
}

func generateJWT(pem []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	claims := jwt.StandardClaims{
		ExpiresAt: now + 60,
		IssuedAt:  now - 60,
		Issuer:    githubIssuer,
	}
	jwtGenerator := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwt, err := jwtGenerator.SignedString(key)
	if err != nil {
		return "", err
	}
	return jwt, nil
}

func errored(msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("::error::%s", msg)
	os.Exit(1)
}
