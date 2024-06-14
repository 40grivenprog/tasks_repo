package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	ALLOW string = "ALLOW"
)

func IsAllowed(rules [][2]string, ip string) bool {
	for _, rule := range rules {
		if rule[1] == ALLOW {
			if isIPAllowed(rule[0], ip) {
				return true
			}
		}
	}
	return false
}

func isIPAllowed(ruleIP string, ip string) bool {
	pattern, err := getCommonIpPattern(ruleIP)
	if err != nil {
		fmt.Println("Error parsing rule IP:", err)
		return false
	}

	comparableValue := string([]rune(ip)[0:len(pattern)])
	fmt.Println(pattern, comparableValue)
	return comparableValue == pattern
}

func getCommonIpPattern(ip string) (string, error) {
	splittedIp := strings.Split(ip, "/")
	if len(splittedIp) == 1 {
		return splittedIp[0], nil
	}

	ipPart := splittedIp[0]
	maskPart := splittedIp[1]

	mask, err := strconv.Atoi(maskPart)
	if err != nil {
		return "", err
	}

	_, ipNet, err := net.ParseCIDR(ipPart + "/" + maskPart)
	if err != nil {
		return "", err
	}

	ipScope := mask / 8
	commonIpPattern := strings.Join(strings.Split(ipNet.IP.String(), ".")[0:int(ipScope)], ".")

	return commonIpPattern, nil
}

func main() {
	rules := [][2]string{
		{"192.168.31.0/24", "ALLOW"},
		{"10.0.0.0/16", "DENY"},
		{"1.2.3.4/16", "ALLOW"},
	}

	fmt.Println(IsAllowed(rules, "1.2.3.4"))      // true
	fmt.Println(IsAllowed(rules, "192.168.31.5")) // true
	fmt.Println(IsAllowed(rules, "10.0.1.1"))     // false
}
