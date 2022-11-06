package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type apiFoldersList []apiFolder

type apiFolder struct {
	Uid   string `json:"uid,omitempty"`
	Title string `json:"title"`
}

func (c *client) createFolder(name string) (string, error) {
	data, err := json.Marshal(apiFolder{Title: name})
	if err != nil {
		return "", err
	}
	resp, err := c.cli.Post(c.url+"/api/folders", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s", resp.Body)
	}
	var folder apiFolder
	if err := json.NewDecoder(resp.Body).Decode(&folder); err != nil {
		return "", err
	}
	return folder.Uid, nil
}

func (c *client) EnsureFolder(name string) (string, error) {
	resp, err := c.cli.Get(c.url + "/api/folders")
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%s", resp.Body)
	}
	var folders apiFoldersList
	if err := json.NewDecoder(resp.Body).Decode(&folders); err != nil {
		return "", err
	}
	if len(folders) > 0 {
		return folders[0].Uid, nil
	}
	return c.createFolder(name)
}
