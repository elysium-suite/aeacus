package main

import (
	"github.com/iamacarpet/go-win64api/shared"
    "fmt"

    wapi "github.com/iamacarpet/go-win64api"
)

func main() {
    // serviceUp(`AdobeARMservice`)
    fmt.Println(getLocalServiceStatus(`AdobeARMservice`))
    serviceUpAndEnabled(`AdobeARMservice`)
}

func getLocalServiceStatus(serviceName string) (shared.Service, error) {
    serviceDataList, err := wapi.GetServices()
    var serviceStatusData shared.Service
	if err != nil {
        fmt.Println("Couldn't get local service: " + err.Error())
        return serviceStatusData, err
    }
    for _, v := range serviceDataList {
        if v.SCName == serviceName {
            serviceStatusData = v
        }
    }
    return serviceStatusData, nil
}

func serviceUpAndEnabled(serviceName string) {
    serviceStatus, err := getLocalServiceStatus(serviceName)
    if err != nil {
        fmt.Println(`reeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee`)
    }
    
}