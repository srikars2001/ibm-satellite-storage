package mustachejsonvalidation

import "fmt"

type Warnings struct {
	Parameter string `json:"paramete"`
	Filename  string `json:"filename"`
	Message   string `json:"message"`
}

type ErrorsStruct struct {
	Parameter string `json:"parameter"`
	Filename  string `json:"filename"`
	Message   string `json:"message"`
}

type OutputJSON struct {
	TemplateName    string         `json:"templateName"`
	TemplateVersion string         `json:"templateVersion"`
	Warnings        []Warnings     `json:"warnings"`
	Errors          []ErrorsStruct `json:"errors"`
}

var errorsArray []ErrorsStruct
var warningsArray []Warnings

func printReport() {
	fmt.Println(errorsArray)
	fmt.Println(warningsArray)
}
