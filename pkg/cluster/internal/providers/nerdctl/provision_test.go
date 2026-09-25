/*
Copyright 2026 The Kubernetes Authors.

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

package nerdctl

import (
	"reflect"
	"testing"

	"sigs.k8s.io/kind/pkg/internal/apis/config"
)

func TestRunArgsForNodeGPUs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		node         *config.Node
		expectedGPUs []string
	}{
		{
			name: "no gpus",
			node: &config.Node{
				Role:  config.WorkerRole,
				Image: "kindest/node:latest",
			},
			expectedGPUs: nil,
		},
		{
			name: "gpus all",
			node: &config.Node{
				Role:  config.WorkerRole,
				Image: "kindest/node:latest",
				GPUs:  "all",
			},
			expectedGPUs: []string{"--gpus", "all"},
		},
		{
			name: "gpus specific id",
			node: &config.Node{
				Role:  config.WorkerRole,
				Image: "kindest/node:latest",
				GPUs:  "0,1",
			},
			expectedGPUs: []string{"--gpus", "0,1"},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			args, err := runArgsForNode(tc.node, config.IPv4Family, "node-1", nil)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			var gpusArgs []string
			for i := 0; i < len(args); i++ {
				if args[i] == "--gpus" && i+1 < len(args) {
					gpusArgs = []string{args[i], args[i+1]}
					break
				}
			}

			if !reflect.DeepEqual(gpusArgs, tc.expectedGPUs) {
				t.Errorf("expected gpu args %v, got %v", tc.expectedGPUs, gpusArgs)
			}
		})
	}
}
