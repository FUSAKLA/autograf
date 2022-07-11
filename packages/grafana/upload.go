package grafana

import (
	"context"
	"net/http"
	"net/url"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
)

func UpsertDashboard(ctx context.Context, grafanaURL *url.URL, grafanaToken string, folder string, dashboard *dashboard.Builder) (string, error) {
	cli := grabana.NewClient(http.DefaultClient, grafanaURL.String(), grabana.WithAPIToken(grafanaToken))
	gFolder, err := cli.FindOrCreateFolder(ctx, folder)
	if err != nil {
		return "", err
	}
	cli.UpsertDashboard(ctx, gFolder, *dashboard)
	gDashboard, err := cli.GetDashboardByTitle(ctx, dashboard.Internal().Title)
	if err != nil {
		return "", err
	}
	return grafanaURL.String() + gDashboard.URL, nil
}
