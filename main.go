package main

import (
	"fmt"
	"log"
	"sort"
	"vsts-coverage-badge/vsts"
)

func main() {
	fmt.Printf("Code Coverage: %.2f%%", calculateCodeCoverage())
}

func calculateCodeCoverage() float32 {
	builds := vsts.GetBuilds()
	if len(builds) == 0 {
		log.Fatal("No builds found.")
	}

	sort.Sort(vsts.ByFinishTimeDesc(builds))
	testRuns := vsts.GetCodeCoverageStatistics(&builds[0])

	covered := 0
	notCovered := 0
	for _, stat := range testRuns {
		covered += stat.BlocksCovered
		notCovered += stat.BlocksNotCovered
	}
	return (float32(covered) / float32(notCovered+covered)) * 100
}
