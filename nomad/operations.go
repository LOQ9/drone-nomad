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

// RegisterJob starts a job task on the nomad cluster
func (d *Driver) RegisterJob(job *api.Job) (*api.JobRegisterResponse, error) {

	// Register the task
	res, _, err := d.client.Jobs().Register(job, nil)
	return res, err
}
