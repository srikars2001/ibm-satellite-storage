package mustachejsonvalidation

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

var (
	statusOk = true // statusok is false when errors or warnings are produced
)

// reads json file and shows the line which has "{{"
func (v *validator) readJSONAndShowBrackets(templatePath string, filePairs string) []string {
	var lists []string

	file, err := os.Open(templatePath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "{{") || strings.Contains(line, "}}") {
			lists = append(lists, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	if v.isDebug {
		fmt.Println("=====================LINES WiTH BRACKETS====================\n", templatePath)
		for _, value := range lists {
			fmt.Println(value)
		}
		fmt.Println("=====================END====================")
	}

	var clean []string

	for _, line := range lists {
		//startingComment := strings.Index(line,"#")

		startingIndex := strings.Index(line, "{{")

		if startingIndex == -1 {
			log.Println("unable to clean --- '{{' missing in .yaml file: ", templatePath, line)
			ew := ErrorsStruct{
				Parameter: line,
				Filepath:  v.filePath,
				Message:   "brackets mismatch {{ " + "unable to clean --- '{{' missing in .yaml file line " + line,
				Filename:  filePairs,
			}

			errorsArray = append(errorsArray, ew)
			continue
		}

		lastIndex := strings.LastIndex(line, "}}")

		if lastIndex == -1 {
			log.Println("unable to clean --- '}}' missing in .yaml file: ", templatePath, line)

			ew := ErrorsStruct{
				Parameter: line,
				Filepath:  v.filePath,
				Message:   "brackets-mismatch }} " + "unable to clean --- '}}' missing in .yaml file line " + line,
				Filename:  filePairs,
			}

			errorsArray = append(errorsArray, ew)
			continue
		}

		tempString := line[startingIndex : lastIndex+2]

		clean = append(clean, tempString)

	}

	if v.isDebug {
		fmt.Println("=====================CLEAN  BRACKETS====================\n", templatePath)
		for _, value := range clean {
			fmt.Println(value)
		}
		fmt.Println("=====================END====================")
	}

	return clean
}

func (v *validator) handleCustomParamsJSON(jsonPath string, filePairs string) map[string]struct{} {
	type DefinedJSON struct {
		Name string `json:"name"`
	}

	l := make(map[string]struct{})
	file, err := os.Open(jsonPath)

	if err != nil {
		log.Println(err)
	}

	data, err := io.ReadAll(file)

	if err != nil {
		log.Println(err)
	}

	var values []DefinedJSON

	err = json.Unmarshal(data, &values)

	if err != nil {
		log.Println(err)
	}

	for _, j := range values {
		l[j.Name] = struct{}{}
	}

	if v.isDebug {
		fmt.Println("============================json values====================\n", jsonPath)
		for _, j := range values {
			fmt.Println(j.Name)
		}
		fmt.Println("============================END====================")
	}

	//fmt.Println(l)
	return l
}

func (v *validator) CheckFaultsJSONYAML(templateValues []string, jsonValues map[string]struct{}, filePairs string) {
	var stack Stack

	if v.isDebug {
		fmt.Println("==========================check fault json yaml============================")
	}
	for _, value := range templateValues {
		pk := v.parseAndInit(value, filePairs)

		if pk.prefix == "#" || pk.prefix == "^" {
			stack.Push(pk)
		} else if pk.prefix == "/" {
			temp, ok := stack.Peek()
			if !ok {
				log.Println(fmt.Sprintln("stack is empty ==> the key: ", pk.value, "has missing opening statement either '#' or '^' "))
				ew := ErrorsStruct{
					Parameter: pk.value,
					Filepath:  v.filePath,
					Message:   fmt.Sprintln("stack is empty ==> the key: ", pk.value, "has missing opening statement either '#' or '^' "),
					Filename:  filePairs,
				}

				errorsArray = append(errorsArray, ew)
			}

			if (temp.prefix == "#" || temp.prefix == "^") && (temp.value == pk.value) && (temp.suffix == pk.suffix) {
				stack.Pop()
			} else {
				log.Println(string(colorRed), "error in validating yaml mustache template \n=== ")
				log.Println("current mustache structure ", pk)
				log.Println("previous stack value", temp, string(colorReset))
				statusOk = false
				ew := ErrorsStruct{
					Parameter: pk.value,
					Filepath:  v.filePath,
					Message:   fmt.Sprintln("error in validating yaml mustache template \n===\n ", "current mustache structure : ", pk, "\n\n previous stack value", temp),
					Filename:  filePairs,
				}
				errorsArray = append(errorsArray, ew)
			}
		} else {
			_, ok := jsonValues[pk.value]
			if !ok {

				temp, ok := stack.Peek()
				if !ok || temp.prefix != "#" {
					log.Println(fmt.Sprintln(string(colorRed), "error", pk.value, "===> missing in json", string(colorReset)))
					statusOk = false
					ew := ErrorsStruct{
						Parameter: pk.value,
						Filepath:  v.filePath,
						Message:   pk.value + " == missing in json",
						Filename:  filePairs,
					}

					errorsArray = append(errorsArray, ew)
				} else {
					log.Println(fmt.Sprintln(string(colorYellow), "warning", pk.value, "===> missing in json", string(colorReset)))
					statusOk = false
					ew := Warnings{
						Parameter: pk.value,
						Filepath:  v.filePath,
						Message:   pk.value + " == missing in json" + " ,parent value : " + temp.value,
						Filename:  filePairs,
					}

					warningsArray = append(warningsArray, ew)
				}
			}
		}

		if v.isDebug {
			fmt.Println(pk)
		}
	}

	//check stack length and report errors

	if !stack.IsEmpty() {
		//error report
		log.Println(fmt.Sprintln(string(colorRed), "Error : ", "the values ", stack, " doesnt have closing brackets", string(colorReset)))

		ew := ErrorsStruct{
			Parameter: "multiple",
			Filepath:  v.filePath,
			Message:   fmt.Sprintln("Error : ", "the values ", stack, " doesnt have closing brackets"),
			Filename:  filePairs,
		}
		errorsArray = append(errorsArray, ew)

		statusOk = false

	}

	if statusOk {
		log.Println(string(colorGreen), "No errors found in ", filePairs, string(colorReset))
	} else {
		log.Println(string(colorRed), "errors or warning found in filepath", string(colorYellow), v.filePath, " in file pairs :", filePairs, string(colorReset))
	}

	if v.isDebug {
		fmt.Println("==========================end============================")
	}

}

func (v *validator) parseAndInit(value string, filePairs string) ParamKeys {
	pk := ParamKeys{
		prefix: "$$$",
		suffix: "$$$",
	}

	var numberOfbrackets int

	var buffer string

	for _, char := range value {
		if char == '{' {
			numberOfbrackets = numberOfbrackets + 1
		} else if char == '}' {
			numberOfbrackets = numberOfbrackets - 1
		} else if char == '#' || char == '^' || char == '/' {
			pk.prefix = string(char)
		} else if char == ' ' {
			continue
		} else if char == '?' {
			pk.suffix = string(char)
		} else {
			buffer = buffer + string(char)
		}
	}

	pk.value = buffer

	if numberOfbrackets != 0 {
		log.Println(string(colorRed), "brackets mismatch for the YAML mustache value :  ", value, string(colorReset))

		statusOk = false
		ew := ErrorsStruct{
			Parameter: pk.value,
			Filepath:  v.filePath,
			Message:   "brackets mismatch for the YAML mustache value :  " + value,
			Filename:  filePairs,
		}

		errorsArray = append(errorsArray, ew)
	}

	return pk
}
