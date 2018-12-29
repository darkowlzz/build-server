package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/darkowlzz/build-server/cmd/client/connection"
	pb "github.com/darkowlzz/build-server/pkg/build"
	"github.com/darkowlzz/build-server/pkg/config"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Status of a build",
	Long:  `Status of a build.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("require at least one arg")
		}
		return nil
	},
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

		r, err := client.BuildStatus(ctx, &pb.BuildStatusRequest{
			Id: args[0],
		})
		if err != nil {
			log.Fatalf("could not inspect build container: %v", err)
		}
		fmt.Printf("Build inspect: %s %s\n", r.ContainerID, r.Status)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
