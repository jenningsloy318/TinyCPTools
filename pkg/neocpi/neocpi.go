package neocpi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type neoCPIResp struct {
	D struct {
		Results []struct {
			Metadata struct {
				ID   string `json:"id"`
				URI  string `json:"uri"`
				Type string `json:"type"`
			} `json:"__metadata"`
			Technicalname        string      `json:"TechnicalName"`
			Displayname          string      `json:"DisplayName"`
			Shorttext            string      `json:"ShortText"`
			RegID                string      `json:"reg_id"`
			Featured             interface{} `json:"Featured"`
			Scope                interface{} `json:"Scope"`
			Description          string      `json:"Description"`
			Version              interface{} `json:"Version"`
			Category             string      `json:"Category"`
			Mode                 string      `json:"Mode"`
			Vendor               string      `json:"Vendor"`
			Orgname              interface{} `json:"OrgName"`
			Supportedplatforms   string      `json:"SupportedPlatforms"`
			Products             interface{} `json:"Products"`
			Industries           interface{} `json:"Industries"`
			Lineofbusiness       interface{} `json:"LineOfBusiness"`
			Keywords             interface{} `json:"Keywords"`
			Countries            interface{} `json:"Countries"`
			Avgrating            string      `json:"AvgRating"`
			Ratingcount          string      `json:"RatingCount"`
			Publishedat          interface{} `json:"PublishedAt"`
			Publishedby          interface{} `json:"PublishedBy"`
			Createdat            string      `json:"CreatedAt"`
			Createdby            string      `json:"CreatedBy"`
			Modifiedat           string      `json:"ModifiedAt"`
			Modifiedby           string      `json:"ModifiedBy"`
			Partnercontent       interface{} `json:"PartnerContent"`
			Certifiedbysap       interface{} `json:"CertifiedBySap"`
			Additionalattributes interface{} `json:"AdditionalAttributes"`
			Artifacts            struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Artifacts"`
			Files struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Files"`
			Urls struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"Urls"`
			Medialinks struct {
				Deferred struct {
					URI string `json:"uri"`
				} `json:"__deferred"`
			} `json:"MediaLinks"`
		} `json:"results"`
	} `json:"d"`
}

type NeoCPIClient struct {
	httpClient *http.Client
	userName   string
	passWord   string
	tenantURL  string
}

func NewNeoCPIClient(username string, password string, tenantURL string) *NeoCPIClient {

	return &NeoCPIClient{
		httpClient: &http.Client{},
		userName:   username,
		passWord:   password,
		tenantURL:  tenantURL,
	}
}

func (nclient *NeoCPIClient) GetPkgRegIDList() map[string]string {

	apiURL := fmt.Sprintf("%s/itspaces/odata/1.0/workspace.svc/ContentEntities.ContentPackages?$format=json", nclient.tenantURL)

	httpReqest, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Failed to create request")
	}

	httpReqest.SetBasicAuth(nclient.userName, nclient.passWord)

	httpResp, err := nclient.httpClient.Do(httpReqest)

	if err != nil {
		fmt.Println("Failed to get response")
	}
	httpResBody := httpResp.Body

	if httpResBody != nil {
		defer httpResBody.Close()
	}

	bodyBytes, err := ioutil.ReadAll(httpResBody)

	if err != nil {
		fmt.Printf("Failed to get response body bytes from response,%v", err)
	}
	var neoCPIRespJson neoCPIResp

	err = json.Unmarshal(bodyBytes, &neoCPIRespJson)
	if err != nil {
		fmt.Printf("Failed to unmarshal json data from response bytes,%v", err)
	}
	var regIDs = map[string]string{}

	for _, pkginfo := range neoCPIRespJson.D.Results {
		regIDs[pkginfo.Technicalname] = pkginfo.RegID
	}

	return regIDs

}

func (nclient *NeoCPIClient) ExportPkg(pkg_reg_id string, filePath string) error {
	apiURL := fmt.Sprintf("https://%s//itspaces/api/1.0/workspace/%s?export=true", nclient.tenantURL, pkg_reg_id)

	httpReqest, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Failed to create request")
	}
	httpReqest.SetBasicAuth(nclient.userName, nclient.passWord)

	httpResp, err := nclient.httpClient.Do(httpReqest)

	if err != nil {
		fmt.Println("Failed to get response")
		return err
	}
	httpResBody := httpResp.Body

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer outFile.Close()
	_, err = io.Copy(outFile, httpResBody)
	return err

}
