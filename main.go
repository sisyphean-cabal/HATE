package main

import (
	"context"
	"net/http"

	"github.com/docker/docker/api/types"
	dockerClient "github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	var run = &cobra.Command{
		Use:   "run",
		Short: "run docker container",
		Long:  "run docker container but longer phrasing",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Test")
			httpHeaders := map[string]string{}
			host := "localhost"
			client := &http.Client{}
			version := ""

			dc, err := dockerClient.NewClientWithOpts(
				dockerClient.WithHTTPClient(client),
				dockerClient.WithHTTPHeaders(httpHeaders),
				dockerClient.WithHost(host),
				dockerClient.WithVersion(version),
			)
			if err != nil {
				log.Error().Err(err)
				return
			}
			defer dc.Close()

			containers, err := dc.ContainerList(context.Background(), types.ContainerListOptions{})

		},
	}

	run.Execute()

}
