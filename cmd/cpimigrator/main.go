package main

import (
	"fmt"

	"github.com/jenningsloy318/TinyCPTools/pkg/cfcpi"
	"github.com/jenningsloy318/TinyCPTools/pkg/neocpi"
)

func main() {

	cfClientID := ""
	cfCientSecret := ""
	cfUaaTokenURL := ""

	fmt.Printf("Starting to get access token from  %s\n", cfUaaTokenURL)
	newCFClient := cfcpi.NewCFClient(cfClientID, cfCientSecret, cfUaaTokenURL)
	newCFClient.GetAccessToken()

	cfCpiAPI := ""
	fmt.Printf("Starting to get cf cpi pcakge list  %s\n", cfCpiAPI)
	newCFClient.SetCpiAPI(cfCpiAPI)
	packages := newCFClient.GetIntegrationPackages()
	for _, id := range packages {
		fmt.Printf("Starting to get package %s info  \n", id)

		newCFClient.GetIntegrationPackage(id)

		fmt.Printf("Starting to get all iflows inside package %s\n", id)
		newCFClient.GetIntegrationPackageIflows(id)
	}

		neoUser := ""
		neoPasswd := ""
		neoCPIURL := ""
		fmt.Printf("Starting to get package reg id  from tenant %s\n", neoCPIURL)
	
		newNeoClient := neocpi.NewNeoCPIClient(neoUser, neoPasswd, neoCPIURL)
	
		pkgRegIDList := newNeoClient.GetIntegrationPackageRegIDList()
		fmt.Printf("package reg id list: %v\n", pkgRegIDList)
	
		for name, id := range pkgRegIDList {
			filePath := id + ".zip"
			fmt.Printf("Starting to proceed package %s with reg_id %s\n", name, id)
			err := newNeoClient.ExportIntegrationPackage(id, filePath)
			if err != nil {
				fmt.Printf("error export the package %s, err is %s\n", name, err)
			}
		}
	
}
