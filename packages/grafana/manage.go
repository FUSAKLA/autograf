package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/fusakla/sdk"
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
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%s %s", resp.Status, b)
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

type apiDashboard struct {
	FolderUid string      `json:"folderUid"`
	Overwrite bool        `json:"overwrite"`
	Dashboard interface{} `json:"dashboard"`
}

type apiDashbaordResp struct {
	Url string `json:"url"`
}

func (c *client) UpsertDashboard(folderUid string, dashboard *sdk.Board) (string, error) {
	dashboard.ID = 0
	data, err := json.Marshal(apiDashboard{FolderUid: folderUid, Overwrite: true, Dashboard: dashboard})
	if err != nil {
		return "", err
	}
	resp, err := c.cli.Post(c.url+"/api/dashboards/db", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		err, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("%s %s", resp.Status, err)
	}
	var dashboardResp apiDashbaordResp
	if err := json.NewDecoder(resp.Body).Decode(&dashboardResp); err != nil {
		return "", err
	}
	u, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}
	u.Path = dashboardResp.Url
	return u.String(), nil
}
