package cfcpi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CFCPIClient struct {
	client        *http.Client
	cfUaaTokenURL string
	clientID      string
	clientSecret  string
	accessToken   string
	cfCpiAPI      string
}
type authResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}

func NewCFClient(clientID string, clientSecret string, cfUaaTokenURL string) *CFCPIClient {
	return &CFCPIClient{
		cfUaaTokenURL: cfUaaTokenURL,
		clientID:      clientID,
		clientSecret:  clientSecret,
	}

}
func (cfCPIClient *CFCPIClient) SetCpiAPI(cfCpiAPI string) {
	cfCPIClient.cfCpiAPI = strings.TrimRight(cfCpiAPI, "/")
}

func (cfCPIClient *CFCPIClient) GetAccessToken() {

	httpClient := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: 10 * time.Second,
	}

	params := url.Values{}
	params.Add("client_id", cfCPIClient.clientID)
	params.Add("client_secret", cfCPIClient.clientSecret)
	params.Add("grant_type", "client_credentials")
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", cfCPIClient.cfUaaTokenURL, body)
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
	cfCPIClient.accessToken = authRespJson.AccessToken
	cfCPIClient.client = httpClient
}

type IntegrationPackagesResp struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID          string `json:"id"`
				URI         string `json:"uri"`
				Type        string `json:"type"`
				ContentType string `json:"content_type"`
				MediaSrc    string `json:"media_src"`
				EditMedia   string `json:"edit_media"`
			} `json:"__metadata"`
			ID                             string      `json:"Id"`
			Name                           string      `json:"Name"`
			Description                    string      `json:"Description"`
			Shorttext                      string      `json:"ShortText"`
			Version                        string      `json:"Version"`
			Vendor                         string      `json:"Vendor"`
			Partnercontent                 bool        `json:"PartnerContent"`
			Updateavailable                bool        `json:"UpdateAvailable"`
			Mode                           string      `json:"Mode"`
			Supportedplatform              string      `json:"SupportedPlatform"`
			Modifiedby                     string      `json:"ModifiedBy"`
			Creationdate                   string      `json:"CreationDate"`
			Modifieddate                   string      `json:"ModifiedDate"`
			Createdby                      string      `json:"CreatedBy"`
			Products                       string      `json:"Products"`
			Keywords                       string      `json:"Keywords"`
			Countries                      string      `json:"Countries"`
			Industries                     string      `json:"Industries"`
			Lineofbusiness                 string      `json:"LineOfBusiness"`
			Packagecontent                 interface{} `json:"PackageContent"`
			Integrationdesigntimeartifacts struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"IntegrationDesigntimeArtifacts"`
			Valuemappingdesigntimeartifacts struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"ValueMappingDesigntimeArtifacts"`
			Customtags struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"CustomTags"`
		} `json:"results"`
	} `json:"d"`
}

func (cfCPIClient *CFCPIClient) GetIntegrationPackages() map[string]string {

	cfCpiIntPkgsURL := fmt.Sprintf("%s/IntegrationPackages", cfCPIClient.cfCpiAPI)
	httpReqest, err := http.NewRequest("GET", cfCpiIntPkgsURL, nil)

	if err != nil {
		fmt.Printf("Error occurs when creating http request to all packages,%v", err)
	}

	httpReqest.Header.Set("Accept", "application/json")

	bearerToken := fmt.Sprintf("Bearer %s", cfCPIClient.accessToken)

	httpReqest.Header.Set("Authorization", bearerToken)

	httpClient := cfCPIClient.client
	httpResp, err := httpClient.Do(httpReqest)

	if err != nil {
		fmt.Printf("Failed to get response,%v", err)
	}
	httpResBody := httpResp.Body

	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}
	bodyBytes, err := ioutil.ReadAll(httpResBody)
	statusCode := httpResp.StatusCode
	if statusCode != 200 {
		fmt.Errorf("Getting workspace info failed,response boy: %s\n", bodyBytes)
		return map[string]string{}
	}

	if err != nil {
		fmt.Printf("Failed to retrieve content bytes  data from response ,%v", err)
	}

	var cfIntPackgesResp IntegrationPackagesResp

	err = json.Unmarshal(bodyBytes, &cfIntPackgesResp)
	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}

	var packageList = map[string]string{}
	for _, packageInfo := range cfIntPackgesResp.D.Results {

		packageList[packageInfo.Name] = packageInfo.ID
	}
	return packageList
}

type PackageInfoResp struct {
	D struct {
		Metadata struct {
			ID          string `json:"id"`
			URI         string `json:"uri"`
			Type        string `json:"type"`
			ContentType string `json:"content_type"`
			MediaSrc    string `json:"media_src"`
			EditMedia   string `json:"edit_media"`
		} `json:"__metadata"`
		ID              string      `json:"Id"`
		Version         string      `json:"Version"`
		Packageid       string      `json:"PackageId"`
		Name            string      `json:"Name"`
		Description     string      `json:"Description"`
		Sender          string      `json:"Sender"`
		Receiver        string      `json:"Receiver"`
		Artifactcontent interface{} `json:"ArtifactContent"`
		Configurations  struct {
			Deferred struct {
				URI string `json:"uri"`
			} `json:"__deferred"`
		} `json:"Configurations"`
		Resources struct {
			Deferred struct {
				URI string `json:"uri"`
			} `json:"__deferred"`
		} `json:"Resources"`
	} `json:"d"`
}

func (cfCPIClient *CFCPIClient) GetIntegrationPackage(packageID string) {

	cfCpiIntPkgURL := fmt.Sprintf("%s/IntegrationPackages('%s')", cfCPIClient.cfCpiAPI, packageID)
	httpReqest, err := http.NewRequest("GET", cfCpiIntPkgURL, nil)

	if err != nil {
		fmt.Printf("Error occurs when creating http request to all packages,%v", err)
	}

	httpReqest.Header.Set("Accept", "application/json")

	bearerToken := fmt.Sprintf("Bearer %s", cfCPIClient.accessToken)

	httpReqest.Header.Set("Authorization", bearerToken)

	httpClient := cfCPIClient.client
	httpResp, err := httpClient.Do(httpReqest)

	if err != nil {
		fmt.Printf("Failed to get response,%v", err)
	}
	httpResBody := httpResp.Body

	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}
	bodyBytes, err := ioutil.ReadAll(httpResBody)
	statusCode := httpResp.StatusCode

	if statusCode != 200 {
		fmt.Printf("Getting package info failed,status code %d,response boy: %s\n", statusCode, bodyBytes)
		return
	}
	if err != nil {
		fmt.Printf("Failed to retrieve content bytes  data from response ,%v", err)
	}

	var packageInfoResp PackageInfoResp

	err = json.Unmarshal(bodyBytes, &packageInfoResp)
	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}

	fmt.Printf("Detail info about package %s and its version %s \n", packageID, packageInfoResp.D.Version)

}

type PackageIflowsResp struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID          string `json:"id"`
				URI         string `json:"uri"`
				Type        string `json:"type"`
				ContentType string `json:"content_type"`
				MediaSrc    string `json:"media_src"`
				EditMedia   string `json:"edit_media"`
			} `json:"__metadata"`
			ID              string      `json:"Id"`
			Version         string      `json:"Version"`
			Packageid       string      `json:"PackageId"`
			Name            string      `json:"Name"`
			Description     string      `json:"Description"`
			Sender          string      `json:"Sender"`
			Receiver        string      `json:"Receiver"`
			Artifactcontent interface{} `json:"ArtifactContent"`
			Configurations  struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Configurations"`
			Resources struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Resources"`
		} `json:"results"`
	} `json:"d"`
}

func (cfCPIClient *CFCPIClient) GetIntegrationPackageIflows(packageID string) map[string]string {
	cfCpiIntPkgURL := fmt.Sprintf("%s/IntegrationPackages('%s')/IntegrationDesigntimeArtifacts", cfCPIClient.cfCpiAPI, packageID)
	httpReqest, err := http.NewRequest("GET", cfCpiIntPkgURL, nil)

	if err != nil {
		fmt.Printf("Error occurs when creating http request to get all iflows in package %s ,%v", packageID, err)
	}

	httpReqest.Header.Set("Accept", "application/json")

	bearerToken := fmt.Sprintf("Bearer %s", cfCPIClient.accessToken)

	httpReqest.Header.Set("Authorization", bearerToken)

	httpClient := cfCPIClient.client
	httpResp, err := httpClient.Do(httpReqest)

	if err != nil {
		fmt.Printf("Failed to get response,%v", err)
	}
	httpResBody := httpResp.Body

	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}
	bodyBytes, err := ioutil.ReadAll(httpResBody)
	statusCode := httpResp.StatusCode

	if statusCode != 200 {
		fmt.Printf("Getting package info failed,status code %d,response boy: %s\n", statusCode, bodyBytes)
		return map[string]string{}
	}
	if err != nil {
		fmt.Printf("Failed to retrieve content bytes  data from response ,%v", err)
	}

	var packageIflowsResp PackageIflowsResp

	err = json.Unmarshal(bodyBytes, &packageIflowsResp)
	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}

	var iflowInfo = map[string]string{}

	for _, iflow := range packageIflowsResp.D.Results {
		iflowInfo[iflow.ID] = iflow.Version
	}

	fmt.Printf("iflown list %v\n", iflowInfo)
	return iflowInfo

}
