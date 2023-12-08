# Cyanic 
![GitHub tag](https://img.shields.io/github/v/tag/JureBevc/cyanic)
![GitHub release](https://img.shields.io/github/v/release/JureBevc/cyanic)
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/JureBevc/cyanic/unit_tests.yml?label=tests)
![Static Badge](https://img.shields.io/badge/license-MIT-blue)



Blue-Green Deployment Tool

## Overview

Cyanic is a lightweight and easy-to-use tool written in Go for automating [blue-green](https://en.wikipedia.org/wiki/Blue%e2%80%93green_deployment) deployments. With Cyanic, you can seamlessly deploy your application to staging and production environments, perform swaps between the environments without downtime, and remove deployments when needed. It is designed with simplicity in mind, making it particularly suitable for small to medium-sized apps or your personal proof of concept deployments.

This tool proves particularly advantageous for users who do not utilize Docker and seek a straightforward deployment solution behind an Nginx proxy.

## Installation

If you are not using a precompiled binary, you can use Cyanic as a Go package. Ensure you have Go installed on your system, then install Cyanic using the following command:

```bash
go get -u github.com/JureBevc/cyanic
```

## Requirements

Cyanic relies heavily on nginx and requires some commands to be available on the deployment server:

- `nginx`
- `systemctl`
- `disown`
- `fuser`

## Usage

Show a list of available commands:
```bash
cyanic help
```

Deploy your application to the staging environment:
```bash
cyanic deploy-staging
```

Swap staging and production environments:
```bash
cyanic swap
```

Run health check for staging:
```bash
cyanic health-staging
```

Run health check for production:
```bash
cyanic health-production
```

Run a fully automatic deploy - deploys to staging, performs health check and swaps with production and health checks again:
```bash
cyanic full-deploy
```

Remove the staging deployment when it is no longer needed:
```bash
cyanic remove-staging
```

Remove the production deployment when it is no longer needed:
```bash
cyanic remove-production
```

Show the status of ports listed in the configuration:
```bash
cyanic port-status
```

## Configuration

Cyanic uses a configuration file (cyanic.yaml by default) for specifying deployment settings. Make sure to configure the file according to your application's requirements.

This repository also contains examples, so you can take a look at the configuration in `cyanic.yaml` and example projects inside the `examples` directory.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or create a pull request.

## License

This tool is open-source and available under the MIT License. 
