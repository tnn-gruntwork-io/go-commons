package git

import (
	"fmt"

	"github.com/gruntwork-io/go-commons/shell"
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
)

// ConfigureForceHTTPS configures git to force usage of https endpoints instead of SSH based endpoints for the three
// primary VCS platforms (GitHub, GitLab, BitBucket).
func ConfigureForceHTTPS(logger *logrus.Logger) error {
	opts := shell.NewShellOptions()
	if logger != nil {
		opts.Logger = logger
	}

	var allErr error

	for _, host := range []string{"github.com", "gitlab.com", "bitbucket.org"} {
		if err := shell.RunShellCommand(
			opts,
			"git", "config", "--global",
			fmt.Sprintf("url.https://%s.insteadOf", host),
			fmt.Sprintf("ssh://git@%s", host),
		); err != nil {
			allErr = multierror.Append(allErr, err)
		}

		if err := shell.RunShellCommand(
			opts,
			"git", "config", "--global",
			fmt.Sprintf("url.https://%s/.insteadOf", host),
			fmt.Sprintf("git@%s:", host),
		); err != nil {
			allErr = multierror.Append(allErr, err)
		}
	}
	return allErr
}

// ConfigureHTTPSAuth configures git with username and password to authenticate with the given VCS host when interacting
// with git over HTTPS. This uses the cache credentials store to configure the credentials. Refer to the git
// documentation on credentials storage for more information:
// https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage
func ConfigureHTTPSAuth(logger *logrus.Logger, gitUsername string, gitOauthToken string, vcsHost string) error {
	opts := shell.NewShellOptions()
	if logger != nil {
		opts.Logger = logger
	}

	if err := shell.RunShellCommand(
		opts,
		"git", "config", "--global",
		"credential.helper",
		"cache --timeout 3600",
	); err != nil {
		return err
	}

	if gitUsername == "" {
		gitUsername = "git"
	}
	credentialsStoreInput := fmt.Sprintf(`protocol=https
host=%s
username=%s
password=%s`, vcsHost, gitUsername, gitOauthToken)
	return shell.RunShellCommandWithInput(opts, credentialsStoreInput, "git", "credential-cache", "store")
}