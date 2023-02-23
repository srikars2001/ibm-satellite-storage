package mustachejsonvalidation

import (
	"fmt"
	"log"
	"os"
)

type validator struct {
	filePath   string
	isDebug    bool
	driverName string
	version    string
}

type validationPairs struct {
	jsonPath     string
	templatePath string
}

// constructor
func Construct(filePath string, driverName string, version string) *validator {
	return &validator{
		filePath:   filePath,
		isDebug:    false,
		driverName: driverName,
		version:    version,
	}
}

// starts validation process calls handleTemplateJSONValidation for each pair
func (v *validator) ValidateFiles() {
	// validate files for each combination
	defer printReport()

	filePairs := v.checkFilesForValidation()
	fmt.Println("filepairs : ", filePairs)
	fmt.Println()

	for _, pair := range filePairs {
		v.HandleTemplateJSONValidation(pair.jsonPath, pair.templatePath)
	}

}

// responsible for validating the given files
func (v *validator) HandleTemplateJSONValidation(jsonPath string, templatePath string) {
	templateValues := v.readJSONAndShowBrackets(templatePath)
	jsonValues := v.handleCustomParamsJSON(jsonPath)

	fmt.Println("templateValues : ", templateValues)
	fmt.Println("\n\njsonValues : ", jsonValues)

	v.CheckFaultsJSONYAML(templateValues, jsonValues)

}

// check which files to validate
func (v *validator) checkFilesForValidation() []validationPairs {
	var filePairs []validationPairs

	mp := make(map[string]struct{})

	entries, err := os.ReadDir(v.filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		mp[e.Name()] = struct{}{}
	}

	_, ok := mp["deployment.yaml"]

	if !ok {
		log.Fatal("deployment.yaml not found in path -> ", v.filePath)
	} else {

		_, ok := mp["custom-parameters.json"]

		if !ok {
			log.Fatal("custom-parameters.json not found in path -> ", v.filePath)
		} else {
			temp := validationPairs{
				jsonPath:     v.filePath + "custom-parameters.json",
				templatePath: v.filePath + "deployment.yaml",
			}

			filePairs = append(filePairs, temp)
		}
	}

	_, ok = mp["storage-class-template.yaml"]

	if ok {
		_, ok = mp["storage-class-parameters.json"]
		if !ok {
			log.Fatal("storage-class-parameters.json not found in path -> ", v.filePath, " but storage-class-template.yaml file is present")
		} else {
			temp := validationPairs{
				jsonPath:     v.filePath + "storage-class-parameters.json",
				templatePath: v.filePath + "storage-class-template.yaml",
			}

			filePairs = append(filePairs, temp)
		}

	}

	return filePairs
}
