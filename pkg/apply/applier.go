package apply

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

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

func (a *Applier) applyNamespaceDirs(wg *sync.WaitGroup, chunkFolder []string) {
	defer wg.Done()

	for _, folder := range chunkFolder {

		a.apply_kubernetes_files(folder)
		a.apply_terraform(folder)
	}

}

func (a *Applier) apply_kubernetes_files(folder string) {

	var outb, errb bytes.Buffer

	kubectlArgs := []string{"-n", filepath.Base(folder), "apply", "-f", folder}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Stdout = &outb
	kubectlCommand.Stderr = &errb
	kubectlCommand.Run()

	result := &Result{Response: errb.String(), Error: errb.String(), Folder: folder}

	a.RunResults <- *result
}
func (a *Applier) apply_terraform(folder string) {

	// Get the value of an Environment Variable
	bucket := os.Getenv("PIPELINE_STATE_BUCKET")
	key_prefix := os.Getenv("PIPELINE_STATE_KEY_PREFIX")
	lock_table := os.Getenv("PIPELINE_TERRAFORM_STATE_LOCK_TABLE")
	region := os.Getenv("PIPELINE_STATE_REGION")
	cluster := os.Getenv("PIPELINE_CLUSTER")

	// //Checking that an environment variable is present or not.

	key := key_prefix + cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

	// // err := os.RemoveAll(filepath.Join(folder, ".terraform"))
	// // if err != nil {
	// // 	result = Result{Error: "Cant remove .terraform folders", Response: "", Folder: folder}
	// // 	results <- result
	// // 	return
	// // }

	var outb, errb bytes.Buffer

	kubectlArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", bucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", lock_table),
		fmt.Sprintf("%s=region=%s", "-backend-config", region)}

	Command := exec.Command("terraform", kubectlArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()

	kubectlArgs = []string{"plan"}

	Command = exec.Command("terraform", kubectlArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()

	kubectlArgs = []string{"apply"}

	Command = exec.Command("terraform", kubectlArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()

	result := &Result{Response: outb.String(), Error: errb.String(), Folder: folder}

	a.RunResults <- *result
}
