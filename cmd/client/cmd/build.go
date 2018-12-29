package cmd

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/darkowlzz/build-server/cmd/client/connection"
	pb "github.com/darkowlzz/build-server/pkg/build"
	"github.com/darkowlzz/build-server/pkg/config"
	"github.com/darkowlzz/build-server/pkg/util"
	"github.com/spf13/cobra"
)

const (
	defaultArtifactsDir = "out"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the current project",
	Long:  `Build the current project.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := config.ReadConfig()
		if err != nil {
			log.Fatal(err)
		}

		client, err := connection.GetBuildClient(*config)
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var buf bytes.Buffer

		absPath, err := filepath.Abs(".")
		if err != nil {
			log.Fatal(err)
		}

		if err := util.Pack(absPath, &buf, []string{defaultArtifactsDir}); err != nil {
			log.Fatalf("error file walking the tree: %v", err)
		}

		log.Println("Context size:", len(buf.Bytes()))
		log.Println("MOUNTPATH:", config.MountPath)
		r, err := client.StartBuild(ctx, &pb.StartBuildRequest{
			Image:     config.Image,
			Command:   []string{"/bin/sh", "-c", config.Command},
			BuildCtx:  buf.Bytes(),
			MountPath: config.MountPath,
		})
		if err != nil {
			log.Fatalf("could not start build: %v", err)
		}
		fmt.Printf("Build started with id: %s\n", r.GetId())
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
