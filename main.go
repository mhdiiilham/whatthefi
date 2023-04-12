package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	ip   = flag.Bool("ip", false, "what is my ip address?")
	down = flag.String("down", "", "is this website down for just me?")
	help = flag.Bool("help", false, "print available command")
)

const (
	WHAT_MY_IP        = "https://api.my-ip.io/ip.json"
	IS_IT_DOWN_FOR_ME = "https://monitor-api.vercel.app/api/public?url=%s"
)

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *ip {
		ipAddress, err := getMyIPAdress()
		if err != nil {
			log.Fatalf("failed to get ip addres: %v", err)
		}
		fmt.Println(ipAddress)
	}

	if *down != "" {
		isDown, err := isItDownCheck(*down)
		if err != nil {
			log.Fatalf("failed to check if the website %s is down: %v", *down, err)
		}

		if !isDown {
			fmt.Printf("%s seems down according to people.\n", *down)
		} else {
			fmt.Printf("%s seems up according to people.\n", *down)
		}
	}

	fmt.Println("thank you!")
}

func getMyIPAdress() (ipAddress string, err error) {
	fmt.Println("getting your ip adress ...")
	httpResp := struct {
		Success bool   `json:"success"`
		IP      string `json:"ip"`
		Type    string `json:"type"`
	}{}

	resp, err := http.Get(WHAT_MY_IP)
	if err != nil {
		return "", err
	}

	if err := json.NewDecoder(resp.Body).Decode(&httpResp); err != nil {
		return "", err
	}

	return fmt.Sprintf("IP: %s Type: %s", httpResp.IP, httpResp.Type), nil
}

func isItDownCheck(url string) (bool, error) {
	fmt.Printf("checking is %s is down ...\n", url)
	httpResp := struct {
		IsDown bool `json:"isDown"`
	}{}

	resp, err := http.Get(fmt.Sprintf(IS_IT_DOWN_FOR_ME, url))
	if err != nil {
		return httpResp.IsDown, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&httpResp); err != nil {
		return httpResp.IsDown, err
	}

	return httpResp.IsDown, nil
}
