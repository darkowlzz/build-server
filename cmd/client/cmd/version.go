package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/darkowlzz/build-server/cmd/client/connection"
	pb "github.com/darkowlzz/build-server/pkg/build"
	"github.com/darkowlzz/build-server/pkg/config"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the build server",
	Long:  `Version of the build server.`,
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

		r, err := client.GetInfo(ctx, &pb.InfoRequest{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("SERVER INFO: %s %s\n", r.Name, r.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
