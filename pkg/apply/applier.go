package apply

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/config"
	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
)

// Applier makes apply calls for a list of files.
// type Applier struct {
// 	Config     EnvPipelineConfig
// 	RunResults chan<- Results
// }

//Results collects the list of ApplyAttempt
type Results struct {
	Start     string
	Finish    string
	Successes []ApplyAttempt
	Failures  []ApplyAttempt
}

// ApplyAttempt collects each of folder path, output and error message
type ApplyAttempt struct {
	FilePath     string
	Output       string
	ErrorMessage string
}

func FullRun(config *config.EnvPipelineConfig) (results []Results) {

	//wg := &sync.WaitGroup{}

	folderChunks := sysutil.PrepareFolders(config)

	//results := []Results{}
	for i := 0; i < len(folderChunks); i++ {
		// wg.Add(1)

		start := time.Now().String()
		successes, failures := applyNamespaceDirs(config, folderChunks[i])

		finish := time.Now().String()

		newRun := Results{start, finish, successes, failures}

		results = append(results, newRun)
		// results <- *newRun

	}
	return results
}

func applyNamespaceDirs(config *config.EnvPipelineConfig, chunkFolder []string) (successes []ApplyAttempt, failures []ApplyAttempt) {

	for _, folder := range chunkFolder {

		log.Printf("RUN : Applying file %v", folder)
		outputKubectl, errKubectl := applyKubectl(config, folder)
		successKubectl := (errKubectl == "")
		outputTf, errTf := planTerraform(config, folder)
		successTf := (errTf == "")

		output := outputKubectl + "\n" + outputTf
		err := errKubectl + errTf
		success := successKubectl || successTf

		appliedFile := ApplyAttempt{folder, output, ""}
		if success {
			successes = append(successes, appliedFile)
			// log.Printf("RUN : %v", output)
		} else {
			appliedFile.ErrorMessage = err
			failures = append(failures, appliedFile)
			// log.Printf("RUN : %v\n%v", output, appliedFile.ErrorMessage)
		}
	}
	return successes, failures
}

// Apply attempts to "kubectl apply" the file located at path.
// It returns the full apply command and its output.
func applyKubectl(config *config.EnvPipelineConfig, folder string) (output string, err string) {

	var outb, errb bytes.Buffer

	kubectlArgs := []string{"-n", filepath.Base(folder), "apply", "-f", folder}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Stdout = &outb
	kubectlCommand.Stderr = &errb
	kubectlCommand.Run()

	return outb.String(), errb.String()

}

func applyTerraform(config *config.EnvPipelineConfig, folder string) (output string, err string) {

	key := config.StateKeyPrefix + config.Cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

	var outb, errb bytes.Buffer

	kubectlArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", config.StateLockTable),
		fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

	Command := exec.Command("terraform", kubectlArgs...)

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
	return outb.String(), errb.String()

}

func planTerraform(config *config.EnvPipelineConfig, folder string) (output string, err string) {

	key := config.StateKeyPrefix + config.Cluster + "/" + filepath.Base(folder) + "/terraform.tfstate"

	var outb, errb bytes.Buffer

	kubectlArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", config.StateLockTable),
		fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

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
