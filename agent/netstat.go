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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	PROC_TCP  = "/proc/net/tcp"
	PROC_UDP  = "/proc/net/udp"
	PROC_TCP6 = "/proc/net/tcp6"
	PROC_UDP6 = "/proc/net/udp6"

	ESTABLISHED = "ESTABLISHED"
	SYN_SENT    = "SYN_SENT"
	SYN_RECV    = "SYN_RECV"
	FIN_WAIT1   = "FIN_WAIT1"
	FIN_WAIT2   = "FIN_WAIT2"
	TIME_WAIT   = "TIME_WAIT"
	CLOSE       = "CLOSE"
	CLOSE_WAIT  = "CLOSE_WAIT"
	LAST_ACK    = "LAST_ACK"
	LISTEN      = "LISTEN"
	CLOSING     = "CLOSING"
)

var STATE = map[string]string{
	"01": ESTABLISHED,
	"02": SYN_SENT,
	"03": SYN_RECV,
	"04": FIN_WAIT1,
	"05": FIN_WAIT2,
	"06": TIME_WAIT,
	"07": CLOSE,
	"08": CLOSE_WAIT,
	"09": LAST_ACK,
	"0A": LISTEN,
	"0B": CLOSING,
}

type Process struct {
	User        int64
	Name        string
	Pid         string
	Exe         string
	State       string
	Ip          string
	Port        int64
	ForeignIp   string
	ForeignPort int64
}

func getData(t string) []string {
	// Get data from tcp or udp file.

	var proc_t string

	switch t {
	case "tcp":
		proc_t = PROC_TCP
	case "udp":
		proc_t = PROC_UDP
	case "tcp6":
		proc_t = PROC_TCP6
	case "udp6":
		proc_t = PROC_UDP6
	default:
		log.Fatalf("%s is a invalid type, tcp and udp only!\n", t)
	}

	data, err := ioutil.ReadFile(proc_t)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(data), "\n")

	// Return lines without Header line and blank line on the end
	return lines[1 : len(lines)-1]

}

func hexToDec(h string) int64 {
	// convert hexadecimal to decimal.
	d, err := strconv.ParseInt(h, 16, 32)
	if err != nil {
		log.Fatal(err)
	}

	return d
}

func convertIp(ip string) string {
	// Convert the ipv4 to decimal. Have to rearrange the ip because the
	// default value is in little Endian order.

	var out string

	// Check ip size if greater than 8 is a ipv6 type
	if len(ip) > 8 {
		i := []string{ip[30:32],
			ip[28:30],
			ip[26:28],
			ip[24:26],
			ip[22:24],
			ip[20:22],
			ip[18:20],
			ip[16:18],
			ip[14:16],
			ip[12:14],
			ip[10:12],
			ip[8:10],
			ip[6:8],
			ip[4:6],
			ip[2:4],
			ip[0:2]}
		out = fmt.Sprintf("%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v:%v%v",
			i[14], i[15], i[13], i[12],
			i[10], i[11], i[8], i[9],
			i[6], i[7], i[4], i[5],
			i[2], i[3], i[0], i[1])

	} else {
		i := []int64{hexToDec(ip[6:8]),
			hexToDec(ip[4:6]),
			hexToDec(ip[2:4]),
			hexToDec(ip[0:2])}

		out = fmt.Sprintf("%v.%v.%v.%v", i[0], i[1], i[2], i[3])
	}
	return out
}

func findPid(inode string) string {
	// Loop through all fd dirs of process on /proc to compare the inode and
	// get the pid.

	pid := "-"

	d, err := filepath.Glob("/proc/[0-9]*/fd/[0-9]*")
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(inode)
	for _, item := range d {
		path, _ := os.Readlink(item)
		out := re.FindString(path)
		if len(out) != 0 {
			pid = strings.Split(item, "/")[2]
		}
	}
	return pid
}

func getProcessExe(pid string) string {
	exe := fmt.Sprintf("/proc/%s/exe", pid)
	path, _ := os.Readlink(exe)
	return path
}

func getProcessName(exe string) string {
	n := strings.Split(exe, "/")
	name := n[len(n)-1]
	return strings.Title(name)
}

func getUser(uid string) string {
	u, _ := user.LookupId(uid)
	return u.Username
}

func removeEmpty(array []string) []string {
	// remove empty data from line
	var new_array []string
	for _, i := range array {
		if i != "" {
			new_array = append(new_array, i)
		}
	}
	return new_array
}

func netstat(t string, buffer chan Process) []Process {
	// Return a array of Process with Name, Ip, Port, State .. etc
	// Require Root acess to get information about some processes.

	var Processes []Process

	data := getData(t)

	for _, line := range data {
		// local ip and port
		line_array := removeEmpty(strings.Split(strings.TrimSpace(line), " "))

		ip_port := strings.Split(line_array[1], ":")
		ip := convertIp(ip_port[0])
		port := hexToDec(ip_port[1])

		// foreign ip and port
		fip_port := strings.Split(line_array[2], ":")
		fip := convertIp(fip_port[0])
		fport := hexToDec(fip_port[1])

		state := STATE[line_array[3]]
		//getUser(line_array[7]) // panic 1
		uid, err := strconv.ParseInt(line_array[7], 10, 64)
		if err != nil {
			uid = 0
		}
		//pid := findPid(line_array[9])
		//exe := getProcessExe(pid)
		//name := getProcessName(exe)

		p := Process{uid, "name", "pid", "exe", state, ip, port, fip, fport}

		buffer <- p
		Processes = append(Processes, p)

	}
	close(buffer)

	return Processes
}

func Tcp(buffer chan Process) []Process {
	// Get a slice of Process type with TCP data
	data := netstat("tcp", buffer)
	return data
}

func Udp(buffer chan Process) []Process {
	// Get a slice of Process type with UDP data
	data := netstat("udp", buffer)
	return data
}

func Tcp6(buffer chan Process) []Process {
	// Get a slice of Process type with TCP6 data
	data := netstat("tcp6", buffer)
	return data
}

func Udp6(buffer chan Process) []Process {
	// Get a slice of Process type with UDP6 data
	data := netstat("udp6", buffer)
	return data
}
