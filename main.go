package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/JureBevc/cyanic/handlers"
	"github.com/JureBevc/cyanic/util"
	"github.com/goccy/go-yaml"
)

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("\t help")
	fmt.Println("\t deploy-staging [path-to-config]")
	fmt.Println("\t swap [path-to-config]")
	fmt.Println("\t remove-staging [path-to-config]")
	fmt.Println("\t remove-production [path-to-config]")
}

func handleAction(stepConfig handlers.StepConfig, action string) {
	slog.Info("Running cyanic", "action", action)
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
		slog.Error("Invalid action", "name", action)
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
		slog.Error("Not enough parameters", "required", count, "got", len(params))
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
		slog.Error("No valid configuration path found")
		os.Exit(1)
	}

	content, err := os.ReadFile(configFilePath)

	if err != nil {
		slog.Error("Could not read configuration file:")
		slog.Error(err.Error())
		return
	}

	cyConfig := handlers.CyanicConfig{}

	if err = yaml.Unmarshal(content, &cyConfig); err != nil {
		slog.Error("Error parsing configuration file:")
		slog.Error(err.Error())
		return
	}

	slog.Debug(util.StructToString(cyConfig))

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
