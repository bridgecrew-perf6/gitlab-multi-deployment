package main

import (
	"errors"
	"os"
	"testing"
)

var (
	ErrTest = errors.New("some error")
)

const (
	testDir            = "test-pim/"
	testFile           = "test-pim/test-deployment.json"
	testAppInsightsKey = "TEST-APIK"
)

func setup(fileContent string) func() {
	os.Mkdir(testDir, os.ModePerm)
	os.WriteFile(testFile, []byte(fileContent), os.ModePerm)
	return func() {
		os.RemoveAll(testDir)
	}
}

func TestDeploymentSuccess(t *testing.T) {
	defer setup("[{\"repository\": \"testServiceA\",\"pipeline\": 1},{\"repository\": \"testServiceB\",\"pipeline\": 2},{\"repository\": \"component/testServiceC\",\"pipeline\": 3}]")()

	gitMock := &mockGitClient{cloneDir: testDir, shouldFail: false}
	gitlabMock := &mockGitlabClient{shouldFail: false}
	kubectlMock := &mockKubectl{deployed: []string{}, shouldFail: false}
	deployment := Deployment{
		repository:     gitMock,
		pipeline:       gitlabMock,
		kubernetes:     kubectlMock,
		appInsightsKey: testAppInsightsKey,
	}

	err := deployment.Run(testFile)
	if err != nil {
		t.Errorf("this deployment should succeed")
	}
	if len(kubectlMock.deployed) != 3 {
		t.Errorf("number of deployed services not correct")
	}

	file, _ := os.ReadFile(testDir + "testServiceB/application/kubernetes/prod/deployment.yaml")
	if string(file) != "2 + "+testAppInsightsKey {
		t.Errorf("variables should be injected to yaml files")
	}
	file, _ = os.ReadFile(testDir + "testServiceC/application/kubernetes/prod/deployment.yaml")
	if string(file) != "3 + "+testAppInsightsKey {
		t.Errorf("variables should be injected to yaml files")
	}
}

func TestDeploymentVerifyFailFileA(t *testing.T) {
	defer setup("[{\"repository\": \"testServiceA\",\"pipeline\": 1},{\"repository\": \"testServiceB\",\"pipeline\": },{\"repository\": \"component/testServiceC\",\"pipeline\": 3}]")()

	gitMock := &mockGitClient{cloneDir: testDir, shouldFail: false}
	gitlabMock := &mockGitlabClient{shouldFail: false}
	kubectlMock := &mockKubectl{deployed: []string{}, shouldFail: false}
	deployment := Deployment{
		repository:     gitMock,
		pipeline:       gitlabMock,
		kubernetes:     kubectlMock,
		appInsightsKey: testAppInsightsKey,
	}

	err := deployment.Run(testFile)
	if err == nil {
		t.Errorf("this deployment should fail")
	}
	if len(kubectlMock.deployed) != 0 {
		t.Errorf("number of deployed services not correct")
	}
}

func TestDeploymentVerifyFailFileB(t *testing.T) {
	defer setup("[{\"repository\": \"testServiceA\",\"pipeline\": 1},{\"repository\": \"testServiceB\" },{\"repository\": \"component/testServiceC\",\"pipeline\": 3}]")()

	gitMock := &mockGitClient{cloneDir: testDir, shouldFail: false}
	gitlabMock := &mockGitlabClient{shouldFail: false}
	kubectlMock := &mockKubectl{deployed: []string{}, shouldFail: false}
	deployment := Deployment{
		repository:     gitMock,
		pipeline:       gitlabMock,
		kubernetes:     kubectlMock,
		appInsightsKey: testAppInsightsKey,
	}

	err := deployment.Run(testFile)
	if err == nil {
		t.Errorf("this deployment should fail")
	}
	if len(kubectlMock.deployed) != 0 {
		t.Errorf("number of deployed services not correct")
	}
}

func TestDeploymentVerifyFailGitlab(t *testing.T) {
	defer setup("[{\"repository\": \"testServiceA\",\"pipeline\": 1},{\"repository\": \"testServiceB\",\"pipeline\": 2},{\"repository\": \"component/testServiceC\",\"pipeline\": 3}]")()

	gitMock := &mockGitClient{cloneDir: testDir, shouldFail: false}
	gitlabMock := &mockGitlabClient{shouldFail: true}
	kubectlMock := &mockKubectl{deployed: []string{}, shouldFail: false}
	deployment := Deployment{
		repository:     gitMock,
		pipeline:       gitlabMock,
		kubernetes:     kubectlMock,
		appInsightsKey: testAppInsightsKey,
	}

	err := deployment.Run(testFile)
	if err == nil {
		t.Errorf("this deployment should fail")
	}
	if len(kubectlMock.deployed) != 0 {
		t.Errorf("number of deployed services not correct")
	}
}

// Mocks
type mockKubectl struct {
	deployed   []string
	shouldFail bool
}

func (kubectl *mockKubectl) Apply(artifactName string, basePath string) error {
	kubectl.deployed = append(kubectl.deployed, artifactName)
	if kubectl.shouldFail {
		return ErrTest
	}

	return nil
}

type mockGitlabClient struct {
	shouldFail bool
}

func (api mockGitlabClient) GetCommitHash(repository string, pipelineID int) (string, error) {
	commitHash := repository + "-commit"
	if api.shouldFail {
		return "", ErrTest
	}
	return commitHash, nil
}

type mockGitClient struct {
	cloneDir   string
	shouldFail bool
}

func (gitClient *mockGitClient) CheckoutCommit(artifactName string, repository string, commit string) (string, error) {
	checkoutDir := gitClient.cloneDir + artifactName
	yamlFile := checkoutDir + "/application/kubernetes/prod/deployment.yaml"
	os.MkdirAll(checkoutDir+"/application/kubernetes/prod/", os.ModePerm)
	os.WriteFile(yamlFile, []byte("__IMAGE_VERSION__ + __APP_INSIGHTS_INTRUMENTATION_KEY__"), os.ModePerm)
	if gitClient.shouldFail {
		return "", ErrTest
	}

	return checkoutDir, nil
}

func (gitClient *mockGitClient) TagCommit(artifactName string, tag string) error {
	if gitClient.shouldFail {
		return ErrTest
	}

	return nil
}
