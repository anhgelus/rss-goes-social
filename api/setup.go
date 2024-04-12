package api

import (
	"encoding/json"
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
	println("Now you have to click on the link below and login to your account. It will answer you if you " +
		"really want to add this application. Click on 'Allow' and copy the token.")
	fmt.Printf(utils.GetFullUrl(
		url,
		fmt.Sprintf(
			"/oauth/authorize?%s=client_id&redirect_uri=urn:ietf:wg:oauth:2.0:oob&response_type=code",
			appCreated.ClientID),
	) + "\n")
	var tmpToken string
	_, err = fmt.Scanln(&tmpToken)
	if err != nil {
		println("Error", "client_id: "+appCreated.ClientID, " and client_secret: "+appCreated.ClientSecret)
		panic(err)
	}
	b, err = json.Marshal(setupRequestToken{
		RedirectUris: RedirectUris,
		ClientID:     appCreated.ClientID,
		ClientSecret: appCreated.ClientSecret,
		GrantType:    "authorization_code",
		Code:         tmpToken,
	})
	if err != nil {
		println("Error", "client_id: "+appCreated.ClientID, " and client_secret: "+appCreated.ClientSecret)
		panic(err)
	}
	req, err = newRequestNoToken(http.MethodPost, utils.GetFullUrl(url, "/oauth/token"), strings.NewReader(string(b)))
	if err != nil {
		println("Error", "client_id: "+appCreated.ClientID, " and client_secret: "+appCreated.ClientSecret)
		panic(err)
	}
	res, err = doRequest(req)
	if err != nil {
		println("Error", "client_id: "+appCreated.ClientID, " and client_secret: "+appCreated.ClientSecret)
		panic(err)
	}
	b, err = io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var tokenGot setupTokenGot
	err = json.Unmarshal(b, &appCreated)
	if err != nil {
		println("Error", "client_id: "+appCreated.ClientID, " and client_secret: "+appCreated.ClientSecret)
		panic(err)
	}
	err = VerifyToken(url, tokenGot.AccessToken)
	if err != nil {
		println(
			"Error while verifying the validity of the token",
			"client_id: "+appCreated.ClientID, " client_secret: "+appCreated.ClientSecret,
			" token: "+tokenGot.AccessToken,
		)
		panic(err)
	}
	fmt.Printf("These information are important and sensible. Keep them secret and safe in the right place!\n"+
		"Your token is: %s\nYour client_id is: %s\nYour client_secret is: %s",
		tokenGot.AccessToken,
		appCreated.ClientID,
		appCreated.ClientSecret,
	)
}
