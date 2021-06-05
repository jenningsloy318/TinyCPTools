package cfcpi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

type CFCPIResp []struct {
	ID                 string   `json:"id"`
	Modifiedby         string   `json:"modifiedBy"`
	Creationdate       string   `json:"creationDate"`
	Modifieddate       string   `json:"modifiedDate"`
	Createdby          string   `json:"createdBy"`
	Products           []string `json:"products"`
	Keywords           []string `json:"keywords"`
	Supportedplatforms []string `json:"supportedPlatforms"`
	Countries          []string `json:"countries"`
	Industries         []string `json:"industries"`
	Lineofbusiness     []string `json:"lineOfBusiness"`
	Vendor             string   `json:"vendor"`
	Version            string   `json:"version"`
	Description        string   `json:"description"`
	Privilegestate     string   `json:"privilegeState"`
	Technicalname      string   `json:"technicalName"`
	Shortdescription   string   `json:"shortDescription"`
	Partnercontent     bool     `json:"partnerContent"`
	Updateavailable    bool     `json:"updateAvailable"`
	Title              string   `json:"title"`
	Type               string   `json:"type"`
}

type CFCPIClient struct {
	httpClient   *http.Client
	uaaURL       string
	clientID     string
	clientSecret string
	userName     string
	passWord     string
	accessToken  string
	oauthClient  *http.Client
}

type authResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}

func NewCFClient(clientID string, clientSecret string, username string, password string, uaaURL string) *CFCPIClient {
	return &CFCPIClient{
		uaaURL:       uaaURL,
		clientID:     clientID,
		userName:     username,
		passWord:     password,
		clientSecret: clientSecret,
	}

}

func (cfCPIClient *CFCPIClient) GetAccessTokenHttp() {

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 10 * time.Second,
	}

	apiURL := fmt.Sprintf("%s/oauth/token", cfCPIClient.uaaURL)

	params := url.Values{}
	params.Add("client_id", cfCPIClient.clientID)
	params.Add("client_secret", cfCPIClient.clientSecret)
	params.Add("grant_type", "client_credentials")
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", apiURL, body)
	if err != nil {
		fmt.Printf("Error occurs when creating http request to fetch acess token,%v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpResp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error occurs when geting http response for acess token,%v", err)

	}

	httpResBody := httpResp.Body
	statusCode := httpResp.StatusCode
	if httpResBody != nil {
		defer httpResBody.Close()
	}

	if statusCode != 200 {
		fmt.Printf("Getting access code failed ")
	}

	bodyBytes, err := ioutil.ReadAll(httpResBody)
	if err != nil {
		fmt.Printf("Failed to get response body,%v", err)
	}

	var authRespJson authResp

	err = json.Unmarshal(bodyBytes, &authRespJson)

	if err != nil {
		fmt.Printf("Error occurs when decode response to json,%v", err)
	}
	token := authRespJson.AccessToken
	cfCPIClient.accessToken = token
	cfCPIClient.httpClient = httpClient
}
func (cfCPIClient *CFCPIClient) GetCFWorkspaceHttp(cfCPIWorkspaceURL string) {

	httpReqest, err := http.NewRequest("POST", cfCPIWorkspaceURL, nil)

	if err != nil {
		fmt.Printf("Error occurs when creating http request to get workspace infomation,%v", err)
	}

	httpReqest.Header.Set("Accept", "application/json")

	bearerToken := fmt.Sprintf("Bearer %s", cfCPIClient.accessToken)

	httpReqest.Header.Set("Authorization", bearerToken)
	fmt.Println(httpReqest.Header)
	httpClient := cfCPIClient.httpClient
	httpResp, err := httpClient.Do(httpReqest)

	if err != nil {
		fmt.Printf("Failed to get response,%v", err)
	}

	httpResBody := httpResp.Body
	statusCode := httpResp.StatusCode

	bodyBytes, err := ioutil.ReadAll(httpResBody)
	if err != nil {
		fmt.Printf("Failed to get response body,%v", err)
	}

	if statusCode != 200 {
		fmt.Printf("Getting workspace info failed,response boy: %s\n", bodyBytes)
	}

	httpResBody.Close()
}

func (cfCPIClient *CFCPIClient) GetAccessTokenOauth() {

	cfAuthconfig := &clientcredentials.Config{
		ClientID:     cfCPIClient.clientID,
		ClientSecret: cfCPIClient.clientSecret,
		TokenURL:     cfCPIClient.uaaURL + "/oauth/token",
	}

	ctx := context.Background()
	token, err := cfAuthconfig.Token(ctx)

	if err != nil {
		fmt.Printf("Error when generating token from client credentials, %v\n", err)
	}
	fmt.Printf("Access token %v\n", token.AccessToken)
	oauthClient := cfAuthconfig.Client(ctx)
	cfCPIClient.accessToken = token.AccessToken
	cfCPIClient.oauthClient = oauthClient
}
func (cfCPIClient *CFCPIClient) GetCFWorkspaceOauth(cfCPIWorkspaceURL string) {
	httpReqest, err := http.NewRequest("GET", cfCPIWorkspaceURL, nil)
	if err != nil {
		fmt.Printf("Error when creating http request, %v\n", err)
	}
	httpClient := cfCPIClient.oauthClient

	httpReqest.Header.Set("Accept", "application/json")

	bearerToken := fmt.Sprintf("Bearer %s", cfCPIClient.accessToken)

	httpReqest.Header.Set("Authorization", bearerToken)
	httpResp, err := httpClient.Do(httpReqest)

	if err != nil {
		fmt.Printf("Failed to get response,%v\n", err)
	}

	httpResBody := httpResp.Body
	statusCode := httpResp.StatusCode

	bodyBytes, err := ioutil.ReadAll(httpResBody)
	if err != nil {
		fmt.Printf("Failed to get response body,%v\n", err)
	}
	fmt.Printf("Getting workspace info failed,response body: %v\n", httpResBody)

	if statusCode != 200 {
		fmt.Printf("Getting workspace info failed,response body")
	}
	fmt.Printf("Response Body %v\n", string(bodyBytes))

	httpResBody.Close()
}
