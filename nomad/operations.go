package nomad

import (
	"errors"
	"fmt"
	"github.com/hashicorp/nomad/api"
	"time"
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

func (d *Driver) WatchDeployment(job *api.JobRegisterResponse, timeout time.Duration) error {
	// get the deployment id for this alloc
	eval, _, err := d.client.Evaluations().Info(job.EvalID, nil)
	if err != nil {
		return err
	}

	if eval.DeploymentID == "" {
		// sometimes Nomad initially returns eval
		// info with an empty deploymentID; and a retry is required in order to get the
		// updated response from Nomad.
		evalInfoTimeout := time.NewTicker(time.Second * 60)
		defer evalInfoTimeout.Stop()
		for {
			select {
			case <-evalInfoTimeout.C:
				return errors.New("timeout reached on attempting to find deployment ID")
			default:
				if eval, _, err = d.client.Evaluations().Info(eval.ID, nil); err != nil {
					return err
				}
				if eval.DeploymentID != "" {
					break
				}
				time.Sleep(time.Second * 2)
				continue
			}
		}
	}

	timer := time.NewTimer(timeout)
	for {
		select {
		case <-timer.C:
			return fmt.Errorf("deployment watcher timedout. Deployment: %s", eval.DeploymentID)
		default:
			deployment, _, err := d.client.Deployments().Info(eval.DeploymentID, nil)
			if err != nil {
				return err
			}
			switch deployment.Status {
			case "successful":
				return nil
			case "running":
				time.Sleep(time.Second * 5)
			default:
				return fmt.Errorf("deployment: %s failed. Status: %s : %s", eval.DeploymentID, deployment.Status, deployment.StatusDescription)
			}
		}
	}
}
