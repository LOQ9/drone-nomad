package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/drone/envsubst"

	"drone-nomad/nomad"
)

const (
	// SanitizeSuffix used when matching vars to sanitize
	SanitizeSuffix = "_SANITIZE"
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
		Parent  int    `json:"parent" env:"DRONE_BUILD_PARENT"`
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
		Address                string        `json:"address" env:"PLUGIN_ADDR"`
		Token                  string        `json:"token" env:"PLUGIN_TOKEN"`
		ConsulToken            string        `json:"consul_token" env:"PLUGIN_CONSUL_TOKEN"`
		VaultToken             string        `json:"vault_token" env:"PLUGIN_VAULT_TOKEN"`
		Region                 string        `json:"region" env:"PLUGIN_REGION"`
		Namespace              string        `json:"namespace" env:"PLUGIN_NAMESPACE"`
		Template               string        `json:"template" env:"PLUGIN_TEMPLATE"`
		TLSCACert              string        `json:"tls_ca_cert" env:"PLUGIN_TLS_CA_CERT"`
		TLSCACertPem           string        `json:"tls_ca_cert_pem" env:"PLUGIN_TLS_CA_CERT_PEM"`
		TLSCAPath              string        `json:"tls_ca_path" env:"PLUGIN_TLS_CA_PATH"`
		TLSClientCert          string        `json:"tls_client_cert" env:"PLUGIN_TLS_CLIENT_CERT"`
		TLSClientCertPem       string        `json:"tls_client_cert_pem" env:"PLUGIN_TLS_CLIENT_CERT_PEM"`
		TLSClientKey           string        `json:"tls_client_key" env:"PLUGIN_TLS_CLIENT_KEY"`
		TLSClientKeyPem        string        `json:"tls_client_key_pem" env:"PLUGIN_TLS_CLIENT_KEY_PEM"`
		TLSServerName          string        `json:"tls_servername" env:"PLUGIN_TLS_SERVERNAME"`
		TLSInsecure            bool          `json:"tls_insecure" env:"PLUGIN_TLS_INSECURE"`
		PreserveCounts         bool          `json:"preserve_counts" env:"PLUGIN_PRESERVE_COUNTS"`
		Debug                  bool          `json:"debug" env:"PLUGIN_DEBUG"`
		DryRun                 bool          `json:"dry_run" env:"PLUGIN_DRY_RUN"`
		WatchDeployment        bool          `json:"watch_deployment" env:"PLUGIN_WATCH_DEPLOYMENT"`
		WatchDeploymentTimeout time.Duration `json:"watch_deployment_timeout" env:"PLUGIN_WATCH_DEPLOYMENT_TIMEOUT"`
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

	nomadClient := &nomad.Client{
		Address:   p.Config.Address,
		Region:    p.Config.Region,
		Namespace: p.Config.Namespace,
		Token:     p.Config.Token,
		TLSConfig: &nomad.ClientTLSConfig{
			CACert:        p.Config.TLSCACert,
			CACertPEM:     []byte(p.Config.TLSCACertPem),
			CAPath:        p.Config.TLSCAPath,
			ClientCert:    p.Config.TLSClientCert,
			ClientCertPEM: []byte(p.Config.TLSClientCertPem),
			ClientKey:     p.Config.TLSClientKey,
			ClientKeyPEM:  []byte(p.Config.TLSClientKeyPem),
			TLSServerName: p.Config.TLSServerName,
			Insecure:      p.Config.TLSInsecure,
		},
	}

	// Connect to Nomad
	nomad, err := nomad.New(nomadClient)

	if err != nil {
		return err
	}

	// Read Template File
	nomadTemplateFile, err := ioutil.ReadFile(p.Config.Template)
	if err != nil {
		return fmt.Errorf("Could not read nomad template file")
	}

	// Perform substitution
	nomadTemplateSubst := p.replaceEnv(string(nomadTemplateFile))

	// Parse template
	nomadTemplate, err := nomad.ParseTemplate(nomadTemplateSubst)
	if err != nil {
		return err
	}

	// Set consul ACL token
	if p.Config.ConsulToken != "" {
		*nomadTemplate.ConsulToken = p.Config.ConsulToken
	}

	// Set vault ACL token
	if p.Config.VaultToken != "" {
		*nomadTemplate.VaultToken = p.Config.VaultToken
	}

	// Log template to STDOUT when debugging is enabled
	if p.Config.Debug {
		fmt.Println("Nomad template:")
		fmt.Println(nomadTemplateSubst)
	}

	// Plan deployment
	_, err = nomad.PlanJob(nomadTemplate)
	if err != nil {
		return err
	}

	if !p.Config.DryRun {
		// Launch deployment
		nomadJob, err := nomad.RegisterJob(nomadTemplate, p.Config.PreserveCounts)
		if err != nil {
			return err
		}

		if len(nomadJob.Warnings) > 0 {
			fmt.Printf("Nomad job deployed with %d warning(s)\n", len(nomadJob.Warnings))
			fmt.Printf("%s\n", nomadJob.Warnings)
		} else if p.Config.WatchDeployment {
			if p.Config.WatchDeploymentTimeout == 0 {
				p.Config.WatchDeploymentTimeout = time.Minute * 5
			}
			// block and watch deployment
			if err = nomad.WatchDeployment(nomadJob, p.Config.WatchDeploymentTimeout); err != nil {
				return err
			}
		} else {
			fmt.Printf("Nomad job deployed successfuly!\n")
		}
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
				structVal = append(structVal, fieldTag)
			}
		}
	}

	return structVal
}

// replaceEnv changes vars from template
func (p Plugin) replaceEnv(template string) string {

	// Get current passed vars
	templateVars := p.envMap()

	// Regular expression matching var expression ${...}
	reVars := regexp.MustCompile(`(?m)\$\{(.+?)\}`)

	// Replace all matches of regular expression and check if it can be replaced
	template = reVars.ReplaceAllStringFunc(template, func(s string) string {
		// Find the exact var name inside match, e.g. ${DRONE_TAG=latest} becomes "DRONE_TAG=latest"
		matches := reVars.FindStringSubmatch(s)

		// Check string sufix
		if strings.HasSuffix(matches[1], SanitizeSuffix) {

			// Remove Suffix
			envName := strings.TrimSuffix(matches[1], SanitizeSuffix)

			// Get var content
			envValue := os.Getenv(envName)

			r, _ := regexp.Compile(`[^a-z0-9]`)
			replacedString := r.ReplaceAllString(strings.ToLower(envValue), "-")

			if p.Config.Debug {
				fmt.Printf("Replacing (var: %s, value: %s) to %s\n", envName, strings.ToLower(envValue), replacedString)
			}

			return replacedString
		}

		// Loop over our known template vars if they can be replaced, otherwise return the original string
		for i := range templateVars {

			// Check if the found var starts with one of our known vars
			if strings.Index(matches[1], templateVars[i]) == 0 {
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
