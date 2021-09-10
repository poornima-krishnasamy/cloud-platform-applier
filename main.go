// Package main Declaration
package main

// Importing packages
import (
	"fmt"
	"runtime"
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

	results := apply.FullRun(config)

	fmt.Printf("Printing Failures \n")
	for _, result := range results {
		if len(result.Failures) > 0 {
			fmt.Println(result.Failures)
		}
	}

	fmt.Printf("END TIME %s \n", time.Now().String())
}
