package main

import (
	"fmt"
	"strings"
	"log"
	"bytes"
	"os"
	"encoding/json"
	"net/http"
	"io/ioutil"
	forecast "github.com/mlbright/forecast/v2"
	"github.com/urfave/cli"
)

// Ipinfo golint
type Ipinfo struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Region   string `json:"region"`
}

func main() {

	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()
	//io.Copy(os.Stdout, resp.Body)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	url := fmt.Sprintf("http://ipinfo.io/%s/json",strings.Trim(newStr, "\n"))
	req, err := http.Get(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	defer req.Body.Close()

	var b Ipinfo
	if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
		log.Println(err)
	}

	keybytes, err := ioutil.ReadFile("api_key.txt")
	if err != nil {
	    log.Fatal(err)
	}
	key := string(keybytes)
	//key := string("API KEY")
	key = strings.TrimSpace(key)

	loc := strings.Split(b.Loc, ",")
	lat := loc[0]
	long := loc[1]

	f, err := forecast.Get(key, lat, long, "now", forecast.CA, forecast.English)
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	//app.Name = "darksky.go"
	//app.Usage = "fight the loneliness!"
	//app.Action = func(c *cli.Context) error {
	//  fmt.Println(b, f)
	//  return nil
	//}

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "lang, l",
			Value: "english",
			Usage: "Language for the greeting",
		},
		cli.StringFlag{
			Name: "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "ipinfo",
			Aliases: []string{"i"},
			Usage:   "complete a task on the list",
			Action:  func(c *cli.Context) error {
				fmt.Println("City 	= ",	b.City)
				fmt.Println("IP   	= ",	b.IP)
				fmt.Println("Loc  	= ",	b.Loc)
				fmt.Println("Country 	= ",	b.Country)
				fmt.Println("Hostname 	= ",	b.Hostname)
				fmt.Println("Org 	= ",	b.Org)
				fmt.Println("Region 	= ",	b.Region)
				return nil
			},
		},
		{
			Name:    "sky",
			Aliases: []string{"s"},
			Usage:   "add a task to the list",
			Action:  func(c *cli.Context) error {
				fmt.Printf("%s: %s\n", f.Timezone, f.Currently.Summary)
				fmt.Printf("humidity: %.2f\n", f.Currently.Humidity)
				fmt.Printf("temperature: %.2f Celsius\n", f.Currently.Temperature)
				fmt.Printf("wind speed: %.2f\n", f.Currently.WindSpeed)
				return nil
			},
		},
		{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action:  func(c *cli.Context) error {
				fmt.Println(b, f)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
