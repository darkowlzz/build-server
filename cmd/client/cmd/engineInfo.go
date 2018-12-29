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

// engineInfoCmd represents the engineInfo command
var engineInfoCmd = &cobra.Command{
	Use:   "engineInfo",
	Short: "EngineInfo of the build engine",
	Long:  `EngineInfo of the build engine.`,
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

		r, err := client.GetEngineInfo(ctx, &pb.EngineInfoRequest{})
		if err != nil {
			log.Fatalf("could not get engine info: %v", err)
		}
		fmt.Printf("ENGINE INFO: %s %s\n", r.Name, r.Version)
	},
}

func init() {
	rootCmd.AddCommand(engineInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// engineInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// engineInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
