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
	"fmt"

	"github.com/poornima-krishnasamy/cloud-platform-applier/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// planCmd represents the plan command
var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("plan called")
	},
}

func init() {

	rootCmd.AddCommand(planCmd)
	var config config.EnvPipelineConfig

	planCmd.PersistentFlags().StringVarP(&config.StateBucket, "pipeline-state-bucket", "", "", "State bucket where terraform state file is stored")
	planCmd.PersistentFlags().StringVarP(&config.StateKeyPrefix, "pipeline-state-key-prefix", "", "", "State buucket key prefix location")
	planCmd.PersistentFlags().StringVarP(&config.StateLockTable, "pipeline-state-locktable", "", "", "DynamoDB table to store the state md5")

	planCmd.PersistentFlags().StringVarP(&config.StateRegion, "pipeline-state-region", "", "", "AWS Region")

	planCmd.PersistentFlags().StringVarP(&config.Cluster, "pipeline-cluster", "", "", "Cluster to which the manifest will be applied")

	planCmd.PersistentFlags().StringVarP(&config.RepoPath, "pipeline-repo-path", "", "", "Repository folder path where the namespace manifest are")
	planCmd.PersistentFlags().IntVarP(&config.NumRoutines, "pipeline-routines", "", 2, "Num of go routines to split the folder into")

	planCmd.MarkPersistentFlagRequired("pipeline-state-bucket")
	planCmd.MarkPersistentFlagRequired("pipeline-state-key-prefix")
	planCmd.MarkPersistentFlagRequired("pipeline-state-locktable")
	planCmd.MarkPersistentFlagRequired("pipeline-state-region")
	planCmd.MarkPersistentFlagRequired("pipeline-cluster")
	planCmd.MarkPersistentFlagRequired("pipeline-repo-path")

	addCommonFlags(plan, &options)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// planCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// planCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addCommonFlags(cmd *cobra.Command, o *config.EnvPipelineConfig) {
	viper.AutomaticEnv()

}
