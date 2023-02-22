package pathdiscovery

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// getFilesPath helps to locate the folder where files are present.
func GetFilesPath(driverName string, version string) (string, error) {
	providerName := getProviderName(driverName)
	fmt.Println("provider name : ", providerName)
	return "", nil
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
