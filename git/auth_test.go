package git

import (
	"strings"
	"testing"

	"github.com/tnn-gruntwork-io/terratest/modules/docker"
	terragit "github.com/tnn-gruntwork-io/terratest/modules/git"
	"github.com/tnn-gruntwork-io/terratest/modules/random"
)

func TestIntegrationGitAuth(t *testing.T) {
	t.Parallel()

	tag := "tnn-gruntwork-io/go-commons:" + strings.ToLower(random.UniqueId())
	ref := terragit.GetCurrentGitRef(t)
	docker.Build(t, "./test", &docker.BuildOptions{
		Tags:      []string{tag},
		BuildArgs: []string{"repo_ref=" + ref},
	})
	docker.Run(t, tag, &docker.RunOptions{
		Remove:               true,
		EnvironmentVariables: []string{"GITHUB_OAUTH_TOKEN"},
	})
}
