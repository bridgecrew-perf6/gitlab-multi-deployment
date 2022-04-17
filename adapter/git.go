package adapter

import (
	"errors"
	"os/exec"
)

type GitClient struct {
	cloneDir     string
	repoUrl      string
	repoUser     string
	repoPw       string
	repositories map[string]string
}

func NewGitClient(cloneDir string, repoUrl string, repoUser string, repoPw string) *GitClient {
	client := GitClient{
		cloneDir:     cloneDir,
		repoUrl:      repoUrl,
		repoUser:     repoUser,
		repoPw:       repoPw,
		repositories: make(map[string]string),
	}
	return &client
}

func (gitClient GitClient) CheckoutCommit(serviceName string, repoName string, commit string) (string, error) {
	checkoutDir := gitClient.cloneDir + serviceName
	url := "https://" + gitClient.repoUser + ":" + gitClient.repoPw + gitClient.repoUrl + repoName + ".git"
	_, err := exec.Command("git", "clone", url, checkoutDir).Output()
	if err != nil {
		return checkoutDir, err
	}
	gitClient.repositories[serviceName] = checkoutDir

	_, err = exec.Command("git", "-C", checkoutDir, "checkout", "--detach", commit).Output()
	if err != nil {
		return checkoutDir, err
	}

	return checkoutDir, nil
}

func (gitClient GitClient) TagCommit(serviceName string, tag string) error {
	repo := gitClient.repositories[serviceName]
	if repo == "" {
		return errors.New("service repository not found")
	}

	_, err := exec.Command("git", "-C", repo, "tag", "-f", tag).Output()
	if err != nil {
		return err
	}

	_, err = exec.Command("git", "-C", repo, "push", "-o", "ci.skip", "origin", tag).Output()
	if err != nil {
		return err
	}

	return nil
}
