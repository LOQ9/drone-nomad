package main

import (
	"drone-nomad/nomad"
	"fmt"
	"io/ioutil"

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
		return err
	}

	// Parse template
	nomadTemplate, err := nomad.ParseTemplate(string(nomadTemplateFile))

	if err != nil {
		return err
	}

	// Perform substitions

	// Launch deployment
	nomadJob, err := nomad.RegisterJob(nomadTemplate)

	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("%+v\n", nomadJob)

	return nil
}

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
