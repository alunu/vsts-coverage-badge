package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"vsts-coverage-badge/awsfunctions"
	"vsts-coverage-badge/vsts"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
}

func returnError(err error) (events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "plain/text"
	return events.APIGatewayProxyResponse{
		Body:            err.Error(),
		StatusCode:      500,
		IsBase64Encoded: false,
		Headers:         headers,
	}, nil
}

// Handler handles the request from API Gateway
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processig Lambda request %s", request.RequestContext.RequestID)

	vsts.TenantName = request.QueryStringParameters["Tenant"]
	vsts.ProjectName = request.QueryStringParameters["Project"]
	bucketName := request.QueryStringParameters["Bucket"]
	folderName := request.QueryStringParameters["Folder"]

	if vsts.TenantName == "" || vsts.ProjectName == "" {
		return returnError(fmt.Errorf("TenantName and ProjectName required"))
	}

	if folderName == "" {
		folderName = "badges"
	}

	svg, err := renderSvg()
	if err != nil {
		return returnError(fmt.Errorf("Failed to render svg, %v", err))
	}

	if bucketName != "" {
		err := awsfunctions.UploadToBucket(bucketName, fmt.Sprintf("%s/%s-%s-coverage.svg", folderName, vsts.TenantName, vsts.ProjectName), "image/svg+xml", &svg)
		if err != nil {
			return returnError(err)
		}
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "image/svg+xml"
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(svg),
	}

	return response, nil
}

func renderSvg() ([]byte, error) {
	percentageText := "unknown"
	coverageValue, err := calculateCodeCoverage()
	if err == nil {
		percentageText = fmt.Sprintf("%.f%%", coverageValue)
	}

	svg, err := ioutil.ReadFile("badge.svg")
	if err != nil {
		errorString := fmt.Sprintf("Could not read badge.svg, %s", err.Error())
		return nil, errors.New(errorString)
	}

	svgText := fmt.Sprintf(string(svg), percentageText)
	return []byte(svgText), nil
}

func calculateCodeCoverage() (float32, error) {
	builds, err := vsts.GetBuilds()
	if err != nil {
		return -1, err
	}
	if len(builds) == 0 {
		return -1, fmt.Errorf("No builds found.")
	}

	sort.Sort(vsts.ByFinishTimeDesc(builds))
	testRuns, err := vsts.GetCodeCoverageStatistics(&builds[0])
	if err != nil {
		return -1, err
	}

	covered := 0
	notCovered := 0
	for _, stat := range testRuns {
		covered += stat.BlocksCovered
		notCovered += stat.BlocksNotCovered
	}
	return (float32(covered) / float32(notCovered+covered)) * 100, nil
}
