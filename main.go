// Package main Declaration
package main

// Importing packages
import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/apply"
	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
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

	results := make(chan apply.Result, 5)

	applier := &apply.Applier{}

	// Location of kubeconfig file
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// TODO: create a pool of threads and spread the folders to the given threads. This is because
	// The number of max threads which can call the AWS api should be limited to avoid exceeding the rate limits

	fmt.Println("Number of Chunks", len(folderChunks))

	applier.Apply(folderChunks, config)

	// go monitorResults(wg, results)

	for result := range results {
		fmt.Printf("Folder: %v\n", result.Folder)
		fmt.Printf("Response: %v\n", result.Response)
		if result.Error != "" {
			fmt.Printf("Error: %v", result.Error)
			continue
		}
	}

	fmt.Printf("END TIME %s \n", time.Now().String())
}
