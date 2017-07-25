package main

import (
	"os"
	"testing"

	"gopkg.in/urfave/cli.v1"
)

func TestInitApp(t *testing.T) {

	testEnvs := []struct {
		envVars    map[string]string
		shouldFail bool
		reason     string
	}{
		{
			envVars: map[string]string{
				"PLUGIN_STORAGE_URL":    "s3://test-bucket/prefix",
				"PLUGIN_REPO_URL":       "http://charts.example.com",
				"PLUGIN_AWS_REGION":     "ap-southeast-1",
				"PLUGIN_AWS_ACCESS_KEY": "TESTACCESSKEY",
				"PLUGIN_AWS_SECRET_KEY": "TESTSECRETKEY",
			},
			shouldFail: false,
			reason:     "All env variables provided",
		},
		{
			envVars: map[string]string{
				"PLUGIN_REPO_URL":   "http://charts.example.com",
				"PLUGIN_AWS_REGION": "ap-southeast-1",
			},
			shouldFail: true,
			reason:     "storage-url is missing",
		},
		{
			envVars: map[string]string{
				"PLUGIN_REPO_URL":    "http://charts.example.com",
				"PLUGIN_STORAGE_URL": "s3://test-bucket/prefix",
			},
			shouldFail: true,
			reason:     "aws-region is missing",
		},
	}

	for _, testEnv := range testEnvs {
		for envvar := range testEnv.envVars {
			os.Setenv(envvar, testEnv.envVars[envvar])
		}
		app := initApp(runEnvCheck)
		app.Run([]string{""})
		if envCheckFailed != testEnv.shouldFail {
			if testEnv.shouldFail {
				t.Errorf("Expected envTest to fail because %s - \nEnv: %v", testEnv.reason, testEnv.envVars)
			} else {
				t.Errorf("Expected envTest to pass - \nEnv: %v", testEnv.envVars)
			}
		}
		// clean up
		for envvar := range testEnv.envVars {
			os.Unsetenv(envvar)
		}
	}

}

// Do not communicate by sharing memory; instead, share memory by communicating.

// but here we will share memory :scream:...
var envCheckFailed bool

func runEnvCheck(c *cli.Context) error {
	err := validateConfig(configFromEnv(c))
	envCheckFailed = (err != nil)
	//if runAction returns error, os.Exit(1) causes tests to always fail
	return nil
}
