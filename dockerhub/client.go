package dockerhub

import (
	"context"
	"fmt"

	"github.com/mittwald/goharbor-client/v3/apiv1"
)

type Client struct {
	restClient *apiv1.RESTClient
}

func NewClient(serverAddr, user, password string) (*Client, error) {
	registryClient, err := apiv1.NewRESTClientForHost(
		serverAddr + "/api", // TODO need carefully handle the trailing slash
		user,
		password,
	)

	if err != nil{
		return nil, err
	}

	c := new(Client)
	c.restClient = registryClient

	return c, nil
}

// Auth dockerhub
// TODO another way to verify the availability of credentials? just like mysql ping?
func (c *Client) Auth(ctx context.Context) error {
	_, err := c.restClient.ListProjects(ctx, "")
	if err != nil {
		return fmt.Errorf("docker registry auth error %v", err)
	}

	return nil
}
