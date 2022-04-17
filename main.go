package main

import (
	"gitlab-multi-deployment/adapter"
	"log"
	"os"
)

const (
	AccessTokenKey  = "ACCESS_TOKEN"
	AppInsightsKey  = "APP_INSIGHTS_KEY"
	GitRepoUrlKey   = "GIT_REPO_KEY"
	RegistryUserKey = "REGISTRY_USER"
	RegistryPwKey   = "REGISTRY_PW"
	GitlabUrlKey    = "GITLAB_API_URL"
	CloneDir        = "/tmp/services/"
)

func main() {
	deploymentFile := os.Args[1]

	deployment := Deployment{
		repository:     adapter.NewGitClient(CloneDir, os.Getenv(GitRepoUrlKey), os.Getenv(RegistryUserKey), os.Getenv(RegistryPwKey)),
		pipeline:       adapter.NewGitlabClient(os.Getenv(AccessTokenKey), os.Getenv(GitlabUrlKey)),
		kubernetes:     adapter.NewKubectlClient(os.Getenv(RegistryUserKey), os.Getenv(RegistryPwKey)),
		cloneDir:       CloneDir,
	}

	err := deployment.Run(deploymentFile)
	if err != nil {
		panic("Deployment finished with failures. Check logs above!")
	}

	log.Print("All done.")
}
