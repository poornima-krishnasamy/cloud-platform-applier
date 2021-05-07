// Package main Declaration
package main

// Importing packages
import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
)

const (
	workerCount = 1
	logLevel    = -1
)

// Main function
// env vars: REPO_PATH
func main() {

	fmt.Printf("START TIME %s \n", time.Now().String())

	fmt.Println("Version", runtime.Version())
	fmt.Println("NumCPU", runtime.NumCPU())
	fmt.Println("GOMAXPROCS", runtime.GOMAXPROCS(0))

	const nRoutines int = 3

	repoPath := sysutil.GetRequiredEnvString("REPO_PATH")
	// clusterName := sysutil.GetRequiredEnvString("TF_VAR_cluster_name")
	// clusterStateBucket := sysutil.GetRequiredEnvString("TF_VAR_cluster_state_bucket")
	// clusterStateKey := sysutil.GetRequiredEnvString("TF_VAR_cluster_state_key")

	// clock := &sysutil.Clock{}

	folders, err := sysutil.ListFolderPaths(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folders {
		fmt.Printf("Found directory %v\n", folder)
	}

	fileSystem := &sysutil.FileSystem{}

	folderChunks, err := fileSystem.ChunkFolders(folders, nRoutines)
	if err != nil {
		log.Fatal(err)
	}

	//results := make(chan Results)

	//wg := &sync.WaitGroup{}

	fmt.Println("Number of Chunks", len(folderChunks))

	// for i := 0; i < len(folderChunks); i++ {
	// 	wg.Add(1)
	// 	go applyNamespaceDirs(wg, results, folderChunks[i])

	// }

	// if err := kube.CheckVersion(); err != nil {
	// 	log.Fatal(err)
	// }

	// successes = []ApplyAttempt{}
	// failures = []ApplyAttempt{}

	// for _, path := range applyList {
	// 	log.Printf("RUN Applying file %v", path)
	// 	cmd, output, err := kube.Apply(path)
	// 	success := (err == nil)
	// 	appliedFile := ApplyAttempt{path, cmd, output, ""}
	// 	if success {
	// 		successes = append(successes, appliedFile)
	// 		log.Printf("RUN %v: %v\n%v", id, cmd, output)
	// 	} else {
	// 		appliedFile.ErrorMessage = err.Error()
	// 		failures = append(failures, appliedFile)
	// 		log.Printf("RUN %v: %v\n%v\n%v", id, cmd, output, appliedFile.ErrorMessage)
	// 	}
	// }

	// go monitorResults(wg, results)

	// for result := range results {
	// 	fmt.Printf("Folder: %v\n", result.Folder)
	// 	fmt.Printf("Response: %v\n", result.Response)
	// 	if result.Error != "" {
	// 		fmt.Printf("Error: %v", result.Error)
	// 		continue
	// 	}
	// }

	fmt.Printf("END TIME %s \n", time.Now().String())

	// fullRunQueue := make(chan bool, 1)
	// runResults := make(chan run.Result, 5)
	// errors := make(chan error)

	// runCount := make(chan int)

	// kubeClient := &kube.Client{
	// 	Server:   server,
	// 	LogLevel: logLevel,
	// }
	// kubeClient.Configure()

	// runner := &run.Runner{
	// 	clock,
	// 	fullRunQueue,
	// 	runResults,
	// 	errors,
	// }

	// // scheduler := &run.Scheduler{fullRunQueue, errors, ""}

	// // scheduler.Start()
	// go runner.StartFullLoop()

	// for err := range errors {
	// 	log.Fatal(err)
	// }

	// ctx = signals.SetupSignalHandler()
	// <-ctx.Done()
	// log.Logger("kube-applier").Info("Interrupted, shutting down...")
	// scheduler.Stop()
	// runner.Stop()
}
