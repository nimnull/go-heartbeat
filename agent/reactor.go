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

	for {
		go stateUpdateExecutor(
			viper.GetString("api"),
			viper.GetString("nodename"),
			viper.GetInt("port"),
			viper.GetBool("debug"))
		time.Sleep(time.Second * time.Duration(sleepSecs))
	}

}

func stateUpdateExecutor(apiHost, nodeName string, nodePort int, debug bool) {

	unique := make(map[string]int)
	proc_ex := make(chan Process, 10)
	defer close(proc_ex)

	go Tcp(proc_ex)

	for proc := range proc_ex {
		if proc.State == ESTABLISHED && proc.Port == int64(nodePort) {
			if _, ok := unique[proc.ForeignIp]; ok {
				unique[proc.ForeignIp] += 1
			} else {
				unique[proc.ForeignIp] = 0
			}
		}
	}
	request := gorequest.New()
	resp, body, errs := request.
		SetDebug(debug).
		Timeout(time.Second*3).
		Set("Accept", "application/json").
		Set("Accept-Language", "en-us").
		Post(apiHost).
		Retry(3, 5*time.Second, http.StatusBadGateway, http.StatusGatewayTimeout).
		SendMap(map[string]string{
			"node":        nodeName,
			"connections": strconv.FormatInt(int64(len(unique)), 10),
			"port":        strconv.FormatInt(int64(nodePort), 10),
		}).
		End()
	if len(errs) > 0 {
		log.Printf("Request to %s failed: %s\n", apiHost, resp.Status)
		for _, err := range errs {
			log.Printf("%s\n", err)
		}
		log.Println(body)
	}
}
