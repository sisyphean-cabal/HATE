package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerClient "github.com/docker/docker/client"
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
		Short: "create docker containers",
		Long:  "create longer docker containers",
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
		},
	}

	list.Execute()
}
