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
	"net/http"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
)

func StartReactor() {
	sleepSecs := viper.GetInt("rtime")
	apiURL := viper.GetString("api")
	nodeName := viper.GetString("nodename")
	watchPort := viper.GetInt("port")
	debugFlag := viper.GetBool("debug")

	log.Printf("Debug:\t%s\n", strconv.FormatBool(debugFlag))
	log.Printf("API:\t%s\n", apiURL)
	log.Printf("Node name:\t%s\n", nodeName)
	log.Printf("Port to watch:\t%s\n", strconv.FormatInt(int64(watchPort), 10))

	ticker := time.NewTicker(time.Duration(sleepSecs) * time.Second)
	quit := make(chan struct{})

	for {
		select {
		case <-ticker.C:
			protect(stateUpdateExecutor)
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

func protect(f func()) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("Runtime panic %v\n", x)
		}
	}()
	f()
}

func stateUpdateExecutor() {
	apiURL := viper.GetString("api")
	nodeName := viper.GetString("nodename")
	watchPort := viper.GetInt("port")
	debugFlag := viper.GetBool("debug")

	unique := make(map[string]int)
	// channel would be closed as soon Tcp() will finish collecting info
	proc_ex := make(chan Process, 10)
	go Tcp(proc_ex)

	for proc := range proc_ex {
		if proc.State == ESTABLISHED && proc.Port == int64(watchPort) {
			if _, ok := unique[proc.ForeignIp]; ok {
				unique[proc.ForeignIp] += 1
			} else {
				unique[proc.ForeignIp] = 0
			}
		}
	}

	request := gorequest.New()
	resp, body, errs := request.
		SetDebug(debugFlag).
		Timeout(time.Second*3).
		Set("Accept", "application/json").
		Set("Accept-Language", "en-us").
		Set("User-Agent", "node_agent_v1.0").
		Post(apiURL).
		Type("form").
		Send(map[string]string{
			"node":        nodeName,
			"connections": strconv.FormatInt(int64(len(unique)), 10),
			"port":        strconv.FormatInt(int64(watchPort), 10),
		}).
		Retry(3, 5*time.Second, http.StatusBadGateway, http.StatusGatewayTimeout).
		End()
	if len(errs) > 0 {
		log.Printf("Request to %s failed: %s\n", apiURL, resp.Status)
		for _, err := range errs {
			log.Printf("%s\n", err)
		}
		log.Println(body)
	}
}
