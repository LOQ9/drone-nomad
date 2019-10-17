package main

import (
	"drone-nomad/nomad"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/drone/envsubst"
)

type (
	// Repo contains information related to the repository
	Repo struct {
		Owner string
		Name  string
	}

	// Build contains information related to the build
	Build struct {
		Tag     string
		Event   string
		Number  int
		Commit  string
		Ref     string
		Branch  string
		Author  string
		Status  string
		Link    string
		Started int64
		Created int64
	}

	// Job ...
	Job struct {
		Started int64
	}

	// Config ...
	Config struct {
		Address  string
		Token    string
		Region   string
		Template string
	}

	// Plugin ...
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
		Job    Job
	}
)

// Exec initiates the plugin execution
func (p Plugin) Exec() error {

	// Connect to Nomad
	nomad, err := nomad.New(&nomad.Client{
		URL: p.Config.Address,
	})

	if err != nil {
		return err
	}

	// Read Template File
	nomadTemplateFile, err := ioutil.ReadFile(p.Config.Template)
	if err != nil {
		return fmt.Errorf("Could not read nomad template file")
	}

	for _, pair := range os.Environ() {
		fmt.Println(pair)
	}

	// Perform substitions
	//nomadTemplateSubst, err := envsubst.EvalEnv(
	//	string(nomadTemplateFile),
	//)
	nomadTemplateSubst, err := p.replaceEnv(
		string(nomadTemplateFile),
	)

	if err != nil {
		return err
	}

	// Parse template
	nomadTemplate, err := nomad.ParseTemplate(nomadTemplateSubst)

	if err != nil {
		return err
	}

	// Plan deployment
	_, err = nomad.PlanJob(nomadTemplate)

	if err != nil {
		return err
	}

	// Launch deployment
	nomadJob, err := nomad.RegisterJob(nomadTemplate)

	if err != nil {
		return err
	}

	if len(nomadJob.Warnings) > 0 {
		fmt.Printf("Nomad job deployed with %d warning(s)\n", len(nomadJob.Warnings))
		fmt.Printf("%s\n", nomadJob.Warnings)
	} else {
		fmt.Printf("Nomad job deployed successfuly!\n")
	}

	return nil
}

// replaceEnv env changes vars from template
func (p Plugin) replaceEnv(template string) (string, error) {
	names := map[string]bool{}

	template, err := envsubst.Eval(template, func(in string) string {
		names[in] = true
		return in
	})

	if err != nil {
		return "", err
	}

	return template, nil
}
