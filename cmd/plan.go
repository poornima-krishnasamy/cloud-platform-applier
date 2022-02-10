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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// applierPlanCmd represents the plan command
func applierPlanCmd() *cobra.Command {
	var config apply.ApplierConfig
	planCmd := &cobra.Command{
		Use:   "plan",
		Short: "A brief description of your command",
		Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Printf("Starting Plan for namespace %v", config.Folder)

			if output, err := apply.ExecutePlanNamespace(&config); err != nil {
				log.Printf("Error executing Plan for namespace %v: %v", config.Folder, err)
				os.Exit(1)
			} else {
				log.Printf("Executing plan successful with output %v:", output)
			}
			return nil
		},
	}
	addCommonFlags(planCmd, &config)
	return planCmd
}

func addCommonFlags(planCmd *cobra.Command, config *apply.ApplierConfig) {

	viper.AutomaticEnv()

	//viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	var TF_VAR_cluster_name, TF_VAR_cluster_state_bucket, TF_VAR_cluster_state_key, TF_VAR_kubernetes_cluster string

	planCmd.PersistentFlags().StringVarP(&config.StateBucket, "PIPELINE_STATE_BUCKET", "", os.Getenv("PIPELINE_STATE_BUCKET"), "State bucket where terraform state file is stored")
	planCmd.PersistentFlags().StringVarP(&config.StateKeyPrefix, "PIPELINE_STATE_KEY_PREFIX", "", os.Getenv("PIPELINE_STATE_KEY_PREFIX"), "State buucket key prefix location")
	planCmd.PersistentFlags().StringVarP(&config.StateLockTable, "PIPELINE_TERRAFORM_STATE_LOCK_TABLE", "", os.Getenv("PIPELINE_TERRAFORM_STATE_LOCK_TABLE"), "DynamoDB table to store the state md5")

	planCmd.PersistentFlags().StringVarP(&config.StateRegion, "PIPELINE_STATE_REGION", "", os.Getenv("PIPELINE_STATE_REGION"), "AWS Region")

	planCmd.PersistentFlags().StringVarP(&config.Cluster, "PIPELINE_CLUSTER", "", os.Getenv("PIPELINE_CLUSTER"), "Cluster to which the manifest will be applied")

	planCmd.PersistentFlags().StringVarP(&config.RepoPath, "PIPELINE_REPOPATH", "", os.Getenv("PIPELINE_REPOPATH"), "Repository folder path where the namespace manifest are")
	planCmd.PersistentFlags().IntVarP(&config.NumRoutines, "PIPELINE_ROUTINES", "", 2, "Num of go routines to split the folder into")

	planCmd.PersistentFlags().StringVarP(&config.Folder, "PIPELINE_FOLDER", "", os.Getenv("NAMESPACE"), "Name of the folder to do the plan")
	planCmd.PersistentFlags().BoolVarP(&config.Dryrun, "dry-run", "", true, "dryrun option for kubectl")

	planCmd.PersistentFlags().StringVarP(&TF_VAR_cluster_name, "TF_VAR_cluster_name", "", os.Getenv("TF_VAR_cluster_name"), "State bucket where terraform state file is stored")
	planCmd.PersistentFlags().StringVarP(&TF_VAR_cluster_state_bucket, "TF_VAR_cluster_state_bucket", "", os.Getenv("TF_VAR_cluster_state_bucket"), "State bucket name")
	planCmd.PersistentFlags().StringVarP(&TF_VAR_cluster_state_key, "TF_VAR_cluster_state_key", "", os.Getenv("TF_VAR_cluster_state_key"), "State bucket key prefix location")
	planCmd.PersistentFlags().StringVarP(&TF_VAR_kubernetes_cluster, "TF_VAR_kubernetes_cluster", "", os.Getenv("TF_VAR_kubernetes_cluster"), "kubernetes cluster name")

	planCmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			planCmd.PersistentFlags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

// func init() {
// 	rootCmd.AddCommand(applierPlanCmd())
// }
