package temporal

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/client"
)

type Client struct {
	client    client.Client
	namespace string
}

func NewClient(hostPort, namespace string, otelEnabled bool) (*Client, error) {
	options := client.Options{
		HostPort:  hostPort,
		Namespace: namespace,
	}

	c, err := client.Dial(options)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporal client: %w", err)
	}

	return &Client{
		client:    c,
		namespace: namespace,
	}, nil
}

func (c *Client) GetClient() client.Client {
	return c.client
}

func (c *Client) GetNamespace() string {
	return c.namespace
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error) {
	return c.client.ExecuteWorkflow(ctx, options, workflow, args...)
}