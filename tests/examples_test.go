package tests

import (
	"os"
	"os/exec"
	"testing"

	"github.com/JureBevc/cyanic/actions"
)

func TestPythonFullDeploy(t *testing.T) {
	_, err := os.Stat("./examples")
	if os.IsNotExist(err) {
		os.Chdir("..")
	}
	configPath := "./examples/python-server/cyanic.yaml"
	deployErr := actions.ParseCommand("full-deploy", []string{configPath})
	actions.ParseCommand("remove-staging", []string{configPath})
	actions.ParseCommand("remove-production", []string{configPath})

	if deployErr != nil {
		t.Fatalf("Full deploy failed with error: %s\n", deployErr)
	}
}

func TestNpmFullDeploy(t *testing.T) {
	_, err := os.Stat("./examples")
	if os.IsNotExist(err) {
		os.Chdir("..")
	}
	configPath := "./examples/npm-server/cyanic.yaml"
	deployErr := actions.ParseCommand("full-deploy", []string{configPath})
	actions.ParseCommand("remove-staging", []string{configPath})
	actions.ParseCommand("remove-production", []string{configPath})

	if deployErr != nil {
		t.Fatalf("Full deploy failed with error: %s\n", deployErr)
	}
}

func TestNpmFullDeployMustError(t *testing.T) {
	_, err := os.Stat("./examples")
	if os.IsNotExist(err) {
		os.Chdir("..")
	}

	deployErr := exec.Command("sudo", "go", "run", ".", "full-deploy", "./examples/invalid-cyanic.yaml").Run()

	if deployErr == nil {
		t.Fatalf("Expected full deploy to fail")
	}
}
