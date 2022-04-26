# Cloud Platform Environments Applier 


Set these env variables 

export AWS_REGION=eu-west-2
export AWS_PROFILE=moj-cp

export TF_VAR_cluster_name="cp-2004-1705"
export TF_VAR_cluster_state_bucket="cloud-platform-terraform-tfstate"
export TF_VAR_cluster_state_key="cloud-platform/cp-2004-1705/terraform.tfstate"
export TF_VAR_kubernetes_cluster=" https://cluster-endpoint"

The rest of the env variables are set as flag defaults.
You can overwrite by passing the flags or setting an env variable "PIPELINE_FLAGNAME"

