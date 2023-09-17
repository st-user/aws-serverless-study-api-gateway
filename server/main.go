package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-chi/chi"
)

func main() {
	startHttpServer("8080")
}

type TokenResponse struct {
	Region         string `json:"region"`
	IdentityPoolID string `json:"identityPoolId"`
	IdentityID     string `json:"identityId"`
	LoginProvider  string `json:"loginProvider"`
	Token          string `json:"token"`
}

type APIGatewayAccessInfoResponse struct {
	Token *TokenResponse `json:"token"`
	Url   string         `json:"url"`
}

func getToken(loginName string) (*TokenResponse, error) {

	identityPoolID := os.Getenv("AWS_IDENTITY_POOL_ID")
	loginProvider := os.Getenv("AWS_LOGIN_PROVIDER")
	tokenDurationStr := os.Getenv("AWS_TOKEN_DURATION_SECONDS")

	tokenDuration, err := strconv.Atoi(tokenDurationStr)

	if err != nil {
		return nil, fmt.Errorf("AWS_TOKEN_DURATION must be an integer: %w", err)
	}

	fmt.Println("identityPoolID: " + identityPoolID)
	fmt.Println("loginProvider: " + loginProvider)
	fmt.Println("loginName: " + loginName)
	fmt.Println("tokenDuration: " + tokenDurationStr)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}
	svc := cognitoidentity.NewFromConfig(cfg)

	resp, err := svc.GetOpenIdTokenForDeveloperIdentity(context.TODO(), &cognitoidentity.GetOpenIdTokenForDeveloperIdentityInput{
		IdentityPoolId: aws.String(identityPoolID),
		Logins: map[string]string{
			loginProvider: loginName,
		},
		TokenDuration: aws.Int64(int64(tokenDuration)),
	})

	if err != nil {
		return nil, fmt.Errorf("unable to get token: %w", err)
	}

	identityID := resp.IdentityId
	token := resp.Token

	return &TokenResponse{
		Region:         cfg.Region,
		IdentityPoolID: identityPoolID,
		IdentityID:     *identityID,
		LoginProvider:  loginProvider,
		Token:          *token,
	}, nil
}

func startHttpServer(port string) {
	router := chi.NewRouter()
	router.Use(cors())

	router.HandleFunc("/accessInfo", func(w http.ResponseWriter, r *http.Request) {

		loginName := r.URL.Query().Get("userId")
		if loginName == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err := getToken(loginName)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		apiGatewayURL := os.Getenv("AWS_API_GATEWAY_URL")
		fmt.Println("AWS_API_GATEWAY_URL: " + apiGatewayURL)

		resp := APIGatewayAccessInfoResponse{
			Token: token,
			Url:   apiGatewayURL,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&resp)

		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func cors() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Access-Control-Allow-Origin", "*")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
