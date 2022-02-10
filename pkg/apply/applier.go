package apply

import (
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
	Dryrun                                                                      bool
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
		output, err := applyNamespace(config)

		appliedStatus := ApplyAttempt{folder, output, ""}
		if err != nil {
			successes = append(successes, appliedStatus)
			// log.Printf("RUN : %v", output)
		} else {
			appliedStatus.ErrorMessage = err.Error()
			failures = append(failures, appliedStatus)
			//log.Printf("RUN : %v\n%v", output, appliedFile.ErrorMessage)
		}

	}
	return successes, failures
}

func applyNamespace(config *ApplierConfig) (string, error) {
	log.Printf("RUN :  file %v", config.Folder)
	outputKubectl, err := applyKubectl(config)
	if err != nil {
		err := fmt.Errorf("Error running kubectl dry-run on namespace %s: %v", config.Folder, err)
		return "", err

	}

	outputInitTf, err := initTerraform(config)
	if err != nil {
		err := fmt.Errorf("Error running terraform init on namespace %s: %v", config.Folder, err)
		return "", err

	}

	outputApplyTf, err := applyTerraform(config)
	if err != nil {
		err := fmt.Errorf("Error running terraform plan  on namespace %s: %v", config.Folder, err)
		return "", err

	}
	output := outputKubectl + "\n" + outputInitTf + "\n" + outputApplyTf

	return output, err
}

func planNamespace(config *ApplierConfig) (string, error) {
	log.Printf("RUN :  file %v", config.Folder)
	config.Dryrun = true
	outputKubectl, err := applyKubectl(config)
	if err != nil {
		err := fmt.Errorf("error running kubectl dry-run on namespace %s: %v", config.Folder, err)
		return "", err

	}

	outputInitTf, err := initTerraform(config)
	if err != nil {
		err := fmt.Errorf("error running terraform init on namespace %s: %v", config.Folder, err)
		return "", err

	}

	outputPlanTf, err := planTerraform(config)
	if err != nil {
		err := fmt.Errorf("error running terraform plan  on namespace %s: %v", config.Folder, err)
		return "", err

	}
	output := outputKubectl + "\n" + outputInitTf + "\n" + outputPlanTf

	return output, err
}

func ExecutePlanNamespace(config *ApplierConfig) (output string, err error) {
	output, err = planNamespace(config)
	if err != nil {
		return "", err
	}
	return output, nil
}

// applyKubectl attempts to dryn-run of "kubectl apply" to the files in the given folder.
// It returns the apply command output and err.
func applyKubectl(config *ApplierConfig) (output string, err error) {

	kubectlArgs := []string{"-n", filepath.Base(config.Folder), "apply", "-f", "."}

	if config.Dryrun {
		kubectlArgs = append(kubectlArgs, "--dry-run")
	}

	kubectlCommand := exec.Command("kubectl", kubectlArgs...)

	kubectlCommand.Dir = config.RepoPath + "/" + config.Folder
	log.Printf("RUN :  command %v on folder %v", kubectlCommand, config.Folder)
	outb, err := kubectlCommand.Output()
	if err != nil {
		return "", err
	}

	return string(outb), nil

}

func applyTerraform(config *ApplierConfig) (output string, err error) {

	tfArgs := []string{"apply"}
	return runTerraform(config, tfArgs)

}

func planTerraform(config *ApplierConfig) (output string, err error) {

	tfArgs := []string{"plan"}
	return runTerraform(config, tfArgs)

}

func initTerraform(config *ApplierConfig) (output string, err error) {

	key := config.StateKeyPrefix + config.Cluster + "/" + filepath.Base(config.Folder) + "/terraform.tfstate"

	// tfArgs := []string{
	// 	"init",
	// 	fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
	// 	fmt.Sprintf("%s=key=%s", "-backend-config", key),
	// 	fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", config.StateLockTable),
	// 	fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

	tfArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

	fmt.Println(tfArgs)
	return runTerraform(config, tfArgs)
}

func runTerraform(config *ApplierConfig, tfArgs []string) (output string, err error) {

	Command := exec.Command("terraform", tfArgs...)

	Command.Dir = config.RepoPath + "/" + config.Folder + "/resources"
	outb, err := Command.Output()
	if err != nil {
		return "", err
	}

	return string(outb), nil
}
