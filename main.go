package main

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func HConfig() {
	type Config struct {
		Conn struct {
			Protocol string `yaml: "protocol"`
			Host string `yaml: "host"`
			Port string `yaml: "port"`
		}
	}

	f, err := os.Open("hate.yml")
	if err != nil {
		processError(err)
	}

	defer f.Close()

	var conf Config
	decoder := yaml.NewDecoder((f))
	err = decoder.Decode(&conf)
	if err != nil {
		log.Error().Err(err).Msg("Error when opening configuration file.")
		processError(err)
	}

}



func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	list := &cobra.Command{
		Use:   "list",
		Short: "list docker containers",
		Long:  "list docker containers, butt longer",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Test")
			host := "tcp://0.0.0.0:5555"

			dc, err := dockerClient.NewClientWithOpts(
				dockerClient.WithHost(host),
				dockerClient.WithAPIVersionNegotiation(),
			)
			if err != nil {
				log.Error().Err(err).Msg("Error while creating docker client")
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

		start := &cobra.Command{
			Use: "start",
			Short: "start a docker container",
			Long: "start docker containers but longer.",
			Run: func(cmd, *cobra.Command, args []string) {
				log.Info().Msg("Test")
				host := "tcp://0.0.0.0:5555"
			}
		}
	}

	list.Execute()

}
