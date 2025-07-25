package grafana

import (
	"fmt"
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/gosimple/slug"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	grafanaClient "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/folders"
	"github.com/grafana/grafana-openapi-client-go/models"
)

func NewClient(urlStr string, token string) *client {
	url, _ := url.Parse(urlStr)
	fmt.Println("Connecting to Grafana at", url.String())
	return &client{
		grafanaCli: grafanaClient.NewHTTPClientWithConfig(strfmt.Default, &grafanaClient.TransportConfig{
			Host:     url.Host,
			BasePath: url.Path + "/api",
			Schemes:  []string{url.Scheme},
			APIKey:   token,
			// Debug:    true,
		}),
		url: url.String(),
	}
}

type client struct {
	url        string
	grafanaCli *grafanaClient.GrafanaHTTPAPI
}

func (c *client) DatasourceIDByName(name string) (string, error) {
	resp, err := c.grafanaCli.Datasources.GetDataSourceByName(name)
	if err != nil {
		return "", fmt.Errorf("error getting data sources: %w", err)
	}
	return resp.Payload.UID, nil
}

func (c *client) EnsureFolder(name string) (string, error) {
	resp, err := c.grafanaCli.Folders.GetFolders(&folders.GetFoldersParams{})
	if err != nil {
		return "", fmt.Errorf("error getting folders: %w", err)
	}
	for _, folder := range resp.Payload {
		if folder.Title == name {
			return folder.UID, nil
		}
	}
	newFolderResp, err := c.grafanaCli.Folders.CreateFolder(&models.CreateFolderCommand{
		Title: name,
	})
	if err != nil {
		return "", fmt.Errorf("error creating folder: %w", err)
	}
	return newFolderResp.Payload.UID, nil
}

func (c *client) UpsertDashboard(folderUid string, dashboard dashboard.Dashboard) (string, error) {
	dashboard.Uid = cog.ToPtr(slug.Make(*dashboard.Title))
	resp, err := c.grafanaCli.Dashboards.PostDashboard(&models.SaveDashboardCommand{
		FolderUID: folderUid,
		Overwrite: true,
		Dashboard: dashboard,
	})
	if err != nil {
		return "", fmt.Errorf("error saving dashboard: %w", err)
	}
	if resp == nil {
		return "", fmt.Errorf("error saving dashboard: response is nil")
	}
	u, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}
	u.Path = *resp.Payload.URL
	return u.String(), nil
}
