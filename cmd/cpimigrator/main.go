package main

import (
	"fmt"

	"github.com/jenningsloy318/TinyCPTools/pkg/cfcpi"
	"github.com/jenningsloy318/TinyCPTools/pkg/neocpi"
)

func main() {

	cfClientID := 
	cfCientSecret := 
	cfUsername := 
	cfPassword := 
	cfUaaURL := 
	newCFClient := cfcpi.NewCFClient(cfClientID, cfCientSecret, cfUsername, cfPassword, cfUaaURL)
	fmt.Printf("Starting to get access token from  %s\n", cfUaaURL)

	newCFClient.GetAccessTokenOauth()
	cfCPIWorkspaceURL := 
	fmt.Printf("Starting to get cf cpi workspace info  %s\n", cfCPIWorkspaceURL)
	newCFClient.GetCFWorkspaceOauth(cfCPIWorkspaceURL)

	neoUser := 
	neoPasswd :=
	neoCPIURL := 
	fmt.Printf("Starting to get package reg id  from tenant %s\n", neoCPIURL)

	newNeoClient := neocpi.NewNeoCPIClient(neoUser, neoPasswd, neoCPIURL)

	pkgRegIDList := newNeoClient.GetPkgRegIDList()
	fmt.Printf("package reg id list: %v\n", pkgRegIDList)

	for name, id := range pkgRegIDList {
		filePath := id + ".zip"
		fmt.Printf("Starting to proceed package %s with reg_id %s\n", name, id)
		err := newNeoClient.ExportPkg(id, filePath)
		if err != nil {
			fmt.Printf("error export the package %s, err is %s\n", name, err)
		}
	}

}
