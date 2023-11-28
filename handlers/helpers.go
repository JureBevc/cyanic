package handlers

import (
	"fmt"
	"io/fs"
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
		fmt.Println("Error creating directories:", err)
		return
	}

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// File does not exist, create it
		file, err := os.Create(filePath)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		fmt.Println("File created:", filePath)
	} else {
		// File already exists
		// fmt.Println("File already exists:", filePath)
	}

}

func deleteNginxConfig(configName string){
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
		fmt.Printf("Could not read template file %s\n", templatePath)
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
	err := exec.Command("sudo", "systemctl", "restart", "nginx").Run()
	return err
}

func killProcessOnPort(portString string) error {
	err := exec.Command("sudo", "fuser", "-k", "-n", "tcp", portString).Run()
	return err
}
