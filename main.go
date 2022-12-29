package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var hconn = "tcp://0.0.0.0:5555"

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	list := &cobra.Command{
		Use:   "list",
		Short: "list docker containers",
		Long:  "list docker containers, butt longer",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Test")
			host := hconn

			dc, err := dockerClient.NewClientWithOpts(
				dockerClient.WithHost(host),
				dockerClient.WithAPIVersionNegotiation(),
			)
			log.Error().Err(err).Msg("Error while creating docker client")
			if err != nil {
				return
			}
			defer dc.Close()

			containers, err := dc.ContainerList(context.Background(), types.ContainerListOptions{})
			if err != nil {
				log.Error().Err(err).Msg("Error while getting containers")
				return
			}
			d, err := json.Marshal(containers)
			if err != nil {
				log.Error().Err(err).Msg("Error while marshalling container list")
				return
			}
			log.Info().RawJSON("containers", d).Msg("here are the containers")

		},
	}

	create := &cobra.Command{
		Use:   "create",
		Short: "create and start docker containers",
		Long:  "create and start longer docker containers",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Starting create")
			host := hconn

			dc, err := dockerClient.NewClientWithOpts(
				dockerClient.WithHost(host),
				dockerClient.WithAPIVersionNegotiation(),
			)

			if err != nil {
				log.Error().Err(err).Msg("Error when starting docker container.")
				return
			}

			defer dc.Close()

			resp, err := dc.ContainerCreate(context.Background(), &container.Config{
				Image: "alpine",
				Cmd:   []string{"echo", "hello world"},
				Tty:   false,
			}, nil, nil, nil, "")

			if err != nil {
				log.Error().Err(err).Msg("why the fuck is our logging command so long?")
			}
			if err := dc.ContainerStart(context.Background(),
				resp.ID, types.ContainerStartOptions{}); err != nil {
				log.Error().Err(err).Msg("Seriously what the fuck?")
			}
			stats, errstat := dc.ContainerWait(context.Background(), resp.ID, container.WaitConditionNotRunning)
			select {
			case err := <-errstat:
				if err != nil {
					log.Error().Err(err).Msg("Writing logs has never been more painful")
				}
			case <-stats:
			}

			out, err := dc.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				log.Error().Err(err).Msg("We need to fix this oh my god. We have to write these after almost ever statement. Jesus.")
			}
			stdcopy.StdCopy(os.Stdout, os.Stderr, out)
		},
	}
	create.Execute()
	list.Execute()
}
