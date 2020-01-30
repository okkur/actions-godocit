package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type PR struct {
	Labels []Label `json:"labels"`
	User   User    `json:"user"`
	Title  string  `json:"title"`
	Body   string  `json:"body"`
	ID     int     `json:"number"`
	jwt    string
}

type Label struct {
	Name string `json:"name"`
}

type User struct {
	Login string `json:"login"`
}

type Installation struct {
	Token string `json:"token"`
}

func main() {
	// Get PR ID
	// GITHUB_REF format: refs/pull/:prNumber/merge
	PRcodeString := strings.Split(os.Getenv("GITHUB_REF"), "/")[2]
	PRcode, err := strconv.Atoi(PRcodeString)
	if err != nil {
		log.Fatalf("[godocit]: Couldn't convert PR id to int: %s", err.Error())
	}

	PRResp, err := http.Get(fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d", os.Getenv("GITHUB_REPOSITORY"), PRcode))
	if err != nil {
		log.Fatalf("[godocit]: Couldn't fetch PR's labels: %s", err.Error())
	}

	PRBody, err := ioutil.ReadAll(PRResp.Body)
	if err != nil {
		log.Fatalf("[godocit]: Couldn't read PR response's body: %s", err.Error())
	}

	var PR PR
	if err := json.Unmarshal(PRBody, &PR); err != nil {
		log.Fatalf("[godocit]: Couldn't parse the json body: %s", err.Error())
	}

	for _, label := range PR.Labels {
		if label.Name == "needs documentation" {
			jwt, err := GenerateJWT()
			if err != nil {
				log.Fatalf("[godocit]: Couldn't generate the JWT token: %s", err.Error())
			}

			token, err := GetToken(*jwt)
			if err != nil {
				log.Fatalf("[godocit]: Couldn't get the access token: %s", err.Error())
			}

			if err := CreateIssue(os.Getenv("INPUT_TARGETREPO"), PR, *token); err != nil {
				log.Fatalf("[godocit]: Couldn't create the documentation issue: %s", err.Error())
			}
		}
	}
}

func CreateIssue(repo string, PR PR, token string) error {
	var requestBody = []byte(fmt.Sprintf(`{
		"title": "[document]: %s",
		"body": "Write documentation for %s#%d"
	}`, PR.Title, os.Getenv("GITHUB_REPOSITORY"), PR.ID))

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/issues", repo), bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("Couldn't create the POST request for issue: %s", err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("Couldn't send the POST request to create issue: %s", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Couldn't create the issue. Received status code: %d", resp.StatusCode)
	}

	return nil
}

func GenerateJWT() (*string, error) {
	data, _ := pem.Decode([]byte(os.Getenv("PRIVATE_KEY")))
	parsedKey, err := x509.ParsePKCS1PrivateKey(data.Bytes)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Unix() + (10 * 60),
		"iss": 48309,
	})

	tokenString, err := token.SignedString(parsedKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't sign the JWT token: %s", err.Error())
	}

	return &tokenString, nil
}

func GetToken(jwt string) (*string, error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", os.Getenv("INSTALLATION_ID")),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Couldn't create the request for access token: %s", err.Error())
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't send the POST request to get the access token: %s", err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[godocit]: Couldn't read access token response's body: %s", err.Error())
	}

	var installation Installation
	if err := json.Unmarshal(body, &installation); err != nil {
		log.Fatalf("[godocit]: Couldn't parse the json body: %s", err.Error())
	}

	return &installation.Token, nil
}
