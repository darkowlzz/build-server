package connection

import (
	"google.golang.org/grpc"

	pb "github.com/darkowlzz/build-server/pkg/build"
	"github.com/darkowlzz/build-server/pkg/config"
)

// GetBuildClient returns a Buld gRPC client.
func GetBuildClient(config config.BuildConfig) (pb.BuildClient, error) {
	conn, err := grpc.Dial(
		config.Remote,
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*50)), // Avoid this using Stream API
	)
	if err != nil {
		return nil, err
	}

	return pb.NewBuildClient(conn), nil
}
