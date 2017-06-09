// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package main

import (
	//"fmt"
	//"log"
	//"os"

	//"github.com/ddliu/go-httpclient"
	"node_agent/cmd"
)

func main() {
	cmd.Execute()
	////var apiURL = "https://vpnx.pw/api/v1/nodes/list"
	//var apiURL = "https://vpnx.pw/api/v1/set_node"
	//unique := make(map[string]int)
	//proc_ex := make(chan Process, 10)
	//
	////go Tcp(proc_ex)
	//close(proc_ex)
	//
	//for proc := range proc_ex {
	//	if proc.State == "ESTABLISHED" && proc.Port == 3000 {
	//		unique[proc.ForeignIp] = 0
	//	}
	//}
	//fmt.Println(len(unique))
	////for _, proc := range tcp_data {
	////	if proc.State == "ESTABLISHED" && proc.Port == 3000 {
	////		unique[proc.ForeignIp] = proc
	////	}
	////}
	//httpclient.Defaults(httpclient.Map{
	//	httpclient.OPT_USERAGENT:      "vpnx_agent",
	//	"Accept-Language":             "en-us",
	//	httpclient.OPT_DEBUG:          true,
	//	httpclient.OPT_CONNECTTIMEOUT: 3,
	//})
	//hostname, err := os.Hostname()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//res, err := httpclient.Post(apiURL, map[string]string{
	//	"node": hostname,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("Resonse status: %d\n", res.StatusCode)
}

//@app.route('/')
//def main_route():
//return 'Forbidden to route.', 403
//
//@app.route('/api')
//def main_api_route():
//return 'Forbidden to route.', 403
//
//@app.route('/api/v1/')
//def get_api_v1():
//return 'Forbidden to route.', 403
//
//@app.route('/api/v1/proxy/users',  methods=['GET'] )
