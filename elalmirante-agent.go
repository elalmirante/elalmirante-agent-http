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

func (c conf) ScriptLine(ref string) string {
	cmd := strings.Join(c.Script, " && ")
	cmd = strings.Replace(cmd, "$REF", ref, -1)
	cmd = strings.Replace(cmd, "${REF}", ref, -1)
	return cmd
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
				ref := r.URL.Query().Get("ref")

				// only accept post and ref cannot be blank
				if r.Method != http.MethodPost || ref == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				log.Printf("Deploying	REF: %s ...", ref)

				cmd := exec.Command("/bin/sh", "-c", conf.ScriptLine(ref))
				outputBytes, err := cmd.CombinedOutput()
				output := fmt.Sprintf("\n%s\n\n Error: %v\n", string(outputBytes), err)

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
