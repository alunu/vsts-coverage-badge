package vsts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"vsts-coverage-badge/rest"
)

var projectName = url.PathEscape("")
var tenantName = url.PathEscape("")
var accessToken = ""

// GetBuilds gets a list of all the builds for the current project
func GetBuilds() []Build {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/build/builds?api-version=5.0-preview.4", tenantName, projectName)
	body := rest.Get(url, accessToken)
	var builds buildList
	if err := json.Unmarshal(body, &builds); err != nil {
		log.Fatal("Failed to unmarshal Build JSON, ", err)
	}

	log.Printf("Received %v builds", builds.Count)
	return builds.Value
}

// GetTimeline gets the timeline (the list of each task executed) for the build
func GetTimeline(build *Build) []Timeline {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/build/builds/%s/timeline?api-version=5.0-preview.2", tenantName, projectName, url.PathEscape(build.BuildNumber))
	body := rest.Get(url, accessToken)
	var timelines timelineList
	if err := json.Unmarshal(body, &timelines); err != nil {
		log.Fatal("Failed to unmarshal Timeline JSON, ", err)
	}

	return timelines.Records
}

// GetTestRuns Gets information regarding the test runs for the build
func GetTestRuns(build *Build) []TestRun {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/test/runs?api-version=5.0-preview.2&includeRunDetails=true", tenantName, projectName)
	body := rest.Get(url, accessToken)
	var testRuns testRunList
	if err := json.Unmarshal(body, &testRuns); err != nil {
		log.Fatal("Failed to unmarshal Test Run JSON ", err)
	}
	log.Printf("Received %v Test Runs", len(testRuns.TestRuns))
	return findBuildTestRuns(&testRuns.TestRuns, build.ID)
}

// GetCodeCoverageStatistics Gets the statistics (blocks/lines covered) for the specified build
func GetCodeCoverageStatistics(build *Build) []CodeCoverageStatistic {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/test/codecoverage?api-version=5.0-preview.1&flags=7&buildId=%v", tenantName, projectName, build.ID)
	body := rest.Get(url, accessToken)
	var codeCoverageList codeCoverageList
	if err := json.Unmarshal(body, &codeCoverageList); err != nil {
		log.Fatal("Failed to unmarshal Code Coverage JSON, ", err)
	}

	if codeCoverageList.Count == 0 {
		log.Fatal("No Code Coverage Results for build ", build.ID)
	}

	if codeCoverageList.Count > 1 {
		log.Fatal("Expected 1 Code Coverage result but found ", codeCoverageList.Count)
	}

	var ret []CodeCoverageStatistic
	for _, module := range codeCoverageList.Value[0].Modules {
		ret = append(ret, module.CoverageStatistic)
	}
	return ret
}

func findBuildTestRuns(testRuns *[]TestRun, buildID int32) []TestRun {
	var ret []TestRun
	buildIDString := fmt.Sprintf("%v", buildID)
	for _, element := range *testRuns {
		if element.Build.ID == buildIDString {
			ret = append(ret, element)
		}
	}
	log.Printf("Build %s has %v Test Run(s)", buildIDString, len(ret))
	return ret
}

func getTestRunCount(build *Build) int32 {
	url := fmt.Sprintf("https://%s.visualstudio.com/%s/_apis/test/runs?api-version=5.0-preview.2", tenantName, projectName)
	body := rest.Get(url, accessToken)
	var count testRunCount
	if err := json.Unmarshal(body, &count); err != nil {
		log.Fatal("Failed to unmarshal Test Run Count JSON ", err)
	}
	return count.Count
}
