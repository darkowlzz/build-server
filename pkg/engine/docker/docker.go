package docker

import (
	"context"
	"log"
	"path/filepath"

	"github.com/fsouza/go-dockerclient"
)

// Client stores docker daemon unix domain socket endpoint path and a docker
// client.
type Client struct {
	endpoint string
	client   *docker.Client
}

// NewClient creates a new Client, given a docker daemon socket endpoint path.
func NewClient(endpoint string) (*Client, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		endpoint: endpoint,
		client:   client,
	}, nil
}

// GetInfo returns the docker daemon information.
func (c Client) GetInfo() (name string, version string, err error) {
	info, err := c.client.Info()
	if err != nil {
		return "", "", err
	}

	return "docker", info.ServerVersion, nil
}

// InspectContainer returns the status of a container given an ID.
func (c Client) InspectContainer(ctx context.Context, id string) (containerID, status string, err error) {
	i, err := c.client.InspectContainerWithContext(id, ctx)
	if err != nil {
		return "", "", err
	}
	return i.ID, i.State.Status, nil
}

// StartBuild creates a container with the given image, starts the container and
// runs the build. It mounts the source as a point mount.
func (c Client) StartBuild(ctx context.Context, name string, buildRoot string, command []string, image string, mountPath string) (id string, err error) {
	configOpts := docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image: image,
			// AttachStdout: true,
			// AttachStderr: true,
			Cmd:        command,
			WorkingDir: mountPath,
		},
		HostConfig: &docker.HostConfig{
			// AutoRemove: true,
			Mounts: []docker.HostMount{
				{
					Source: filepath.Join(buildRoot, name),
					Target: mountPath,
					Type:   "bind",
					BindOptions: &docker.BindOptions{
						Propagation: "rprivate",
					},
				},
			},
		},
		Context: ctx,
	}

	container, err := c.client.CreateContainer(configOpts)
	if err != nil {
		return "", err
	}

	log.Printf("Starting container %s", container.ID)
	if err := c.client.StartContainerWithContext(container.ID, nil, ctx); err != nil {
		return "", err
	}

	return container.ID, nil
}
