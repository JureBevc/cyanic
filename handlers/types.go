package handlers

type StageConfig struct {
	UniqueName     string `yaml:"unique-name"`
	Nginx          string `yaml:"nginx"`
	HealthCheckUrl string `yaml:"health-check-url"`
}

type StepConfig struct {
	Clone      string      `yaml:"clone"`
	SSH        string      `yaml:"ssh"`
	Ports      []int       `yaml:"ports"`
	Staging    StageConfig `yaml:"staging"`
	Production StageConfig `yaml:"production"`
	Setup      []string    `yaml:"setup"`
}

type CyanicConfig struct {
	Step StepConfig `yaml:"cyanic"`
}
