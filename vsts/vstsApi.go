package vsts

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"vsts-coverage-badge/rest"
)

// ProjectName is the name of the project in VSTS
var ProjectName = ""

//TenantName is the subdomain of visualstudio.com
var TenantName = ""

var accessToken = os.Getenv("VSTS_ACCESS_TOKEN")

// GetBuilds gets a list of all the builds for the current project
func GetBuilds() ([]Build, error) {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/build/builds?api-version=5.0-preview.4", TenantName, ProjectName)
	body := rest.Get(url, accessToken)
	var builds buildList
	if err := json.Unmarshal(body, &builds); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Build JSON, %v", err)
	}

	log.Printf("Received %v builds", builds.Count)
	return builds.Value, nil
}

// GetCodeCoverageStatistics Gets the statistics (blocks/lines covered) for the specified build
func GetCodeCoverageStatistics(build *Build) ([]CodeCoverageStatistic, error) {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/test/codecoverage?api-version=5.0-preview.1&flags=7&buildId=%v", TenantName, ProjectName, build.ID)
	body := rest.Get(url, accessToken)
	var codeCoverageList codeCoverageList
	if err := json.Unmarshal(body, &codeCoverageList); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Code Coverage JSON, %v", err)
	}

	if codeCoverageList.Count == 0 {
		return nil, fmt.Errorf("No Code Coverage Results for build %v", build.ID)
	}

	if codeCoverageList.Count > 1 {
		return nil, fmt.Errorf("Expected 1 Code Coverage result but found %v", codeCoverageList.Count)
	}

	var ret []CodeCoverageStatistic
	for _, module := range codeCoverageList.Value[0].Modules {
		ret = append(ret, module.CoverageStatistic)
	}
	return ret, nil
}
