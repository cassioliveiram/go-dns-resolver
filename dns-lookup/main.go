package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	fileName := "hostfiles/hosts.txt"

	hosts, err := readFromFile(fileName)
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}
	fmt.Println("\n")
	fmt.Println("\n")
	fmt.Println("\n")
	fmt.Println("List of hosts to be validated:")
	for _, hostsList := range hosts {
		fmt.Println("-",hostsList)
	}

	fmt.Println("------------------------------")
	fmt.Println("| Starting dns-resolver Tool |")
	fmt.Println("------------------------------")
	for {
		for _, host := range hosts {

			ips, erro := net.LookupIP(host)
			if erro != nil {
				fmt.Println("######################################################")
				fmt.Println("Failure detected")
				log.Println(erro)
				fmt.Println("######################################################")
			}

			for _, ip := range ips {
				log.Println(host, "resolved to", ip)
				time.Sleep(time.Second)
			}

		}
	}
}

func readFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}
