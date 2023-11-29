package handlers

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func createFile(filePath string) {

	// Extract the directory from the file path
	dir := filepath.Dir(filePath)

	// Create directories if they don't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		slog.Error("Error creating directories:", "msg", err.Error())
		return
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create it
		file, err := os.Create(filePath)
		if err != nil {
			slog.Error("Error creating file:", "msg", err.Error())
			return
		}
		defer file.Close()

		slog.Info("File created:", "path", filePath)
	} else {
		// File already exists
		slog.Debug("File already exists:", "path", filePath)
	}

}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		return false
	}
}

func deleteNginxConfig(configName string) {
	filePath := filepath.Join(nginxSitesPath, configName)
	os.Remove(filePath)
}

func extractProxyPort(nginxConfig string) (string, error) {
	// Find the index of "proxy_pass" in the nginx config
	proxyPassIndex := strings.Index(nginxConfig, "proxy_pass")

	if proxyPassIndex == -1 {
		return "", fmt.Errorf("proxy_pass directive not found in nginx configuration")
	}

	// Find the start index of the URL after "proxy_pass"
	urlStartIndex := strings.Index(nginxConfig[proxyPassIndex:], "http") + proxyPassIndex

	if urlStartIndex == -1 {
		return "", fmt.Errorf("http not found after proxy_pass in nginx configuration")
	}

	// Find the end index of the URL
	urlEndIndex := strings.Index(nginxConfig[urlStartIndex:], ";") + urlStartIndex

	if urlEndIndex == -1 {
		return "", fmt.Errorf("semicolon not found after proxy_pass URL in nginx configuration")
	}

	// Extract the URL
	proxyURL := nginxConfig[urlStartIndex:urlEndIndex]

	// Split the URL by colon to get the port
	urlParts := strings.Split(proxyURL, ":")
	if len(urlParts) < 3 {
		return "", fmt.Errorf("invalid proxy_pass URL format in nginx configuration")
	}

	// The port is the third part of the URL
	proxyPort := urlParts[2]

	return proxyPort, nil
}

var nginxSitesPath string = "/etc/nginx/sites-enabled"

func getPortInNginxConfig(fileName string) string {
	filePath := filepath.Join(nginxSitesPath, fileName)
	fileContent, err := os.ReadFile(filePath)

	if err != nil {
		return ""
	}

	port, err := extractProxyPort(string(fileContent))
	if err == nil {
		return port
	}

	return ""
}

func createNginxConfig(templatePath string, configName string, proxyPort string) error {

	filePath := filepath.Join(nginxSitesPath, configName)
	createFile(filePath)
	templateContent, err := os.ReadFile(templatePath)

	if err != nil {
		slog.Error("Could not read template file", "path", templatePath)
		slog.Error(err.Error())
		return err
	}

	configContent := strings.ReplaceAll(string(templateContent), "${PORT}", proxyPort)

	err = os.WriteFile(filePath, []byte(configContent), fs.FileMode(os.O_RDONLY))

	return err
}

func testNginx() error {
	err := exec.Command("sudo", "nginx", "-t").Run()
	return err
}

func restartNginx() error {
	//err := exec.Command("sudo", "systemctl", "restart", "nginx").Run()
	err := exec.Command("sudo", "nginx", "-s", "reload").Run()
	return err
}

func runSetup(setupCommands []string, deployPort string) {

	// Create and open script file
	scriptFilePath := "./cyanic-scripts/tmp.sh"

	createFile(scriptFilePath)

	err := os.Chmod(scriptFilePath, 0755)
	if err != nil {
		slog.Error("Error changing file permissions", "msg", err.Error())
		return
	}

	shFile, err := os.Create(scriptFilePath)

	if err != nil {
		slog.Error("Error opening file", "msg", err.Error())
		return
	}

	shFile.WriteString("#!/bin/sh -ex\n")
	shFile.WriteString("export PORT=" + deployPort + "\n")

	// Create script content
	for _, line := range setupCommands {
		_, err := shFile.WriteString(line + "\n")
		if err != nil {
			slog.Error("Error writting line", "msg", err.Error())
		}
	}

	shFile.Close()

	// Run script file
	var commandError error
	commandError = exec.Command(scriptFilePath, "&", "disown").Start()
	if commandError != nil {
		slog.Error("Error running setup script")
		slog.Error(commandError.Error())
	}
}
