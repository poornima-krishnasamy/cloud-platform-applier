package apply

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	terraform "github.com/ministryofjustice/cloud-platform-cli/pkg/terraform"
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
		err := fmt.Errorf("error running kubectl dry-run on namespace %s: %v", config.Folder, err)
		return "", err

	}

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

	outputInitTf, err := runTerraform(config, tfArgs)
	if err != nil {
		err := fmt.Errorf("error running terraform init on namespace %s: %v", config.Folder, err)
		return "", err

	}

	tfArgs = []string{"apply"}
	outputApplyTf, err := runTerraform(config, tfArgs)
	if err != nil {
		err := fmt.Errorf("error running terraform plan  on namespace %s: %v", config.Folder, err)
		return "", err

	}
	output := outputKubectl + "\n" + outputInitTf.Stdout + "\n" + outputApplyTf.Stdout

	return output, err
}

func PlanNamespace(config *ApplierConfig) error {
	log.Printf("RUN :  file %v", config.Folder)
	config.Dryrun = true
	outputKubectl, err := applyKubectl(config)
	if err != nil {
		err := fmt.Errorf("error running kubectl dry-run on namespace %s: %v", config.Folder, err)
		return err
	}

	key := config.StateKeyPrefix + config.Cluster + "/" + filepath.Base(config.Folder) + "/terraform.tfstate"

	tfArgs := []string{
		"init",
		fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
		fmt.Sprintf("%s=key=%s", "-backend-config", key),
		fmt.Sprintf("%s=dynamodb_table=%s", "-backend-config", config.StateLockTable),
		fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

	// tfArgs := []string{
	// 	"init",
	// 	fmt.Sprintf("%s=bucket=%s", "-backend-config", config.StateBucket),
	// 	fmt.Sprintf("%s=key=%s", "-backend-config", key),
	// 	fmt.Sprintf("%s=region=%s", "-backend-config", config.StateRegion)}

	outputInitTf, err := runTerraform(config, tfArgs)
	if err != nil {
		err := fmt.Errorf("error running terraform init on namespace %s: %v: %v", config.Folder, err.Error(), outputInitTf.Stderr)
		return err

	}

	tfArgs = []string{"plan"}
	outputPlanTf, err := runTerraform(config, tfArgs)
	if err != nil {
		err := fmt.Errorf("error running terraform plan  on namespace %s: %v: %v", config.Folder, err.Error(), outputPlanTf.Stderr)
		return err

	}
	output := outputKubectl + "\n" + outputInitTf.Stdout + "\n" + outputPlanTf.Stdout

	fmt.Printf("Output of Namespace changes %s", output)
	return nil
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

func runTerraform(config *ApplierConfig, tfArgs []string) (output *terraform.CmdOutput, err error) {

	Command := exec.Command("terraform", tfArgs...)

	Command.Dir = config.RepoPath + "/" + config.Folder + "/resources"

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	var exitCode int

	Command.Stdout = &stdoutBuf
	Command.Stderr = &stderrBuf

	err = Command.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		}
		cmdOutput := terraform.CmdOutput{
			Stdout:   stdoutBuf.String(),
			Stderr:   stderrBuf.String(),
			ExitCode: exitCode,
		}
		return &cmdOutput, err
	} else {
		ws := Command.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	cmdOutput := terraform.CmdOutput{
		Stdout:   stdoutBuf.String(),
		Stderr:   stderrBuf.String(),
		ExitCode: exitCode,
	}

	if cmdOutput.ExitCode != 0 {
		return &cmdOutput, err
	} else {
		return &cmdOutput, nil
	}
}
