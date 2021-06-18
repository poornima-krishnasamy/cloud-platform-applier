package apply

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
)

// Applier makes apply calls for a list of files.
type Applier struct {
	RepoPath   string
	RunResults chan<- Results
}

type Results struct {
	Start     string
	Finish    string
	Successes []ApplyAttempt
	Failures  []ApplyAttempt
}

type ApplyAttempt struct {
	FilePath     string
	Output       string
	ErrorMessage string
}

const (
	nRoutines = 2
)

func (a *Applier) FullRun() {

	// wg := &sync.WaitGroup{}

	folders, err := sysutil.ListFolderPaths(a.RepoPath)
	if err != nil {
		log.Fatal(err)
	}

	folderChunks, err := sysutil.ChunkFolders(folders, nRoutines)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(folderChunks); i++ {
		// wg.Add(1)

		start := time.Now().String()
		successes, failures := a.applyNamespaceDirs(folderChunks[i])

		finish := time.Now().String()

		newRun := &Results{start, finish, successes, failures}

		a.RunResults <- *newRun

	}

}

func (a *Applier) applyNamespaceDirs(chunkFolder []string) (successes []ApplyAttempt, failures []ApplyAttempt) {

	for _, folder := range chunkFolder {

		log.Printf("RUN : Applying file %v", folder)
		output, err := a.ApplyKubectl(folder)
		success := (err == "")
		appliedFile := ApplyAttempt{folder, output, ""}
		if success {
			successes = append(successes, appliedFile)
			log.Printf("RUN : %v", output)
		} else {
			appliedFile.ErrorMessage = err
			failures = append(failures, appliedFile)
			log.Printf("RUN : %v\n%v", output, appliedFile.ErrorMessage)
		}
	}
	return successes, failures
}

// Apply attempts to "kubectl apply" the file located at path.
// It returns the full apply command and its output.
func (a *Applier) ApplyKubectl(folder string) (output string, err string) {

	var outb, errb bytes.Buffer

	kubectlArgs := []string{"-n", filepath.Base(folder), "apply", "-f", folder}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Stdout = &outb
	kubectlCommand.Stderr = &errb
	kubectlCommand.Run()

	return outb.String(), errb.String()

}

func (a *Applier) ApplyTerraform(folder string) (output string, err string) {

	// Get the value of an Environment Variable
	bucket := os.Getenv("PIPELINE_STATE_BUCKET")
	key_prefix := os.Getenv("PIPELINE_STATE_KEY_PREFIX")
	lock_table := os.Getenv("PIPELINE_TERRAFORM_STATE_LOCK_TABLE")
	region := os.Getenv("PIPELINE_STATE_REGION")
	cluster := os.Getenv("PIPELINE_CLUSTER")

	key := key_prefix + cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

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
	return outb.String(), errb.String()

}

func (a *Applier) PlanTerraform(folder string) (output string, err string) {

	// Get the value of an Environment Variable
	bucket := os.Getenv("PIPELINE_STATE_BUCKET")
	key_prefix := os.Getenv("PIPELINE_STATE_KEY_PREFIX")
	lock_table := os.Getenv("PIPELINE_TERRAFORM_STATE_LOCK_TABLE")
	region := os.Getenv("PIPELINE_STATE_REGION")
	cluster := os.Getenv("PIPELINE_CLUSTER")

	key := key_prefix + cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

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

	kubectlArgs = []string{"apply"}

	Command = exec.Command("terraform", kubectlArgs...)
	return outb.String(), errb.String()
}
