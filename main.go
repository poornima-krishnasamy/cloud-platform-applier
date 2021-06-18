// Package main Declaration
package main

// Importing packages
import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/apply"
	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
)

// Main function
// env vars: REPO_PATH
func main() {

	fmt.Printf("START TIME %s \n", time.Now().String())

	fmt.Println("Version", runtime.Version())
	fmt.Println("NumCPU", runtime.NumCPU())
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	const nRoutines int = 3
	wg := &sync.WaitGroup{}

	repoPath := sysutil.GetRequiredEnvString("REPO_PATH")
	// clusterName := sysutil.GetRequiredEnvString("TF_VAR_cluster_name")
	// clusterStateBucket := sysutil.GetRequiredEnvString("TF_VAR_cluster_state_bucket")
	// clusterStateKey := sysutil.GetRequiredEnvString("TF_VAR_cluster_state_key")

	// clock := &sysutil.Clock{}

	runResults := make(chan apply.Results, 5)

	applier := &apply.Applier{RepoPath: repoPath, RunResults: runResults}

	applier.FullRun()

	// TODO Fix channel output

	go monitorResults(wg, runResults)
	go func() {
		for result := range runResults {
			log.Printf("Updating successes Run %v.", result.Successes)
			log.Printf("Updating failure from Run %v.", result.Failures)
		}
	}()

	fmt.Printf("END TIME %s \n", time.Now().String())
}

func monitorResults(wg *sync.WaitGroup, runResults chan apply.Results) {
	wg.Wait()
	close(runResults)
}
