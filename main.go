package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type errorMessage struct {
	err    error
	reason string
}

func dcList(hconn string, ctx context.Context, outChan chan string, errChan chan errorMessage) {
	outChan <- "Starting list containers command"
	host := hconn
	fmt.Println("1")
	dc, err := dockerClient.NewClientWithOpts(
		dockerClient.WithHost(host),
		dockerClient.WithAPIVersionNegotiation(),
	)

	if err != nil {
		errChan <- errorMessage{
			err,
			"Error while creating docker client",
		}
		return
	}
	fmt.Println("2")

	defer dc.Close()
	containers, err := dc.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		fmt.Println("nope")
		errChan <- errorMessage{
			err,
			"Error while getting containers",
		}
		return
	}
	fmt.Println("3")

	d, err := json.Marshal(containers)
	if err != nil {
		errChan <- errorMessage{
			err,
			"failed marshalling",
		}
		return
	}
	fmt.Println("4")

	outChan <- string(d)
}

func dcCreate(hconn string, ctx context.Context, outChan chan string, errChan chan errorMessage) {
	outChan <- "Starting create"
	host := hconn

	dc, err := dockerClient.NewClientWithOpts(
		dockerClient.WithHost(host),
		dockerClient.WithAPIVersionNegotiation(),
	)

	if err != nil {
		errChan <- errorMessage{
			err,
			"error when starting docker container",
		}
		return
	}

	defer dc.Close()

	resp, err := dc.ContainerCreate(ctx, &container.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "hello world"},
		Tty:   false,
	}, nil, nil, nil, "")

	if err != nil {
		errChan <- errorMessage{
			err,
			"failed to get start docker container",
		}
		return
	}
	if err := dc.ContainerStart(ctx,
		resp.ID, types.ContainerStartOptions{}); err != nil {
		log.Error().Err(err).Msg("Seriously what the fuck?")
	}
	stats, errstat := dc.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errstat:
		if err != nil {
			errChan <- errorMessage{
				err,
				"failed to wait for container.",
			}
			return
		}
	case <-stats:
	}

	out, err := dc.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		errChan <- errorMessage{
			err,
			"failed to get container logs",
		}
		return
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)
}

func main() {

	// output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	zerolog.TimeFieldFormat = time.RFC3339
	outChan := make(chan string)
	errChan := make(chan errorMessage)

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		for line := range outChan {
			log.Log().Msg(line)
		}
		wg.Done()
	}()

	go func() {
		for errorMsg := range errChan {
			fmt.Println("got one")
			log.Log().Err(errorMsg.err).Msg(errorMsg.reason)
		}
		wg.Done()
	}()

	hconn := "tcp://0.0.0.0:2375"
	ctx := context.Background()

	list := &cobra.Command{
		Use:   "list",
		Short: "list docker containers",
		Long:  "list docker containers, butt longer",
		Run: func(cmd *cobra.Command, args []string) {
			dcList(hconn, ctx, outChan, errChan)
		},
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create and start docker containers",
		Long:  "create and start longer docker containers",
		Run: func(cmd *cobra.Command, args []string) {
			dcCreate(hconn, ctx, outChan, errChan)
		},
	}

	rootCmd := &cobra.Command{
		Use: "h8",
	}

	rootCmd.AddCommand(list, create)
	rootCmd.Execute()
	close(outChan)
	close(errChan)
	wg.Wait()
	fmt.Println("done")
}
