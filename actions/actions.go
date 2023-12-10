package actions

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/JureBevc/cyanic/handlers"
	"github.com/JureBevc/cyanic/util"
	"github.com/goccy/go-yaml"
)

func PrintHelp() {
	fmt.Println("Available commands:")
	fmt.Println("\t help")
	fmt.Println("\t port-status [path-to-config]")
	fmt.Println("\t deploy-staging [path-to-config]")
	fmt.Println("\t swap [path-to-config]")
	fmt.Println("\t remove-staging [path-to-config]")
	fmt.Println("\t remove-production [path-to-config]")
	fmt.Println("\t health-staging [path-to-config]")
	fmt.Println("\t health-production [path-to-config]")
}

func handleAction(stepConfig handlers.StepConfig, action string) error {
	slog.Info("Running cyanic", "action", action)
	switch action {
	case "full-deploy":
		return handlers.HandleFullDeploy(stepConfig)

	case "deploy-staging":
		return handlers.HandleDeployStaging(stepConfig)

	case "deploy-production":
		return handlers.HandleDeployProduction(stepConfig)

	case "remove-staging":
		return handlers.HandleRemoveStaging(stepConfig)

	case "remove-production":
		return handlers.HandleRemoveProduction(stepConfig)

	case "swap":
		return handlers.HandleSwap(stepConfig)

	case "health-staging":
		health := handlers.HandleHealthCheckStaging(stepConfig)
		if !health {
			return errors.New("Health check failed")
		}

	case "health-production":
		health := handlers.HandleHealthCheckProduction(stepConfig)
		if !health {
			return errors.New("Health check failed")
		}

	case "port-status":
		handlers.HandlePortStatus(stepConfig)

	default:
		slog.Error("Invalid action", "name", action)
		return errors.New("Invalid action: " + action)
	}

	return nil
}

func handleNonConfigActions(action string, params []string) (bool, error) {

	switch action {

	case "kill-port":
		requireParamCount(params, 1)
		err := handlers.KillProcessOnPort(params[0])
		return true, err

	case "help":
		PrintHelp()
		return true, nil

	}

	return false, nil
}

func requireParamCount(params []string, count int) {
	if len(params) < count {
		slog.Error("Not enough parameters", "required", count, "got", len(params))
		os.Exit(1)
	}
}

var defaultConfigPath string = "./cyanic.yaml"

func ParseCommand(action string, params []string) error {

	if didRun, err := handleNonConfigActions(action, params); didRun {
		return err
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
		return err
	}

	cyConfig := handlers.CyanicConfig{}

	if err = yaml.Unmarshal(content, &cyConfig); err != nil {
		slog.Error("Error parsing configuration file:")
		slog.Error(err.Error())
		return err
	}

	slog.Debug(util.StructToString(cyConfig))

	err = handleAction(cyConfig.Step, action)
	return err
}
