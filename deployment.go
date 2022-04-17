package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type kubernetesPort interface {
	Apply(artifactName string, basePath string) error
}

type repositoryPort interface {
	CheckoutCommit(artifactName string, repository string, commit string) (string, error)
	TagCommit(artifactName string, tag string) error
}

type pipelinePort interface {
	GetCommitHash(repository string, pipelineID int) (string, error)
}

type artifact struct {
	Name       string
	Repository string `json:"repository"`
	PipelineID int    `json:"pipeline"`
}

type Deployment struct {
	repository     repositoryPort
	pipeline       pipelinePort
	kubernetes     kubernetesPort
	appInsightsKey string
	cloneDir       string
}

func (deployment Deployment) Run(deploymentFilePath string) error {
	artifacts, err := deployment.parseDeploymentFile(deploymentFilePath)
	if err != nil {
		log.Printf("Invalid deployment.json")
		return err
	}

	log.Printf("Veryfing deployment config...")
	for _, artifact := range artifacts {
		log.Printf("Veryfing %+v", artifact)
		commit, err := deployment.pipeline.GetCommitHash(artifact.Repository, artifact.PipelineID)
		if err != nil {
			log.Printf("Pipeline information not found.")
			return err
		}

		repoPath, err := deployment.repository.CheckoutCommit(artifact.Name, artifact.Repository, commit)
		if err != nil {
			log.Printf("Checking out repository failed.")
			return err
		}

		err = deployment.injectVariables(repoPath, artifact)
		if err != nil {
			log.Printf("Kubernetes config files not found")
			return err
		}
	}

	log.Printf("Verification successfull. Starting deployments...")
	for _, artifact := range artifacts {
		log.Printf("Deploying %+v", artifact)

		err := deployment.kubernetes.Apply(artifact.Name, deployment.cloneDir)
		if err != nil {
			log.Printf("Kubectl deployment failed.")
			return err
		}

		err = deployment.repository.TagCommit(artifact.Name, "PROD-"+time.Now().Format("02-01-2006-15-04-05"))
		if err != nil {
			log.Printf("WARN: Tagging failed! Will continue...")
		}
	}

	return nil
}

func (deployment Deployment) injectVariables(repoPath string, artifact artifact) error {
	const path = "/application/kubernetes/prod/deployment.yaml"
	deploymentFile, err := os.ReadFile(repoPath + path)
	if err != nil {
		return err
	}

	updatedFile := strings.ReplaceAll(string(deploymentFile), "__IMAGE_VERSION__", strconv.Itoa(artifact.PipelineID))
	updatedFile = strings.ReplaceAll(updatedFile, "__APP_INSIGHTS_INTRUMENTATION_KEY__", deployment.appInsightsKey)
	err = os.WriteFile(repoPath+path, []byte(updatedFile), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (deployment Deployment) parseDeploymentFile(path string) ([]artifact, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	artifacts := make([]artifact, 0)
	err = json.Unmarshal(input, &artifacts)
	if err != nil {
		return nil, err
	}
	for i := range artifacts {
		artifact := &artifacts[i]
		if artifact.Repository != "" && artifact.PipelineID > 0 {
			pathParts := strings.Split(artifact.Repository, "/")
			artifact.Name = pathParts[len(pathParts)-1]
		} else {
			return nil, errors.New("deployment file not valid")
		}
	}

	return artifacts, nil
}
