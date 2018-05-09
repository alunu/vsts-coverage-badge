package vsts

import "time"

// Build information about the build
type Build struct {
	ID          int32     `json:"id"`
	BuildNumber string    `json:"buildNumber"`
	FinishTime  time.Time `json:"finishTime"`
}

type buildList struct {
	Count int32   `json:"count"`
	Value []Build `json:"value"`
}

// ByFinishTimeDesc implements the sort interface for Build array
type ByFinishTimeDesc []Build

func (a ByFinishTimeDesc) Len() int           { return len(a) }
func (a ByFinishTimeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByFinishTimeDesc) Less(i, j int) bool { return a[i].FinishTime.After(a[j].FinishTime) }

// Timeline Information about the tasks executed
type Timeline struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Log  struct {
		ID int32 `json:"id"`
	} `json:"log"`
	Task struct {
		Name string `json:"name"`
	}
}

type timelineList struct {
	Records []Timeline `json:"records"`
}

// TestRun Information about a test run
type TestRun struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Build struct {
		ID string `json:"id"`
	} `json:"build"`
}

type testRunList struct {
	TestRuns []TestRun `json:"value"`
}

type testRunCount struct {
	Count int32 `json:"count"`
}

type codeCoverageResult struct {
	TestRun `json:"testRun"`
	Modules []struct {
		CoverageStatistic CodeCoverageStatistic `json:"statistics"`
	} `json:"modules"`
}

// CodeCoverageStatistic contains information about the number of lines/blocks covered by a test
type CodeCoverageStatistic struct {
	BlocksCovered         int `json:"blocksCovered"`
	BlocksNotCovered      int `json:"blocksNotCovered"`
	LinesCovered          int `json:"linesCovered"`
	LinesNotCovered       int `json:"linesNotCovered"`
	LinesPartiallyCovered int `json:"linesPartiallyCovered"`
}

type codeCoverageList struct {
	Value []codeCoverageResult `json:"value"`
	Count int32                `json:"count"`
}
