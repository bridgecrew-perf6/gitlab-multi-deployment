package adapter

import (
	"log"
	"os/exec"
)

type KubectlClient struct {
	registryUser string
	registryPw   string
}

func NewKubectlClient(registryUser string, registryPw string) *KubectlClient {
	client := KubectlClient{
		registryUser: registryUser,
		registryPw:   registryPw,
	}
	return &client
}

func (kubectl KubectlClient) Apply(serviceName string, basePath string) error {
	out, err := exec.Command("./adapter/deploy.sh", serviceName, kubectl.registryUser, kubectl.registryPw,
		basePath+serviceName).CombinedOutput()
	log.Print(string(out))
	if err != nil {
		return err
	}

	return nil
}
