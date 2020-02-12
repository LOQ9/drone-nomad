package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/ulule/deepcopier"
)

// Driver is the abstraction for using the nomad driver
type Driver struct {
	client *api.Client
}

// New creates a new instance of the nomad driver
func New(c *Client) (*Driver, error) {

	// Create a new nomad client config
	config := api.DefaultConfig()
	config.Address = c.Address

	if c.Region != "" {
		config.Region = c.Region
	}

	if c.Token != "" {
		config.SecretID = c.Token
	}

	if c.Namespace != "" {
		config.Namespace = c.Namespace
	}

	if c.TLSConfig != nil {
		deepcopier.Copy(c.TLSConfig).To(config.TLSConfig)
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Driver{client}, nil
}
