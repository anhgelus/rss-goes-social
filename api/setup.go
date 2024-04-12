package api

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/anhgelus/rss-goes-social/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

type setupRegisterApplication struct {
	ClientName   string `json:"client_name"`
	RedirectUris string `json:"redirect_uris"`
	Scopes       string `json:"scopes"`
}

type setupApplicationCreated struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	RedirectUris string `json:"redirect_uris"`
	VapidKey     string `json:"vapid_key"`
	Website      string `json:"website"`
}

type setupRequestToken struct {
	RedirectUris string `json:"redirect_uris"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
}

type setupTokenGot struct {
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

const (
	AppName      = "rss_goes_social"
	RedirectUris = "urn:ietf:wg:oauth:2.0:oob"
)

func Setup() {
	if len(os.Args) < 3 {
		println("An argument is required to setup the application. Use " + os.Args[0] + " help setup for more information")
		return
	}
	url := os.Args[2]

	var clientID string
	var clientSecret string
	flag.StringVar(&clientID, "id", "", "client_id of the application (not required)")
	flag.StringVar(&clientSecret, "secret", "", "client_secret of the application (not required)")
	flag.Parse()

	if len(clientID) != 0 && len(clientSecret) == 0 {
		println("-id is set and not -secret. You have to set these two flags.")
		return
	} else if len(clientID) == 0 && len(clientSecret) != 0 {
		println("-secret is set and not -id. You have to set these two flags.")
		return
	} else if len(clientID) == 0 && len(clientSecret) == 0 {
		newApplication(url, &clientID, &clientSecret)
	}

	println("Now you have to click on the link below and login to your account. It will answer you if you " +
		"really want to add this application. Click on 'Allow' and copy the token.")
	fmt.Printf(utils.GetFullUrl(
		url,
		fmt.Sprintf(
			"/oauth/authorize?%s=client_id&redirect_uri=urn:ietf:wg:oauth:2.0:oob&response_type=code",
			clientID),
	) + "\n")
	print("The token obtained from the link above -> ")
	var tmpToken string
	_, err := fmt.Scanln(&tmpToken)
	if err != nil {
		println("Error", "client_id: "+clientID, " and client_secret: "+clientSecret)
		panic(err)
	}
	println("Getting the token.")
	b, err := json.Marshal(setupRequestToken{
		RedirectUris: RedirectUris,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "authorization_code",
		Code:         tmpToken,
	})
	if err != nil {
		println("Error", "client_id: "+clientID, " and client_secret: "+clientSecret)
		panic(err)
	}
	req, err := newRequestNoToken(http.MethodPost, utils.GetFullUrl(url, "/oauth/token"), strings.NewReader(string(b)))
	if err != nil {
		println("Error", "client_id: "+clientID, " and client_secret: "+clientSecret)
		panic(err)
	}
	res, err := doRequest(req)
	if err != nil {
		println("Error", "client_id: "+clientID, " and client_secret: "+clientSecret)
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var tokenGot setupTokenGot
	err = json.Unmarshal(b, &tokenGot)
	if err != nil {
		println("Error", "client_id: "+clientID, " and client_secret: "+clientSecret)
		panic(err)
	}
	println("Verifying the token.")
	err = VerifyToken(url, tokenGot.AccessToken)
	if err != nil {
		println(
			"Error while verifying the validity of the token",
			"client_id: "+clientID, " client_secret: "+clientSecret,
			" token: "+tokenGot.AccessToken,
		)
		panic(err)
	}
	fmt.Printf("These information are important and sensible. Keep them secret and safe in the right place!\n"+
		"Your token is: %s\nYour client_id is: %s\nYour client_secret is: %s",
		tokenGot.AccessToken,
		clientID,
		clientSecret,
	)
}

func newApplication(url string, clientID *string, clientSecret *string) {
	println("Registering the application.")
	b, err := json.Marshal(setupRegisterApplication{
		ClientName:   AppName,
		RedirectUris: RedirectUris,
		Scopes:       "write read",
	})
	if err != nil {
		panic(err)
	}
	req, err := newRequestNoToken(http.MethodPost, utils.GetFullUrl(url, "/api/v1/apps"), strings.NewReader(string(b)))
	if err != nil {
		panic(err)
	}
	res, err := doRequest(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Invalid status code after creating a new application: %d\n", res.StatusCode)
		return
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var appCreated setupApplicationCreated
	err = json.Unmarshal(b, &appCreated)
	if err != nil {
		panic(err)
	}
	*clientID = appCreated.ClientID
	*clientSecret = appCreated.ClientSecret
}
