package main

import (
	"encoding/base64"
	"fmt"
	"github-app-token/http/github"
	"github.com/dgrijalva/jwt-go"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) == 1 {
		errored("format: {app name/id} {base64 private key} {org name (optional, defaults to current)}", nil)
	}
	app := os.Args[1]
	if app == "" {
		errored("app name missing", nil)
	}
	_, isID := strconv.Atoi(app)
	if isID != nil {
		var err error
		app, err = github.GetApp(app)
		if err != nil {
			errored("invalid app\n", err)
		}
	}

	key := os.Args[2]
	if key == "" {
		errored("base64 private key missing", nil)
	}

	owner := ""
	if len(os.Args) >= 4 {
		owner = os.Args[3]
	}
	if owner == "" {
		owner = os.Getenv("GITHUB_REPOSITORY_OWNER")
		if owner == "" {
			errored("org name missing", nil)
		}
	}

	pem, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		errored("decoding base64 key", err)
	}
	jwt, err := generateJWT(app, pem)
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

func generateJWT(app string, pem []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		return "", err
	}

	now := time.Now().Unix()
	claims := jwt.StandardClaims{
		ExpiresAt: now + 60,
		IssuedAt:  now - 60,
		Issuer:    app,
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
