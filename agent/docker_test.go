// Copyright 2021 tappythumbz development
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agent

import (
	"strings"
	"testing"
)

func TestExecuteDocker(t *testing.T) {
	ctx, err := ExecuteDocker("hello", "hello-world")

	if err != nil {
		t.Errorf("Failed to start local docker container. Error: %s", err.Error())
	}

	if !strings.Contains(ctx.Out, "Hello from Docker!") {
		t.Errorf("Failed to start retrieve logs from started container.")
	}
}
