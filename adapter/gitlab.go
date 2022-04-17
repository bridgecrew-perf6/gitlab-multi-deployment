package adapter

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GitlabClient struct {
	accessToken string
	baseUrl string
}

type commit struct {
	Sha string `json:"sha"`
}

func NewGitlabClient(accessToken string, baseUrl string) *GitlabClient {
	client := GitlabClient{accessToken: accessToken, baseUrl: baseUrl}
	return &client
}

func (api GitlabClient) GetCommitHash(repository string, pipelineID int) (string, error) {
	apiUrl := api.baseUrl + strings.ReplaceAll(repository, "/", "%2F") + "/pipelines/" + strconv.Itoa(pipelineID)
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header = http.Header{
		"PRIVATE-TOKEN": []string{api.accessToken},
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode == 404 {
		return "", &url.Error{}
	}

	var pipelineCommit commit
	err = json.NewDecoder(res.Body).Decode(&pipelineCommit)
	if err != nil {
		return "", err
	}

	return pipelineCommit.Sha, nil
}
