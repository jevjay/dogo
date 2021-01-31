package agents

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// CreateNewContainer used to start new Docker container
func CreateNewContainer(name string, image string) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	r, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, r)
	i := rand.Int()
	rand := strconv.Itoa(i)
	// array of strings.
	str := []string{name, rand}

	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:        image,
			Tty:          true,
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			OpenStdin:    true,
		},
		nil, nil, nil, strings.Join(str, "-"))

	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, cont.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("error when start container", err)
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, cont.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Container %s is started", cont.ID)
	buf := new(bytes.Buffer)
	buf.ReadFrom(out)
	if err != nil {
		panic(err)
	}
	return buf.String(), nil
}
