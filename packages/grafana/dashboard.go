package grafana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fusakla/sdk"
)

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
	return c.url + dashboardResp.Url, nil
}
