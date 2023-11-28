package main

import (
	"JureBevc/cyanic/handlers"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("\tcyanic help")
	fmt.Println("\tcyanic run [file]")
}

func handleAction(stepConfig handlers.StepConfig, action string) {
	fmt.Printf("-- RUNNING ACTION %s --\n", action)
	switch action {
	case "full-deploy":
		handlers.HandleFullDeploy(stepConfig)

	case "deploy-staging":
		handlers.HandleDeployStaging(stepConfig)

	case "deploy-production":
		handlers.HandleDeployProduction(stepConfig)

	case "remove-staging":
		handlers.HandleRemoveStaging(stepConfig)

	case "remove-production":
		handlers.HandleRemoveProduction(stepConfig)

	case "swap":
		handlers.HandleSwap(stepConfig)

	case "check-staging":
		handlers.HandleHealthCheckStaging(stepConfig)

	case "check-production":
		handlers.HandleHealthCheckProduction(stepConfig)

	default:
		fmt.Printf("Invalid action '%s'\n", action)
	}
}

func handleNonConfigActions(action string, params []string) bool {

	switch action {
	case "kill-port":
		requireParamCount(params, 1)
		handlers.KillProcessOnPort(params[0])
		return true
	}

	return false
}

func requireParamCount(params []string, count int) {
	if len(params) < count {
		fmt.Printf("Not enough parameters (required %d got %d)\n", count, len(params))
		os.Exit(1)
	}
}

var defaultConfigPath string = "./cyanic.yaml"

func runConfig(action string, params []string) {

	if handleNonConfigActions(action, params) {
		return
	}

	configFilePath := ""
	if handlers.FileExists(defaultConfigPath) {
		configFilePath = defaultConfigPath
	} else {
		requireParamCount(params, 1)
		configFilePath = params[0]
	}

	if configFilePath == "" {
		fmt.Println("No valid configuration path found")
		os.Exit(1)
	}

	content, err := os.ReadFile(configFilePath)

	if err != nil {
		fmt.Printf("Could not read configuration file:")
		fmt.Printf("%s\n", err)
		return
	}

	cyConfig := handlers.CyanicConfig{}

	if err = yaml.Unmarshal(content, &cyConfig); err != nil {
		fmt.Println("Error parsing configuration file:")
		fmt.Printf("%s\n", err)
		return
	}

	//util.PrintStruct(cyConfig)

	handleAction(cyConfig.Step, action)
}

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	switch command {
	case "help":
		printHelp()
	default:
		runConfig(command, os.Args[2:])
	}

}
