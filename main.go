package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"vsts-coverage-badge/awsfunctions"
	"vsts-coverage-badge/vsts"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(Handler)
	//http.HandleFunc("/", localHTTPHandler)
	//http.ListenAndServe("localhost:8080", nil)
}

func localHTTPHandler(w http.ResponseWriter, r *http.Request) {
	svg, err := generateSvgText(r.URL.Query().Get("string"))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Add("Content-Type", "image/svg+xml")
	w.Write([]byte(svg))
}

func returnError(err error) (events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Content-Type"] = "text/plain"
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

	svgText, err := generateSvgText(percentageText)
	if err != nil {
		return nil, err
	}

	return []byte(svgText), nil
}

func generateSvgText(text string) (string, error) {
	// could do some math to work out proper length based on character
	// but for my use case right now theres only 4 cases, so just do them
	width := 94
	textLength := 250
	x := 755
	if len(text) == 2 { //single digit
		textLength = 195
	} else if len(text) == 4 {
		width = 100
		textLength = 310
		x = 790
	} else if text == "unknown" {
		width = 118
		textLength = 480
		x = 880
	} else if len(text) != 3 {
		return "", fmt.Errorf("Text %s was not recognized", text)
	}
	h := width - 59

	svg, err := ioutil.ReadFile("badge.svg")
	if err != nil {
		return "", fmt.Errorf("Could not read badge.svg, %s", err.Error())
	}

	svgText := fmt.Sprintf(string(svg), width, h, x, textLength, text)
	return svgText, nil
}

func calculateCodeCoverage() (float32, error) {
	builds, err := vsts.GetBuilds()
	if err != nil {
		return -1, err
	}
	if len(builds) == 0 {
		return -1, fmt.Errorf("no builds found")
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
	if notCovered+covered == 0 {
		return -1, fmt.Errorf("denominator was 0")
	}
	return (float32(covered) / float32(notCovered+covered)) * 100, nil
}
