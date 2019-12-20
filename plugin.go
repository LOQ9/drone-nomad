package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	"github.com/drone/envsubst"

	"drone-nomad/nomad"
)

type (
	// Repo contains information related to the repository
	Repo struct {
		Owner string `json:"owner" env:"DRONE_REPO_OWNER"`
		Name  string `json:"name" env:"DRONE_REPO_NAME"`
	}

	// Build contains information related to the build
	Build struct {
		Tag     string `json:"tag" env:"DRONE_TAG"`
		Event   string `json:"event" env:"DRONE_BUILD_EVENT"`
		Number  int    `json:"number" env:"DRONE_BUILD_NUMBER"`
		Commit  string `json:"commit" env:"DRONE_COMMIT_SHA"`
		Ref     string `json:"ref" env:"DRONE_COMMIT_REF"`
		Branch  string `json:"branch" env:"DRONE_COMMIT_BRANCH"`
		Author  string `json:"author" env:"DRONE_COMMIT_AUTHOR"`
		Message string `json:"message" env:"DRONE_COMMIT_MESSAGE"`
		Status  string `json:"status" env:"DRONE_BUILD_STATUS"`
		Link    string `json:"link" env:"DRONE_BUILD_LINK"`
		Started int64  `json:"started" env:"DRONE_BUILD_STARTED"`
		Created int64  `json:"created" env:"DRONE_BUILD_CREATED"`
	}

	// Job ...
	Job struct {
		Started int64 `json:"created" env:"DRONE_JOB_STARTED"`
	}

	// Config ...
	Config struct {
		Address   string `json:"address" env:"PLUGIN_ADDR"`
		Token     string `json:"token" env:"PLUGIN_TOKEN"`
		Region    string `json:"region" env:"PLUGIN_REGION"`
		Namespace string `json:"namespace" env:"PLUGIN_NAMESPACE"`
		Template  string `json:"template" env:"PLUGIN_TEMPLATE"`
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
		Address:   p.Config.Address,
		Region:    p.Config.Region,
		Namespace: p.Config.Namespace,
		Token:     p.Config.Token,
	})

	if err != nil {
		return err
	}

	// Read Template File
	nomadTemplateFile, err := ioutil.ReadFile(p.Config.Template)
	if err != nil {
		return fmt.Errorf("Could not read nomad template file")
	}

	// Perform substitions
	nomadTemplateSubst := p.replaceEnv(string(nomadTemplateFile))

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

func (p Plugin) envMap() []string {
	structVal := make([]string, 0)
	u := reflect.ValueOf(p)

	for k := 0; k < u.NumField(); k++ {
		field := u.Field(k)
		fieldType := field.Type()

		for j := 0; j < fieldType.NumField(); j++ {
			nestedField := fieldType.Field(j)
			fieldName := nestedField.Name
			fieldTag := nestedField.Tag.Get("env")

			f := reflect.Indirect(field).FieldByName(fieldName)
			if f.IsValid() {
				structVal = append(structVal,fieldTag)
			}
		}
	}

	return structVal
}

// replaceEnv changes vars from template
func (p Plugin) replaceEnv(template string) string {

	// Get current passed vars
	templateVars := p.envMap()

	// regular expression matching var expression ${...}
	reVars := regexp.MustCompile(`(?m)\$\{(.+?)\}`)

	// replace all matches of regular expression and check if it can be replaced
	template = reVars.ReplaceAllStringFunc(template, func(s string) string {
		// find the exact var name inside match, e.g. ${DRONE_TAG=latest} becomes "DRONE_TAG=latest"
		matches := reVars.FindStringSubmatch(s)
		// loop over our known template vars if the can be replaced, otherwise return the original string
		for i := range templateVars {
			// check if the found var starts with one of our known vars
			if strings.Index(matches[1], templateVars[i])==0 {
				subst, err := envsubst.EvalEnv(s)
				if err != nil {
					return s
				}
				return subst
			}
		}

		return s
	})

	return template
}
