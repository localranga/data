package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/localranga/tttt/request"
)

type DataPipeline struct {
	apiFile string
	input   <-chan *request.ResponseData
}

func NewDataPipeline(apiFile string, input <-chan *request.ResponseData) *DataPipeline {
	return &DataPipeline{
		apiFile: apiFile,
		input:   input,
	}
}

func (d *DataPipeline) Start() {
	for responseData := range d.input {
		formattedData, err := formatData(d.apiFile, responseData.Data)
		if err != nil {
			fmt.Printf("Error formatting data for %s: %v\n", d.apiFile, err)
			continue
		}

		apiID := strings.TrimSuffix(d.apiFile, ".go")
		err = saveResponseDataToFile(formattedData, apiID)
		if err != nil {
			fmt.Printf("Error saving data for %s: %v\n", apiID, err)
		}
	}
}

func formatData(api API, data []byte) (FormattedResponseData, error) {
	// Execute the API file
	apiResponseData, err := request.ExecuteAPIFile(data)
	if err != nil {
		return FormattedResponseData{}, err
	}

	// Format the API response data
	formattedData, err := formatAPIResponse(apiResponseData)
	if err != nil {
		return FormattedResponseData{}, err
	}

	return formattedData, nil
}

func saveResponseDataToFile(data FormattedResponseData, apiID string) error {
	// Create the output directory for the API if it doesn't exist
	err := os.MkdirAll(apiID, 0755)
	if err != nil {
		return err
	}

	// Generate a timestamped file name
	timestamp := time.Now().Format("20060102150405")
	fileName := fmt.Sprintf("%s/response_%s.json", apiID, timestamp)

	// Save the response data to a JSON file
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
