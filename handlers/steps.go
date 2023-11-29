package handlers

import (
	"errors"
	"log/slog"
	"os/exec"
	"strconv"
)

func HandleFullDeploy(conf StepConfig) {

}

func HandleSwap(conf StepConfig) {
	// Check existing production and staging ports
	slog.Info("Reading existing configurations")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	prodPort := getPortInNginxConfig(conf.Production.UniqueName)

	slog.Info("Existing ports", "staging", stagingPort, "production", prodPort)

	// Continue only if both ports exists, or if only staging exists
	if stagingPort == "" {
		slog.Error("Cannot swap, staging port not found")
		return
	}

	// Overwrite config files of production and staging with swapped ports
	// Create production config with staging port
	slog.Info("Creating production configuration")
	err := createNginxConfig(conf.Production.Nginx, conf.Production.UniqueName, stagingPort)
	if err != nil {
		slog.Error("Could not create configuration for production")
		slog.Error(err.Error())
		return
	}

	// Create staging config with production port, if port exists
	if prodPort != "" {
		slog.Info("Creating staging configuration")
		err = createNginxConfig(conf.Staging.Nginx, conf.Staging.UniqueName, prodPort)
		if err != nil {
			slog.Error("Could not create configuration for staging")
			slog.Error(err.Error())
			return
		}
	} else {
		// Cannot swap with production, delete staging
		// Deleting nginx config
		deleteNginxConfig(conf.Staging.UniqueName)
	}

	// Run nginx test
	slog.Info("Testing nginx config")
	err = testNginx()
	if err != nil {
		slog.Error("Nginx test command failed")
		slog.Error(err.Error())
	}

	// Reset nginx
	slog.Info("Restarting nginx config")
	err = restartNginx()
	if err != nil {
		slog.Error("Nginx restart command failed")
		slog.Error(err.Error())
	}
}

func HandleDeployStaging(conf StepConfig) error {
	// Get existing staging port
	// Shutdown existing process listening on staging port
	slog.Info("Killing existing process")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	slog.Info("Existing staging port", "value", stagingPort)
	if stagingPort != "" {
		err := KillProcessOnPort(stagingPort)
		if err != nil {
			slog.Error("Could not kill process on port", "value", stagingPort)
			slog.Error(err.Error())
		}
	}

	// Get existing production port
	// Choose any other available port as staging port
	slog.Info("Setting deploy port")
	prodPort := getPortInNginxConfig(conf.Production.UniqueName)

	deployPort := stagingPort
	if deployPort == "" {
		for _, port := range conf.Ports {
			portStr := strconv.Itoa(port)
			if portStr != prodPort {
				deployPort = portStr
				break
			}
		}
	}

	if deployPort == "" {
		return errors.New("Could not define a valid port for deployment")
	}

	slog.Info("Chose port", "value", deployPort)
	// Create nginx config
	slog.Info("Creating nginx config")
	err := createNginxConfig(conf.Staging.Nginx, conf.Staging.UniqueName, deployPort)
	if err != nil {
		slog.Error("Could not create nginx config")
		slog.Error(err.Error())
		return err
	}

	// Run setup to start server on chosen staging port
	slog.Info("Running setup")
	runSetup(conf.Setup, deployPort)

	// Restart nginx
	slog.Info("Restarting nginx")
	err = restartNginx()

	return err
}

func HandleRemoveStaging(conf StepConfig) error {
	slog.Info("Killing existing process")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	slog.Info("Existing staging port", "value", stagingPort)
	if stagingPort != "" {
		err := KillProcessOnPort(stagingPort)
		if err != nil {
			slog.Error("Could not kill process on port", "value", stagingPort)
			slog.Error(err.Error())
		}
	}

	// Remove config
	deleteNginxConfig(conf.Staging.UniqueName)

	// Restart nginx
	slog.Info("Restarting nginx")
	err := restartNginx()

	return err
}

func HandleRemoveProduction(conf StepConfig) error {
	slog.Info("Killing existing process...")
	prodPort := getPortInNginxConfig(conf.Production.UniqueName)
	slog.Info("Existing staging port", "value", prodPort)
	if prodPort != "" {
		err := KillProcessOnPort(prodPort)
		if err != nil {
			slog.Error("Could not kill process on port", "value", prodPort)
			slog.Error(err.Error())
		}
	}

	// Remove config
	deleteNginxConfig(conf.Production.UniqueName)

	// Restart nginx
	slog.Info("Restarting nginx")
	err := restartNginx()

	return err
}

func HandleDeployProduction(conf StepConfig) {

}

func HandleHealthCheckStaging(conf StepConfig) {

}

func HandleHealthCheckProduction(conf StepConfig) {

}

func KillProcessOnPort(portString string) error {
	err := exec.Command("sudo", "fuser", "-k", "-n", "tcp", portString).Run()
	return err
}
