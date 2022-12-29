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
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Conn struct {
		Protocol string `yaml:"protocol"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	}
}

func HConfig(params ...string) Config {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	file, err := os.Open(`hate.yml`)

	//TODO: Check if the file exists, if it doesn't make it.
	var hcfg Config

	if err != nil {
		log.Error().Err(err).Msg("Either the hate config is missing or corrupt. Please check your path.")
	}

	defer file.Close()

	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&hcfg); err != nil {
			log.Error().Err(err).Msg("An error occurred while decoding configuration file.")
		}
	}

	return hcfg

}

func main() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	log := zerolog.New(output).With().Timestamp().Logger()

	hconn := fmt.Sprintf("%s%s:%s", HConfig().Conn.Protocol, HConfig().Conn.Host, HConfig().Conn.Port)
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
	list.Execute()
}
