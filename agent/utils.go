// Copyright © 2017 Yehor Nazarkin <nimnull@gmail.com>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		hostUrl.Scheme = "https"
	}

	if hostUrl.Host == "" {
		log.Fatalf("Can't recognize provided API host FQDN: %s\n", apiFQDN)
	}
	hostUrl.Port()
	hostUrl.Path = "/api/v1/nodes/"
	return hostUrl.String()
}
