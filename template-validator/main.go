package main

import (
	"flag"
	"fmt"
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

	pathdiscovery.GetFilesPath(*providerName, *version)

}
