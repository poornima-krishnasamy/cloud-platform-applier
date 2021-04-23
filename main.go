// Package main Declaration
package main

// Importing packages
import (
	"fmt"
	"log"

	sysutil "github.com/cloud-platform-applier/sysutil"
)

const (
	workerCount = 1
	logLevel    = -1
)

// Main function
// env vars: REPO_PATH
func main() {

	repoPath := sysutil.GetRequiredEnvString("REPO_PATH")
	// clusterName := sysutil.GetRequiredEnvString("TF_VAR_cluster_name")
	// clusterStateBucket := sysutil.GetRequiredEnvString("TF_VAR_cluster_state_bucket")
	// clusterStateKey := sysutil.GetRequiredEnvString("TF_VAR_cluster_state_key")

	// clock := &sysutil.Clock{}

	fileSystem := &sysutil.FileSystem{}

	folderList, err := fileSystem.ListFolders(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, folder := range folderList {
		fmt.Printf("Found directory %v\n", folder)
	}

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
