package url

import (
	"fmt"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/tnn-gruntwork-io/go-commons/errors"
)

// Create a URL with the given base, path parts, query string, and fragment. This method will properly URI encode
// everything and handle leading and trailing slashes.
func FormatUrl(baseUrl string, pathParts []string, query url.Values, fragment string) (string, error) {
	parsedUrl, err := url.Parse(stripSlashes(baseUrl))
	if err != nil {
		return "", errors.WithStackTrace(err)
	}

	normalizedPathParts := []string{}

	for _, pathPart := range pathParts {
		normalizedPathPart := stripSlashes(pathPart)
		normalizedPathParts = append(normalizedPathParts, normalizedPathPart)
	}

	if len(normalizedPathParts) > 0 {
		parsedUrl.Path = fmt.Sprintf("%s/%s", stripSlashes(parsedUrl.Path), strings.Join(normalizedPathParts, "/"))
	}

	parsedUrl.RawQuery = mergeQuery(parsedUrl.Query(), query).Encode()
	parsedUrl.Fragment = fragment

	return parsedUrl.String(), nil
}

// Merge the two query params together. The new query will override the original.
func mergeQuery(originalQuery url.Values, newQuery url.Values) url.Values {
	result := map[string][]string{}

	for key, values := range originalQuery {
		result[key] = values
	}

	for key, values := range newQuery {
		result[key] = values
	}

	return result
}

// Remove all leading or trailing slashes in the given string
func stripSlashes(str string) string {
	return strings.Trim(str, "/")
}

// Attempt to open a URL in the user's browser. We use this to open docs, PRs we've
// programmatically opened, etc
func OpenURL(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		return err
	}
	return nil
}
