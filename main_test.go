package main

import (
	"os"
	"testing"

	"gopkg.in/urfave/cli.v1"
)

func TestInitApp(t *testing.T) {

	testEnvs := []struct {
		envVars map[string]string
		pass    bool
	}{
		{
			envVars: map[string]string{
				"PLUGIN_STORAGE_URL":    "s3://test-bucket/prefix",
				"PLUGIN_REPO_URL":       "http://charts.example.com",
				"PLUGIN_AWS_REGION":     "ap-southeast-1",
				"PLUGIN_AWS_ACCESS_KEY": "TESTACCESSKEY",
				"PLUGIN_AWS_SECRET_KEY": "TESTSECRETKEY",
			},
			pass: true,
		},
		{
			envVars: map[string]string{
				"PLUGIN_REPO_URL":       "http://charts.example.com",
				"PLUGIN_AWS_REGION":     "ap-southeast-1",
				"PLUGIN_AWS_ACCESS_KEY": "TESTACCESSKEY",
				"PLUGIN_AWS_SECRET_KEY": "TESTSECRETKEY",
			},
			pass: false, //missing storage url
		},
	}

	for _, testEnv := range testEnvs {
		for envvar := range testEnv.envVars {
			os.Setenv(envvar, testEnv.envVars[envvar])
		}
		app := initApp(runEnvCheck)
		err := app.Run([]string{""})
		if (err == nil) != testEnv.pass {
			if testEnv.pass {
				t.Errorf("Expected %v to pass", testEnv.envVars)
			} else {
				t.Errorf("Expected %v to fail", testEnv.envVars)
			}
		}
		// clean up
		for envvar := range testEnv.envVars {
			os.Unsetenv(envvar)
		}
	}

}

func runEnvCheck(c *cli.Context) error {
	conf := Config{
		SourceDir:    c.String("source-dir"),
		Exclude:      c.StringSlice("exclude"),
		StorageURL:   c.String("storage-url"),
		RepoURL:      c.String("repo-url"),
		Debug:        c.Bool("debug"),
		AWSAccessKey: c.String("aws-acces-key"),
		AWSSecretKey: c.String("aws-secret-key"),
		AWSRegion:    c.String("aws-region"),
	}
	return validateConfig(conf)
}
