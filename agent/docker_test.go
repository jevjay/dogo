package agent

import (
	"strings"
	"testing"
)

func TestExecute(t *testing.T) {
	ctx, err := Execute("hello", "hello-world")

	if err != nil {
		t.Errorf("Failed to start local docker container. Error: %s", err.Error())
	}

	if !strings.Contains(ctx.Out, "Hello from Docker!") {
		t.Errorf("Failed to start retrieve logs from started container.")
	}
}
