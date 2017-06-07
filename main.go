package main

import (
	"fmt"
	"log"

	"github.com/ddliu/go-httpclient"
	"os"
)

func main() {
	//var apiURL = "https://vpnx.pw/api/v1/nodes/list"
	var apiURL = "https://vpnx.pw/api/v1/set_node"
	unique := make(map[string]int)
	proc_ex := make(chan Process, 10)

	//go Tcp(proc_ex)
	close(proc_ex)

	for proc := range proc_ex {
		if proc.State == "ESTABLISHED" && proc.Port == 3000 {
			unique[proc.ForeignIp] = 0
		}
	}
	fmt.Println(len(unique))
	//for _, proc := range tcp_data {
	//	if proc.State == "ESTABLISHED" && proc.Port == 3000 {
	//		unique[proc.ForeignIp] = proc
	//	}
	//}
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT:      "vpnx_agent",
		"Accept-Language":             "en-us",
		httpclient.OPT_DEBUG:          true,
		httpclient.OPT_CONNECTTIMEOUT: 3,
	})
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	res, err := httpclient.Post(apiURL, map[string]string{
		"node": hostname,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Resonse status: %d\n", res.StatusCode)
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
