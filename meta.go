package lighter

import (
	"net/http"
	"os"
	"strings"

	m "cloud.google.com/go/compute/metadata"
)


var (
	agentName   = "lighter"

	projectKeys = []string{
		"GCP_PROJECT",
		"PROJECT",
		"PROJECT_ID",
		"GOOGLE_CLOUD_PROJECT",
		"GCLOUD_PROJECT",
		"CLOUDSDK_CORE_PROJECT",
	}
)

func getClient() *m.Client {
	return m.NewClient(&http.Client{
		Transport: userAgentTransport{
			userAgent: agentName,
			base:      http.DefaultTransport,
		},
	})
}

func getProjectID() (id string, err error) {
	for _, key := range projectKeys {
		if val, ok := os.LookupEnv(key); ok {
			return strings.TrimSpace(val), nil
		}
	}
	return getClient().ProjectID()
}

type userAgentTransport struct {
	userAgent string
	base      http.RoundTripper
}

// RoundTrip implements the transport interface
// // https://godoc.org/cloud.google.com/go/compute/metadata#example-NewClient
func (t userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}
