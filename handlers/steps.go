package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

func HandlePortStatus(conf StepConfig) {
	for _, port := range conf.Ports {
		portString := strconv.Itoa(port)
		runningOnPort := isProcessRunningOnPort(portString)
		portStatus := "Offline"
		if runningOnPort {
			portStatus = "Online"
		}
		infoString := fmt.Sprintf("Status of port %s: %s\n", portString, portStatus)
		slog.Info(infoString)
	}
}

func HandleFullDeploy(conf StepConfig) error {
	// Does a deploy to staging, and swap after successful health check
	// If health check for production fails after swap, the swap is reversed
	HandleDeployStaging(conf)
	stagingHealth := HandleHealthCheckStaging(conf)
	if !stagingHealth {
		slog.Error("Aborting full deploy: Failed staging health check")
		return errors.New("Staging health check failed")
	}

	err := HandleSwap(conf)
	if err != nil {
		return err
	}

	productionHealth := HandleHealthCheckProduction(conf)
	if !productionHealth {
		slog.Error("Aborting full deploy and reversing swap: Failed production health check")
		HandleSwap(conf)
		return errors.New("Production health check failed")
	}

	slog.Info("Production health OK. Full deploy finished.")
	return nil
}

func HandleSwap(conf StepConfig) error {
	// Check existing production and staging ports
	slog.Info("Reading existing configurations")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	prodPort := getPortInNginxConfig(conf.Production.UniqueName)

	slog.Info("Existing ports", "staging", stagingPort, "production", prodPort)

	// Continue only if both ports exists, or if only staging exists
	if stagingPort == "" {
		slog.Error("Cannot swap, staging port not found")
		return errors.New("Staging port not found")
	}

	// Overwrite config files of production and staging with swapped ports
	// Create production config with staging port
	slog.Info("Creating production configuration")
	err := createNginxConfig(conf.Production.Nginx, conf.Production.UniqueName, stagingPort)
	if err != nil {
		slog.Error("Could not create configuration for production")
		slog.Error(err.Error())
		return err
	}

	// Create staging config with production port, if port exists
	if prodPort != "" {
		slog.Info("Creating staging configuration")
		err = createNginxConfig(conf.Staging.Nginx, conf.Staging.UniqueName, prodPort)
		if err != nil {
			slog.Error("Could not create configuration for staging")
			slog.Error(err.Error())
			return err
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
		return err
	}

	// Reset nginx
	slog.Info("Restarting nginx config")
	err = restartNginx()
	if err != nil {
		slog.Error("Nginx restart command failed")
		slog.Error(err.Error())
		return err
	}

	return nil
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

func HandleDeployProduction(conf StepConfig) error {
	// Deploy directly to production
	// Get existing production port
	// Shutdown existing process listening on production port
	slog.Info("Killing existing process")
	productionPort := getPortInNginxConfig(conf.Production.UniqueName)
	slog.Info("Existing production port", "value", productionPort)
	if productionPort != "" {
		err := KillProcessOnPort(productionPort)
		if err != nil {
			slog.Error("Could not kill process on port", "value", productionPort)
			slog.Error(err.Error())
		}
	}

	// Get existing staging port
	// Choose any other available port as production port
	slog.Info("Setting deploy port")
	prodPort := getPortInNginxConfig(conf.Staging.UniqueName)

	deployPort := productionPort
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
	err := createNginxConfig(conf.Production.Nginx, conf.Production.UniqueName, deployPort)
	if err != nil {
		slog.Error("Could not create nginx config")
		slog.Error(err.Error())
		return err
	}

	// Run setup to start server on chosen production port
	slog.Info("Running setup")
	runSetup(conf.Setup, deployPort)

	// Restart nginx
	slog.Info("Restarting nginx")
	err = restartNginx()

	return err
}

func HandleHealthCheckStaging(conf StepConfig) bool {
	// Continuously checks the health of staging
	if conf.Staging.HealthCheckUrl == "" {
		slog.Error("Staging health check url not set")
		return false
	}

	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	if stagingPort == "" {
		slog.Error("No active port found for staging")
		return false
	}

	formattedUrl := strings.ReplaceAll(conf.Staging.HealthCheckUrl, "${PORT}", stagingPort)

	health, err := healthCheckPromise(formattedUrl).Await()
	if err != nil {
		slog.Error("Error checking health check:")
		slog.Error(err.Error())
		return false
	}
	if health {
		slog.Info("Staging health check: PASS")
	} else {
		slog.Info("Staging health check: FAIL")
	}
	return health
}

func HandleHealthCheckProduction(conf StepConfig) bool {
	// Continuously checks the health of production
	if conf.Production.HealthCheckUrl == "" {
		slog.Error("Production health check url not set")
		return false
	}

	productionPort := getPortInNginxConfig(conf.Production.UniqueName)
	if productionPort == "" {
		slog.Error("No active port found for production")
		return false
	}

	formattedUrl := strings.ReplaceAll(conf.Production.HealthCheckUrl, "${PORT}", productionPort)

	health, err := healthCheckPromise(formattedUrl).Await()
	if err != nil {
		slog.Error("Error checking health check:")
		slog.Error(err.Error())
		return false
	}
	if health {
		slog.Info("Production health check: PASS")
	} else {
		slog.Info("Production health check: FAIL")
	}
	return health
}

func KillProcessOnPort(portString string) error {
	err := exec.Command("fuser", "-k", "-n", "tcp", portString).Run()
	return err
}
