package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)


func main() {
	hcfg, err := nConn(*cInfo, hateConfig{
		Protocol: "tcp://",
		Host: "0.0.0.0",
		Port: "5555",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("A fatal error has occurred.")
	}
	defer hcfg.close()
	hConn = fmt.Sprintf("%s%s%s", &hcfg.Protocol, &hcfg.Host, &hcfg.Port)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()
	list := &cobra.Command{
		Use:   "list",
		Short: "list docker containers",
		Long:  "list docker containers, butt longer",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Test")
			host := hConn 

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
