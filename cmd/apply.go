/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"os"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/apply"
	"github.com/spf13/cobra"
)

// applierapplyCmd represents the plan command
func applierApplyCmd() *cobra.Command {
	var config apply.ApplierConfig

	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("Starting Plan for namespace %v", config.Folder)

			if err := apply.ApplyNamespace(&config); err != nil {
				log.Printf("Error executing Plan for namespace %v: %v", config.Folder, err)
				os.Exit(1)
			}
			return nil
		},
	}

	var TF_VAR_cluster_name, TF_VAR_cluster_state_bucket, TF_VAR_cluster_state_key, TF_VAR_kubernetes_cluster string

	applyCmd.PersistentFlags().StringVarP(&config.StateBucket, "state-bucket", "", "cloud-platform-terraform-state", "State bucket where terraform state file is stored")
	applyCmd.PersistentFlags().StringVarP(&config.StateKeyPrefix, "state-key-prefix", "", "cloud-platform-environments/", "State buucket key prefix location")

	applyCmd.PersistentFlags().StringVarP(&config.StateLockTable, "terraform-state-lock-table", "", "cloud-platform-environments-terraform-lock", "DynamoDB table to store the state md5")

	applyCmd.PersistentFlags().StringVarP(&config.StateRegion, "state-region", "", "eu-west-1", "AWS Region")

	applyCmd.PersistentFlags().StringVarP(&config.Cluster, "cluster", "", "cp-2004-1705.cloud-platform.service.justice.gov.uk", "Cluster to which the manifest will be applied")

	applyCmd.PersistentFlags().StringVarP(&config.RepoPath, "repo-path", "", "namespaces/cp-2004-1705.cloud-platform.service.justice.gov.uk", "Repository folder path where the namespace manifest are")
	applyCmd.PersistentFlags().IntVarP(&config.NumRoutines, "routines", "", 2, "Num of go routines to split the folder into")

	applyCmd.PersistentFlags().StringVarP(&config.Folder, "namespace", "", "ns-test-dev", "Name of the folder to do the plan")
	applyCmd.PersistentFlags().BoolVarP(&config.Dryrun, "dry-run", "", false, "dryrun option for kubectl")

	applyCmd.PersistentFlags().StringVarP(&TF_VAR_cluster_name, "TF_VAR_cluster_name", "", os.Getenv("TF_VAR_cluster_name"), "State bucket where terraform state file is stored")
	applyCmd.PersistentFlags().StringVarP(&TF_VAR_cluster_state_bucket, "TF_VAR_cluster_state_bucket", "", os.Getenv("TF_VAR_cluster_state_bucket"), "State bucket name")
	applyCmd.PersistentFlags().StringVarP(&TF_VAR_cluster_state_key, "TF_VAR_cluster_state_key", "", os.Getenv("TF_VAR_cluster_state_key"), "State bucket key prefix location")
	applyCmd.PersistentFlags().StringVarP(&TF_VAR_kubernetes_cluster, "TF_VAR_kubernetes_cluster", "", os.Getenv("TF_VAR_kubernetes_cluster"), "kubernetes cluster name")

	return applyCmd
}
