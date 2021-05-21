package apply

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
	papply "github.com/pytimer/k8sutil/apply"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var config *rest.Config

type Applier struct {
	RunResults chan<- Result
}

type Result struct {
	Response string
	Error    string
	Folder   string
}

const (
	nRoutines = 2
)

func (a *Applier) Apply(folderChunks [][]string, kubeconfig *rest.Config) {
	config = kubeconfig

	// wg := &sync.WaitGroup{}

	for i := 0; i < len(folderChunks); i++ {
		// wg.Add(1)
		applyNamespaceDirs(folderChunks[i])

	}

}

func applyNamespaceDirs(chunkFolder []string) {

	// defer wg.Done()

	for _, folder := range chunkFolder {

		apply_kubernetes_files(folder)
		// a.apply_terraform(folder)
	}
}

func apply_kubernetes_files(folder string) {

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fileSystem := &sysutil.FileSystem{}

	files, err := fileSystem.ListFiles(folder)
	if err != nil {
		panic(err.Error())
	}
	ctx := context.Background()

	for _, file := range files {

		fmt.Println("Applying file %s", file)

		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		applyOptions := papply.NewApplyOptions(dynamicClient, discoveryClient)
		if err := applyOptions.Apply(ctx, []byte(content)); err != nil {
			log.Fatalf("apply error: %v", err)
		}

	}

}

// func (a *Applier) apply_terraform(folder string) {

// 	// Get the value of an Environment Variable
// 	bucket := os.Getenv("PIPELINE_STATE_BUCKET")
// 	key_prefix := os.Getenv("PIPELINE_STATE_KEY_PREFIX")
// 	lock_table := os.Getenv("PIPELINE_TERRAFORM_STATE_LOCK_TABLE")
// 	region := os.Getenv("PIPELINE_STATE_REGION")
// 	cluster := os.Getenv("PIPELINE_CLUSTER")

// 	// //Checking that an environment variable is present or not.

// 	key := key_prefix + cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

// 	// // err := os.RemoveAll(filepath.Join(folder, ".terraform"))
// 	// // if err != nil {
// 	// // 	result = Result{Error: "Cant remove .terraform folders", Response: "", Folder: folder}
// 	// // 	results <- result
// 	// // 	return
// 	// // }

// 	var outb, errb bytes.Buffer

// 	kubectlArgs := []string{
// 		"init",
// 		fmt.Sprintf("%s=bucket=%s", "-backend-config", bucket),
// 		fmt.Sprintf("%s=key=%s", "-backend-config", key),
// 		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", lock_table),
// 		fmt.Sprintf("%s=region=%s", "-backend-config", region)}

// 	Command := exec.Command("terraform", kubectlArgs...)

// 	Command.Dir = folder + "/resources"
// 	Command.Stdout = &outb
// 	Command.Stderr = &errb
// 	Command.Run()

// 	kubectlArgs = []string{"plan"}

// 	Command = exec.Command("terraform", kubectlArgs...)

// 	Command.Dir = folder + "/resources"
// 	Command.Stdout = &outb
// 	Command.Stderr = &errb
// 	Command.Run()

// 	kubectlArgs = []string{"apply"}

// 	Command = exec.Command("terraform", kubectlArgs...)

// 	Command.Dir = folder + "/resources"
// 	Command.Stdout = &outb
// 	Command.Stderr = &errb
// 	Command.Run()

// 	result := &Result{Response: outb.String(), Error: errb.String(), Folder: folder}

// 	a.RunResults <- *result
// }
