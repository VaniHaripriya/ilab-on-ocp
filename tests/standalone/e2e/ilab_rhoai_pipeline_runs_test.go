/*
Copyright 2025.

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

package odh

import (
	"os"
	"testing"

	TestUtil "github.com/opendatahub-io/ilab-on-ocp/tests/standalone/e2e/util"
	"github.com/stretchr/testify/require"
)

func TestPipelineRun(t *testing.T) {

	if os.Getenv("ENABLE_ILAB_PIPELINE_TEST") != "true" {
		t.Skip("Skipping iLab pipeline test. Set ENABLE_ILAB_PIPELINE_TEST=true to enable.")
	}

	pipelineServerURL := os.Getenv("PIPELINE_SERVER_URL")
	require.NotEmpty(t, pipelineServerURL, "PIPELINE_SERVER_URL environment variable must be set")

	bearerToken := os.Getenv("BEARER_TOKEN")
	require.NotEmpty(t, bearerToken, "BEARER_TOKEN environment variable must be set")

	pipelineDisplayName := os.Getenv("PIPELINE_DISPLAY_NAME")
	require.NotEmpty(t, pipelineDisplayName, "PIPELINE_DISPLAY_NAME environment variable must be set")

	// Retrieve the pipeline ID
	pipelineID, err := TestUtil.retrievePipelineId(t, pipelineServerURL, pipelineDisplayName, bearerToken)
	require.NoError(t, err, "Failed to retrieve pipeline ID")

	// Define input parameters for the pipeline
	parameters := map[string]interface{}{
		"final_eval_batch_size":                "auto",
		"final_eval_few_shots":                 5,
		"final_eval_max_workers":               "auto",
		"final_eval_merge_system_user_message": false,
		"k8s_storage_class_name":               "nfs-csi",
		"mt_bench_max_workers":                 "auto",
		"mt_bench_merge_system_user_message":   false,
		"sdg_base_model":                       "s3://rhods-dsp-dev/granite-7b-starter",
		"sdg_max_batch_len":                    5000,
		"sdg_pipeline":                         "simple",
		"sdg_repo_branch":                      "",
		"sdg_repo_pr":                          0,
		"sdg_repo_url":                         "https://github.com/instructlab/taxonomy.git",
		"sdg_sample_size":                      0.00002,
		"sdg_scale_factor":                     30,
		"train_effective_batch_size_phase_1":   3840,
		"train_effective_batch_size_phase_2":   3840,
		"train_learning_rate_phase_1":          0.1,
		"train_learning_rate_phase_2":          0.1,
		"train_max_batch_len":                  20000,
		"train_nnodes":                         1,
		"train_nproc_per_node":                 1,
		"train_num_epochs_phase_1":             1,
		"train_num_epochs_phase_2":             1,
		"train_num_warmup_steps_phase_1":       800,
		"train_num_warmup_steps_phase_2":       800,
		"train_save_samples":                   0,
		"train_seed":                           42,
	}

	// Trigger the pipeline run
	runID, err := TestUtil.triggerPipeline(t, pipelineServerURL, pipelineID, pipelineDisplayName, parameters, bearerToken)
	require.NoError(t, err, "Failed to trigger pipeline")

	// Verify the pipeline's successful completion
	err = TestUtil.verifyPipelineSuccess(t, pipelineServerURL, runID, bearerToken)
	require.NoError(t, err, "Pipeline did not complete successfully")
}
