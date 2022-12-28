package main

import (
	"net/http"

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
			host := ""
			client := &http.Client{}
			version := ""

			dc, err := dockerClient.NewClient(host, version, client, httpHeaders)
			if err != nil {
				log.Error().Err(err)
				return
			}
			defer dc.Close()

		},
	}

	run.Execute()

}
