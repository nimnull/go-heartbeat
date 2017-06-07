package main

import "fmt"

func main() {
	unique := make(map[string]int)
	proc_ex := make(chan Process, 10)

	go Tcp(proc_ex)

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

}
