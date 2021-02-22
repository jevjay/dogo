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
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var ctx context.Context
var cli *client.Client

// Container stores data about Docker container and its output
type Container struct {
	Con container.ContainerCreateCreatedBody // The container
	Out string
}

// ExecuteDocker used to start new Docker container
func ExecuteDocker(name string, image string) (*Container, error) {
	ctx = context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		return nil, err
	}
	// Pull image if such does not exist
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	io.Copy(os.Stdout, out)
	// Generate container name
	b := make([]byte, 4)
	rand.Read(b)
	// array of strings.
	str := []string{name, fmt.Sprintf("%x", b)}
	// Create a Docker container
	resp, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        image,
			Tty:          true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
		},
		nil, nil, strings.Join(str, "-"))

	if err != nil {
		return nil, err
	}
	// Run container
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	// Configure container logger options
	options := types.ContainerLogsOptions{
		ShowStderr: true,
		ShowStdout: true,
		Timestamps: false,
	}
	// Retrieve logs from container
	out, err = cli.ContainerLogs(ctx, resp.ID, options)
	if err != nil {
		return nil, err
	}
	// Read container logs
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	// Construct container struct
	con := &Container{
		Con: resp,
		Out: buf.String(),
	}
	return con, nil
}
