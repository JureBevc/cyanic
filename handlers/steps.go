package handlers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func HandleFullDeploy(conf StepConfig) {

}

func HandleSwap(conf StepConfig) {
	// Check existing production and staging ports
	fmt.Println("Getting existing ports...")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	prodPort := getPortInNginxConfig(conf.Production.UniqueName)

	fmt.Printf("Staging: %s, production: %s\n", stagingPort, prodPort)

	// Continue only if both ports exists, or if only staging exists
	if stagingPort == "" {
		fmt.Println("Error: Cannot swap, staging port not found")
		return
	}

	// Overwrite config files of production and staging with swapped ports
	// Create production config with staging port
	fmt.Println("Creating production nginx config...")
	err := createNginxConfig(conf.Production.Nginx, conf.Production.UniqueName, stagingPort)
	if err != nil {
		fmt.Println("Error: Could not create nginx config for production")
		fmt.Println(err)
		return
	}

	// Create staging config with production port, if port exists
	if prodPort != "" {
		fmt.Println("Creating staging nginx config...")
		err = createNginxConfig(conf.Staging.Nginx, conf.Staging.UniqueName, prodPort)
		if err != nil {
			fmt.Println("Error: Could not create nginx config for staging")
			fmt.Println(err)
			return
		}
	} else {
		// Cannot swap with production, delete staging
		// Deleting nginx config
		deleteNginxConfig(conf.Staging.UniqueName)
	}

	// Run nginx test
	fmt.Println("Testing nginx config...")
	err = testNginx()
	if err != nil {
		fmt.Println("Error: Nginx test command failed")
		fmt.Println(err)
	}

	// Reset nginx
	fmt.Println("Restarting nginx config...")
	err = restartNginx()
	if err != nil {
		fmt.Println("Error: Nginx restart command failed")
		fmt.Println(err)
	}
}

func HandleDeployStaging(conf StepConfig) error {
	// Get existing staging port
	// Shutdown existing process listening on staging port
	fmt.Println("Killing existing process...")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	fmt.Printf("Existing staging port: %s\n", stagingPort)
	if stagingPort != "" {
		err := KillProcessOnPort(stagingPort)
		if err != nil {
			fmt.Printf("Could not kill process on port %s\n", stagingPort)
			fmt.Println(err)
		}
	}

	// Get existing production port
	// Choose any other available port as staging port
	fmt.Println("Setting deploy port...")
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

	fmt.Printf("Chose port %s\n", deployPort)
	// Create nginx config
	fmt.Println("Creating nginx config...")
	err := createNginxConfig(conf.Staging.Nginx, conf.Staging.UniqueName, deployPort)
	if err != nil {
		fmt.Println("Error: Could not create nginx config")
		fmt.Println(err)
		return err
	}

	// Run setup to start server on chosen staging port
	fmt.Println("Running setup...")
	runSetup(conf.Setup, deployPort)

	// Restart nginx
	fmt.Println("Restarting nginx...")
	err = restartNginx()

	return err
}

func HandleRemoveStaging(conf StepConfig) error {
	fmt.Println("Killing existing process...")
	stagingPort := getPortInNginxConfig(conf.Staging.UniqueName)
	fmt.Printf("Existing staging port: %s\n", stagingPort)
	if stagingPort != "" {
		err := KillProcessOnPort(stagingPort)
		if err != nil {
			fmt.Printf("Could not kill process on port %s\n", stagingPort)
			fmt.Println(err)
		}
	}

	// Remove config
	deleteNginxConfig(conf.Staging.UniqueName)

	// Restart nginx
	fmt.Println("Restarting nginx...")
	err := restartNginx()

	return err
}

func HandleRemoveProduction(conf StepConfig) error {
	fmt.Println("Killing existing process...")
	prodPort := getPortInNginxConfig(conf.Production.UniqueName)
	fmt.Printf("Existing staging port: %s\n", prodPort)
	if prodPort != "" {
		err := KillProcessOnPort(prodPort)
		if err != nil {
			fmt.Printf("Could not kill process on port %s\n", prodPort)
			fmt.Println(err)
		}
	}

	// Remove config
	deleteNginxConfig(conf.Production.UniqueName)

	// Restart nginx
	fmt.Println("Restarting nginx...")
	err := restartNginx()

	return err
}

func HandleDeployProduction(conf StepConfig) {

}

func HandleHealthCheckStaging(conf StepConfig) {

}

func HandleHealthCheckProduction(conf StepConfig) {

}

func runSetup(setupCommands []string, deployPort string) {

	// Create and open script file
	scriptFilePath := "./cyanic-scripts/tmp.sh"

	err := os.Chmod(scriptFilePath, 0755)
	if err != nil {
		fmt.Println("Error changing file permissions:", err)
		return
	}

	createFile(scriptFilePath)
	shFile, err := os.Create(scriptFilePath)

	if err != nil {
		fmt.Println("Error opening file")
		fmt.Println(err)
		return
	}

	shFile.WriteString("#!/bin/sh -ex\n")
	shFile.WriteString("export PORT=" + deployPort + "\n")

	// Create script content
	for _, line := range setupCommands {
		_, err := shFile.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writting line")
			fmt.Println(err)
		}
	}

	shFile.Close()

	// Run script file
	var stdout []byte
	var commandError error
	commandError = exec.Command(scriptFilePath, "&", "disown").Start()
	fmt.Println(string(stdout[:]))
	if commandError != nil {
		fmt.Println("Error running setup script:")
		fmt.Println(commandError)
	}
}
