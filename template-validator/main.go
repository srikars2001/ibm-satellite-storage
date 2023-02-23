package main

import (
	"flag"
	"fmt"
	musJsonValidation "templateValidator/mustachejsonvalidation"
	pathdiscovery "templateValidator/pathDiscovery"
)

var (
	driverName = flag.String("driverName", "", "driver_name")
	version    = flag.String("version", "", "version")
)

func main() {
	flag.Parse()
	fmt.Println("Template-validation started")

	fmt.Println("driverName: ", *driverName)
	fmt.Println("version: ", *version)

	filePath := pathdiscovery.GetFilesPath(*driverName, *version)

	tempValidator := musJsonValidation.Construct(filePath, *driverName, *version)

	tempValidator.ValidateFiles()

}
