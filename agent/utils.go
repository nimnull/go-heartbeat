package agent

import (
	"log"
	"net/url"
)

func ApiURIBuilder(apiFQDN string, useSSL bool) string {
	prefix := "http://"
	hostUrl, err := url.Parse(prefix + apiFQDN)

	if err != nil {
		log.Fatal(err)
	}

	if useSSL {
		hostUrl.Scheme = "https://"
	}

	if hostUrl.Host == "" {
		log.Fatalf("Can't recognize provided API host FQDN: %s\n", apiFQDN)
	}
	hostUrl.Path = "/api/v1/nodes/"
	return hostUrl.String()
}
