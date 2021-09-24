package apply

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/sysutil"
)

// Applier makes apply calls for a list of files.
// type Applier struct {
// 	Config     EnvPipelineConfig
// 	RunResults chan<- Results
// }

//Results collects the list of ApplyAttempt

// ApplyAttempt collects each of folder path, output and error message
type ApplyAttempt struct {
	FilePath, Output, ErrorMessage string
}

type Results struct {
	Start, Finish string
	Successes     []ApplyAttempt
	Failures      []ApplyAttempt
}

type ApplierConfig struct {
	StateBucket, StateKeyPrefix, StateLockTable, StateRegion, Cluster, RepoPath string
	NumRoutines                                                                 int
	Folder                                                                      string
}

func FullRun(config *ApplierConfig) (results []Results) {

	//wg := &sync.WaitGroup{}

	folderChunks := sysutil.GetFolderChunks(config.RepoPath, config.NumRoutines)

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

func applyNamespaceDirs(config *ApplierConfig, chunkFolder []string) (successes []ApplyAttempt, failures []ApplyAttempt) {

	for _, folder := range chunkFolder {
		output, err, isSuccess := applyNamespace(config)

		appliedStatus := ApplyAttempt{folder, output, ""}
		if isSuccess {
			successes = append(successes, appliedStatus)
			// log.Printf("RUN : %v", output)
		} else {
			appliedStatus.ErrorMessage = err
			failures = append(failures, appliedStatus)
			//log.Printf("RUN : %v\n%v", output, appliedFile.ErrorMessage)
		}

	}
	return successes, failures
}

func applyNamespace(config *ApplierConfig) (string, string, bool) {
	log.Printf("RUN : Applying file %v", config.Folder)
	outputKubectl, errKubectl := applyKubectl(config)
	successKubectl := (errKubectl == "")

	outputInitTf, errInitTf := initTerraform(config)
	successInitTf := (errInitTf == "")

	outputPlanTf, errPlanTf := planTerraform(config)
	successPlanTf := (errPlanTf == "")

	output := outputKubectl + "\n" + outputInitTf + "\n" + outputPlanTf
	err := errKubectl + "\n" + errInitTf + "\n" + errPlanTf
	isSuccess := successKubectl && successInitTf && successPlanTf

	return output, err, isSuccess
}

func planNamespace(config *ApplierConfig) (string, string, bool) {
	log.Printf("RUN :  file %v", config.Folder)
	outputKubectl, errKubectl := planKubectl(config)
	successKubectl := (errKubectl == "")

	outputInitTf, errInitTf := initTerraform(config)
	successInitTf := (errInitTf == "")

	outputPlanTf, errPlanTf := planTerraform(config)
	successPlanTf := (errPlanTf == "")

	output := outputKubectl + "\n" + outputInitTf + "\n" + outputPlanTf
	err := errKubectl + "\n" + errInitTf + "\n" + errPlanTf
	isSuccess := successKubectl && successInitTf && successPlanTf

	return output, err, isSuccess
}

func ExecutePlanNamespace(config *ApplierConfig) string {
	_, err, _ := planNamespace(config)
	if err != "" {
		return err
	}
	return ""
}

// planKubectl attempts to dryn-run of "kubectl apply" to the files in the given folder.
// It returns the apply command output and err.
func planKubectl(config *ApplierConfig) (output string, err string) {

	var outb, errb bytes.Buffer

	kubectlArgs := []string{"-n", filepath.Base(config.Folder), "apply", "--dry-run", "-f", config.Folder}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Stdout = &outb
	kubectlCommand.Stderr = &errb
	kubectlCommand.Run()

	return outb.String(), errb.String()

}

// applyKubectl attempts to dryn-run of "kubectl apply" to the files in the given folder.
// It returns the apply command output and err.
func applyKubectl(config *ApplierConfig) (output string, err string) {

	var outb, errb bytes.Buffer

	kubectlArgs := []string{"-n", filepath.Base(config.Folder), "apply", "-f", config.Folder}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Stdout = &outb
	kubectlCommand.Stderr = &errb
	kubectlCommand.Run()

	return outb.String(), errb.String()

}

func applyTerraform(config *ApplierConfig, folder string) (output string, err string) {

	tfArgs := []string{"apply"}
	return runTerraform(folder, tfArgs)

}

func planTerraform(config *ApplierConfig) (output string, err string) {

	tfArgs := []string{"plan"}
	return runTerraform(config.Folder, tfArgs)

}

func initTerraform(config *ApplierConfig) (output string, err string) {

	key := config.StateKeyPrefix + config.Cluster + "/" + filepath.Base(config.Folder) + "/terraform.tfstate"

	tfArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", config.StateLockTable),
		fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

	return runTerraform(config.Folder, tfArgs)
}

func runTerraform(folder string, tfArgs []string) (output string, err string) {

	var outb, errb bytes.Buffer

	Command := exec.Command("terraform", tfArgs...)

	Command.Dir = folder + "/resources"
	Command.Stdout = &outb
	Command.Stderr = &errb
	Command.Run()
	return outb.String(), errb.String()
}
