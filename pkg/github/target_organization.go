package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v45/github"
)

type TargetOrganization struct {
	client *github.Client

	Name string
}

func NewTargetOrganization(http *http.Client, name string) *TargetOrganization {
	client := github.NewClient(http)
	return &TargetOrganization{client: client, Name: name}
}

func (t *TargetOrganization) URL() string {
	return fmt.Sprintf("https://github.com/%s", t.Name)
}

func (t *TargetOrganization) GetRunners(
	ctx context.Context, page int, pageSize int,
) ([]*github.Runner, int, error) {
	runners, resp, err := t.client.Actions.ListOrganizationRunners(
		ctx, t.Name,
		&github.ListOptions{Page: page, PerPage: pageSize},
	)
	if err != nil {
		return nil, 0, err
	}

	return runners.Runners, resp.NextPage, nil
}
