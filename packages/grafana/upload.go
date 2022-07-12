package grafana

import (
	"context"
	"net/http"

	"github.com/K-Phoen/grabana"
	"github.com/K-Phoen/grabana/dashboard"
)

func UpsertDashboard(ctx context.Context, grafanaURL string, grafanaToken string, folder string, dashboard *dashboard.Builder) (string, error) {
	cli := grabana.NewClient(http.DefaultClient, grafanaURL, grabana.WithAPIToken(grafanaToken))
	gFolder, err := cli.FindOrCreateFolder(ctx, folder)
	if err != nil {
		return "", err
	}
	_, _ = cli.UpsertDashboard(ctx, gFolder, *dashboard)
	gDashboard, err := cli.GetDashboardByTitle(ctx, dashboard.Internal().Title)
	if err != nil {
		return "", err
	}
	return grafanaURL + gDashboard.URL, nil
}
