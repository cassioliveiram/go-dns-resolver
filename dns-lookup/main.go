package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	fmt.Println("----------------------------------")
	fmt.Println("| Reading resolv.conf settings:  |")
	fmt.Println("----------------------------------")
	resolvConf, err := os.ReadFile("/etc/resolv.conf")

	if err != nil {
		fmt.Println("Error to read %s", resolvConf)
	}
	fmt.Println(string(resolvConf))

	fmt.Println("----------------------------------")
	fmt.Println("|    Reading linux hosts file:   |")
	fmt.Println("----------------------------------")

	hostsFile, err := os.ReadFile("/etc/hosts")

	if err != nil {
		fmt.Println("Error to read %s", hostsFile)
	}
	fmt.Println(string(hostsFile))

	// read host list from file
	fileName := "hostfiles/hosts.txt"

	hosts, err := readFromFile(fileName)
	if err != nil {
		fmt.Println("Error reading from file:", err)
		return
	}
	fmt.Println("----------------------------------")
	fmt.Println("| List of hosts to be validated: |")
	fmt.Println("----------------------------------")

	for _, hostsList := range hosts {
		fmt.Println("-",hostsList)
	}

	fmt.Println("\n")
	fmt.Println("----------------------------------")
	fmt.Println("|   Starting dns-resolver Tool   |")
	fmt.Println("----------------------------------")
	fmt.Println("\n")

	for {
		for _, host := range hosts {

			ips, erro := net.LookupIP(host)
			if erro != nil {
				fmt.Println("######################################################")
				fmt.Println("Failure detected")
				log.Println(erro)
				sendAlertOnSlack(host, erro)
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

var (
	slack_webhook = os.Getenv("slack_webhook")
)

func sendAlertOnSlack(host string, erro error) {
	client := &http.Client{}

	// convert error in string to be used in the payload
	errString := erro.Error()

	data := strings.NewReader(
		`{
    "text": "The DNS Resolver tool found a issue in your environment",
    "blocks": [
    	{
    		"type": "section",
    		"text": {
    			"type": "mrkdwn",
    			"text": "The DNS Resolver tool found a issue:"
    		}
    	},
    	{
    		"type": "section",
    		"text": {
    			"type": "mrkdwn",
    			"text": "The lookup to: ` + host + ` have returned this error: \n`+ errString+`"
    		},
    		"accessory": {
    			"type": "image",
    			"image_url": "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQBM8z1py0lFRs4obVF-bF4h5Ct4oFIVJgHug&usqp=CAU",
    			"alt_text": "maintenance"
    		}
    	},
    ]
}`)
	req, err := http.NewRequest("POST", slack_webhook, data)
	if err != nil {
		log.Fatal(err)
	}
	//req.Header.Set("Authorization", zoom_authorization)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Alert sent to slack channel - %s\n", bodyText)
	fmt.Println("######################################################")
}
