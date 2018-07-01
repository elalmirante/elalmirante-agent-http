package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"gitlab.com/ozkar99/middleware"
	"gopkg.in/yaml.v2"
)

type conf struct {
	Port   string
	User   string
	Pass   string
	Script []string
}

func (c conf) ScriptLine() string {
	return strings.Join(c.Script, " && ")
}

func main() {

	// read config file
	var confPath string
	flag.StringVar(&confPath, "c", "/etc/elalmirante-agent.conf", "The path to the configuration file.")
	flag.Parse()

	confBytes, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal(err)
	}

	// unmarshall config
	var conf *conf
	err = yaml.Unmarshal(confBytes, &conf)
	if err != nil {
		log.Fatal(err)
	}

	// http server

	authPair := conf.User + ":" + conf.Pass
	http.Handle("/deploy",
		middleware.BasicAuth(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// only accept post
				if r.Method != http.MethodPost {
					log.Println("Method not supported.")
					return
				}

				log.Println("Deploying...")

				cmd := exec.Command("/bin/sh", "-c", conf.ScriptLine())
				outputBytes, err := cmd.CombinedOutput()
				output := fmt.Sprintf("%s\n\n Error: %v\n", string(outputBytes), err)

				log.Println(output)

				if err != nil {
					http.Error(w, output, http.StatusInternalServerError)
					return
				}

				fmt.Fprintf(w, output)
				log.Println("Finished deploy.")

			}), authPair))

	bind := ":" + conf.Port
	http.ListenAndServe(bind, nil)
}
