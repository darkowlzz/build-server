package cmd

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/darkowlzz/build-server/cmd/client/connection"
	pb "github.com/darkowlzz/build-server/pkg/build"
	"github.com/darkowlzz/build-server/pkg/config"
	"github.com/darkowlzz/build-server/pkg/util"
	"github.com/spf13/cobra"
)

// artifactsCmd represents the artifacts command
var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "Fetch the build artifacts",
	Long:  `Fetch the build artifacts.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}

		client, err := connection.GetBuildClient(*config)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		r, err := client.GetArtifacts(ctx, &pb.GetArtifactsRequest{
			Id: args[0],
		})
		if err != nil {
			log.Fatalf("could not fetch the build artifacts: %v", err)
		}
		b := r.GetArtifacts()

		if err := util.Unpack(defaultArtifactsDir, bytes.NewBuffer(b)); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(artifactsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// artifactsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// artifactsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
