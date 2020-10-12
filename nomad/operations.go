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
func (d *Driver) RegisterJob(job *api.Job, preserveCounts bool) (*api.JobRegisterResponse, error) {

	// Register the task
	res, _, err := d.client.Jobs().RegisterOpts(job,
		&api.RegisterOptions{
			PreserveCounts: preserveCounts,
		}, nil)
	return res, err
}

// PlanJob plans a job on the nomad cluster
func (d *Driver) PlanJob(job *api.Job) (*api.JobPlanResponse, error) {

	// Register the task
	res, _, err := d.client.Jobs().Plan(job, true, nil)
	return res, err
}
