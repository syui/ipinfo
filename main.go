package main

import (
	"os"
	"fmt"
	"log"
	"bytes"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"github.com/urfave/cli"
	"github.com/hokaccha/go-prettyjson"
	forecast "github.com/mlbright/forecast/v2"
)

var exit = os.Exit
func doSomething() {
    exit(0)
}

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

var apikey string

func fetch(url string) []byte {
  res, err := http.Get(url)
  fmt.Println(url)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  body, err := ioutil.ReadAll(res.Body)
  if err != nil {
    log.Fatal(err)
  }
  return body
}

func main() {

	var b Ipinfo
	var f string

	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(0)
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	app := cli.NewApp()
	app.Version = "0.1"

	if len(os.Args) == 1 {
		fmt.Println(strings.Trim(newStr, "\n"))
		os.Exit(0)
	}

	url := fmt.Sprintf("http://ipinfo.io/%s/json",strings.Trim(newStr, "\n"))
	req, err := http.Get(url)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(0)
	}
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(&b); err != nil {
		log.Println(err)
	}

	loc := strings.Split(b.Loc, ",")
	lat := loc[0]
	long := loc[1]

	app.Commands = []cli.Command{
		{
			Name:    "ipinfo",
			Aliases: []string{"i"},
			Usage:   "ipinfo.io\n\t\tsub-command : ip(i), country(c), city(ci), loc(l), org(o)",
			Action:  func(c *cli.Context) error {
				jb, _ := prettyjson.Marshal(b)
				fmt.Printf("%s", jb)
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "ip",
					Usage:   "ipifo i i",
					Aliases: []string{"i"},
					Action:  func(c *cli.Context) error {
						fmt.Println(b.IP)
						return nil
					},
				},
				cli.Command{
					Name:   "country",
					Usage:   "ipifo i c",
					Aliases: []string{"c"},
					Action:  func(c *cli.Context) error {
						fmt.Println(b.Country)
						return nil
					},
				},
				cli.Command{
					Name:   "city",
					Usage:   "ipifo i ci",
					Aliases: []string{"ci"},
					Action:  func(c *cli.Context) error {
						fmt.Println(b.City)
						return nil
					},
				},
				cli.Command{
					Name:   "loc",
					Usage:   "ipifo i l",
					Aliases: []string{"l"},
					Action:  func(c *cli.Context) error {
						fmt.Println(b.Loc)
						return nil
					},
				},
				cli.Command{
					Name:   "org",
					Usage:   "ipifo i o",
					Aliases: []string{"o"},
					Action:  func(c *cli.Context) error {
						fmt.Println(b.Org)
						return nil
					},
				},
			},
		},
		{
			Name:    "ip",
			Usage:   "grobal ip address",
			Action:  func(c *cli.Context) error {
				fmt.Println(b.IP)
				return nil
			},
		},
		{
			Name:    "loc",
			Usage:   "location gps",
			Action:  func(c *cli.Context) error {
				fmt.Println(b.Loc)
				return nil
			},
		},
		{
			Name:    "sky",
			Aliases: []string{"s"},
			Usage:   "darksky.io\n\t\tsub-command : info(i)",
			Action:  func(c *cli.Context) error {
				f, err := forecast.Get(apikey, lat, long, "now", forecast.CA, forecast.English)
				if err != nil {
					log.Fatal(err)
				}
				jf, _ := prettyjson.Marshal(f)
				fmt.Printf("%s", jf)
				return nil
			},
			Subcommands: cli.Commands{
				cli.Command{
					Name:   "info",
					Usage:   "info",
					Aliases: []string{"i"},
					Action:  func(c *cli.Context) error {
						f, err := forecast.Get(apikey, lat, long, "now", forecast.CA, forecast.English)
						if err != nil {
							log.Fatal(err)
						}
						fmt.Printf("%s: %s\n", f.Timezone, f.Currently.Summary)
						fmt.Printf("humidity: %.2f\n", f.Currently.Humidity)
						fmt.Printf("temperature: %.2f Celsius\n", f.Currently.Temperature)
						fmt.Printf("wind speed: %.2f\n", f.Currently.WindSpeed)
						return nil
					},
				},
			},
		},
		{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "json",
			Action:  func(c *cli.Context) error {
				key := apikey
				loc := strings.Split(b.Loc, ",")
				lat := loc[0]
				long := loc[1]

				ff, err := forecast.Get(key, lat, long, "now", forecast.CA, forecast.English)
				if err != nil {
					log.Fatal(err)
				}
				jf, _ := prettyjson.Marshal(ff)
				jb, _ := prettyjson.Marshal(b)
				fmt.Printf("[\n%s,", jb)
				fmt.Printf("%s\n]", jf)
				fmt.Println(f)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

