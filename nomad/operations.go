package nomad

import (
	"github.com/hashicorp/nomad/api"
)

// ParseTemplate ...
func (d *Driver) ParseTemplate(jobHCL string) (*api.Job, error) {
	job, err := d.client.Jobs().ParseHCL(jobHCL, true)

	if err != nil {
		return &api.Job{}, err
	}

	return job, nil
}

// StartJob starts a job task on the nomad cluster
func (d *Driver) StartJob(job *api.Job) error {

	// Register the task
	_, _, err := d.client.Jobs().Register(job, nil)
	return err
}
