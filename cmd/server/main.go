package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"path/filepath"

	pb "github.com/darkowlzz/build-server/pkg/build"
	dockerBuild "github.com/darkowlzz/build-server/pkg/engine/docker"
	"github.com/darkowlzz/build-server/pkg/util"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port                  = ":50051"
	defaultBuildRoot      = "/tmp/remote-build"
	defaultArtifactsDir   = "out"
	defaultDockerEndpoint = "unix:///var/run/docker.sock"
)

type server struct {
	dockerClient *dockerBuild.Client
	buildRoot    string
}

func (s *server) GetInfo(ctx context.Context, in *pb.InfoRequest) (*pb.InfoReply, error) {
	return &pb.InfoReply{
		Name:    "Build Server",
		Version: "0.0.1",
	}, nil
}

func (s *server) GetEngineInfo(ctx context.Context, in *pb.EngineInfoRequest) (*pb.EngineInfoReply, error) {
	name, version, err := s.dockerClient.GetInfo()
	if err != nil {
		return nil, err
	}
	return &pb.EngineInfoReply{
		Name:    name,
		Version: version,
	}, nil
}

func (s *server) BuildStatus(ctx context.Context, in *pb.BuildStatusRequest) (*pb.BuildStatusReply, error) {
	id, status, err := s.dockerClient.InspectContainer(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.BuildStatusReply{
		Id:          in.Id,
		ContainerID: id,
		Status:      status,
	}, nil
}

func (s *server) GetArtifacts(ctx context.Context, in *pb.GetArtifactsRequest) (*pb.GetArtifactsReply, error) {
	target := filepath.Join(s.buildRoot, in.GetId(), defaultArtifactsDir)

	if _, err := os.Stat(target); err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := util.Pack(target, &buf, []string{}); err != nil {
		return nil, err
	}

	return &pb.GetArtifactsReply{
		Artifacts: buf.Bytes(),
	}, nil
}

func (s *server) StartBuild(ctx context.Context, in *pb.StartBuildRequest) (*pb.StartBuildReply, error) {
	// Generate a build ID.
	uid, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	destination := filepath.Join(s.buildRoot, uid.String())

	// Create a build directory.
	if err := os.MkdirAll(filepath.Join(s.buildRoot, uid.String()), 0755); err != nil {
		return nil, err
	}

	buildCtx := in.GetBuildCtx()

	if err := util.Unpack(destination, bytes.NewBuffer(buildCtx)); err != nil {
		return nil, err
	}

	// Create the out artifact directory.
	if err := os.MkdirAll(filepath.Join(destination, defaultArtifactsDir), 0755); err != nil {
		return nil, err
	}

	_, err = s.dockerClient.StartBuild(ctx, uid.String(), s.buildRoot, in.GetCommand(), in.GetImage(), in.GetMountPath())
	if err != nil {
		return nil, err
	}
	return &pb.StartBuildReply{Id: uid.String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.MaxMsgSize(1024*1024*50), grpc.MaxRecvMsgSize(1024*1024*50))

	client, err := dockerBuild.NewClient(defaultDockerEndpoint)
	if err != nil {
		log.Fatalf("failed to connect to docker daemon: %v", err)
	}

	// Create build root directory.
	if err := os.MkdirAll(defaultBuildRoot, 0755); err != nil {
		log.Fatalf("failed to create build root: %v", err)
	}

	pb.RegisterBuildServer(s, &server{dockerClient: client, buildRoot: defaultBuildRoot})

	reflection.Register(s)

	log.Println("Build Server listening at", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
