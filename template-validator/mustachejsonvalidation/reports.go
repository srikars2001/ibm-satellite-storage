package mustachejsonvalidation

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Warnings struct {
	Parameter string `json:"parameter"`
	Filename  string `json:"filename"`
	Filepath  string `json:"filepath"`
	Message   string `json:"message"`
}

type ErrorsStruct struct {
	Parameter string `json:"parameter"`
	Filename  string `json:"filename"`
	Filepath  string `json:"filepath"`
	Message   string `json:"message"`
}

type OutputJSON struct {
	DriverName    string         `json:"driverName"`
	DriverVersion string         `json:"driverVersion"`
	Warnings      []Warnings     `json:"warnings"`
	Errors        []ErrorsStruct `json:"errors"`
}

var errorsArray []ErrorsStruct
var warningsArray []Warnings

func printReport() {
	fmt.Println(errorsArray)
	fmt.Println(warningsArray)

	b, err := json.MarshalIndent(errorsArray, "", "\t")
	if err != nil {
		log.Println("unable to convert struct to json ")
	}
	os.Stdout.Write(b)

	b, err = json.MarshalIndent(warningsArray, "", "\t")
	if err != nil {
		log.Println("unable to convert struct to json ")
	}
	os.Stdout.Write(b)

}

func (v *validator) SaveReport() {
	opJson := OutputJSON{
		DriverName:    v.driverName,
		DriverVersion: v.version,
		Warnings:      warningsArray,
		Errors:        errorsArray,
	}

	b, err := json.MarshalIndent(opJson, "", "\t")
	if err != nil {
		log.Println("unable to convert struct to json ")
	}
	fmt.Println("\n\nERRORS AND WARNINGS")
	os.Stdout.Write(b)

	outputFilePath := "./RESULTS/" + v.driverName + "_" + v.version + ".json"

	f, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Results are written in  ", outputFilePath)
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

// func CreateErrorsStruct(parameter string, filename string, message string, filepath string) ErrorsStruct {
// 	return ErrorsStruct{
// 		Parameter: parameter,
// 		Filename:  filename,
// 		Filepath:  filepath,
// 		Message:   message,
// 	}
// }
