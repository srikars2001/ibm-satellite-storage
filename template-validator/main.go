package main

import (
	"flag"
	"fmt"
	musJsonValidation "templateValidator/mustachejsonvalidation"
	pathdiscovery "templateValidator/pathDiscovery"
)

var (
	providerName = flag.String("driverName", "", "driver_name")
	version      = flag.String("version", "", "version")
)

func main() {
	flag.Parse()
	fmt.Println("Template-validation started")

	fmt.Println("driverName: ", *providerName)
	fmt.Println("version: ", *version)

	filePath := pathdiscovery.GetFilesPath(*providerName, *version)

	tempValidator := musJsonValidation.Construct(filePath)

	tempValidator.ValidateFiles()

}
