package pathdiscovery

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// getFilesPath helps to locate the folder where files are present.
func GetFilesPath(driverName string, version string) string {
	providerName := getProviderName(driverName)
	// fmt.Println("provider name : ", providerName)
	tempFilePath := "./config-templates/" + providerName + "/" + driverName + "/"

	cv := checkVersion(tempFilePath, version)

	if !cv {
		log.Fatal("Version : ", version, " not found in driver-name :", driverName, " || provider :", providerName)
	}

	filePath := tempFilePath + version + "/"

	fmt.Println("file path = ", filePath)

	return filePath
}

// checkVersion takes path as argument and checks whether the version is present,
// the program will quit if the version is not present in the pathName folder
func checkVersion(pathName string, version string) bool {
	entries, err := os.ReadDir(pathName)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		if e.Name() == version {
			return true
		}
	}

	return false
}

// takes driverName as parameter and gets the provider name.
func getProviderName(driverName string) string {

	type templateJSON struct {
		Name     string `json:"name"`
		Provider string `json:"provider"`
	}

	file, err := os.Open("./config-templates/template_list.json")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)

	if err != nil {
		log.Panic(err)
	}

	var values []templateJSON

	err = json.Unmarshal(data, &values)

	if err != nil {
		log.Println(err)
	}

	for _, v := range values {
		if v.Name == driverName {
			return v.Provider
		}
	}

	log.Fatal(driverName, " ==> not found in template_list.json \n exiting the program")

	return ""

}
