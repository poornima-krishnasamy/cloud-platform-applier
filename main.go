// Package main Declaration
package main

// Importing packages
import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/apply"
	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/config"
)

// Main function
func main() {

	fmt.Printf("START TIME %s \n", time.Now().String())

	fmt.Println("Version", runtime.Version())
	fmt.Println("NumCPU", runtime.NumCPU())
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	config := config.NewEnvPipelineConfig()

	//wg := &sync.WaitGroup{}

	// clock := &sysutil.Clock{}

	// runResults := make(chan apply.Results, 5)

	apply.FullRun(config)

	// TODO Fix channel output

	// go monitorResults(wg, runResults)
	// go func() {
	// 	for result := range runResults {
	// 		log.Printf("Updating successes Run %v.", result.Successes)
	// 		log.Printf("Updating failure from Run %v.", result.Failures)
	// 	}
	// }()

	fmt.Printf("END TIME %s \n", time.Now().String())
}

func monitorResults(wg *sync.WaitGroup, runResults chan apply.Results) {
	wg.Wait()
	close(runResults)
}
