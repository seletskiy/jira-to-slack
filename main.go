package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/docopt/docopt-go"
)

const usage = `jira-to-slack
	
Missing link between Jira and Slack Integration.
	
Usage:
	$0 -h | --help
	$0 <slack-url> -L <address> -t <text-template> [options]

Options:
    -h --help             Show this help.
    -L                    Listen to specified address.
      -t <text-template>  Template which will be rendered from incoming JSON
                          and sent to Slack.
      -e <emoji>          Emoji to use as icon.
      -c <channel>        Channel to send to.
      -u <username>       Username to show.
    -v                    Show incoming JSON from Jira to stdout.
`

type webHookHandler struct {
	slackURL string
	template *template.Template
	channel  interface{}
	emoji    interface{}
	username interface{}
	debug    bool
}

func (handler *webHookHandler) ServeHTTP(
	writer http.ResponseWriter, request *http.Request,
) {
	jiraBody := map[string]interface{}{}
	err := json.NewDecoder(request.Body).Decode(&jiraBody)
	if err != nil {
		log.Println("error while decoding JSON from WebHook:", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	textBuffer := &bytes.Buffer{}
	err = handler.template.Execute(textBuffer, jiraBody)
	if err != nil {
		log.Println("error while executing template:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if handler.debug {
		debugBody, _ := json.MarshalIndent(jiraBody, "", "  ")
		fmt.Println(string(debugBody))
	}

	slackBody := map[string]interface{}{
		"text": textBuffer.String(),
	}

	if handler.channel != nil {
		slackBody["channel"] = handler.channel.(string)
	}

	if handler.emoji != nil {
		slackBody["icon_emoji"] = handler.emoji.(string)
	}

	if handler.username != nil {
		slackBody["username"] = handler.username.(string)
	}

	slackBuffer := &bytes.Buffer{}

	err = json.NewEncoder(slackBuffer).Encode(slackBody)
	if err != nil {
		log.Println("error while encoding request for Slack:", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = http.Post(handler.slackURL, "application/json", slackBuffer)
	if err != nil {
		log.Println("error while performing request to Slack:", err)
		writer.WriteHeader(http.StatusBadGateway)
	}

	writer.WriteHeader(http.StatusOK)
}

func main() {
	args, err := docopt.Parse(
		strings.Replace(usage, "$0", os.Args[0], -1),
		nil, true, "jira-to-slack 1.0", false,
	)
	if err != nil {
		panic(err)
	}

	switch {
	case args["-L"]:
		textTemplate, err := template.New("text").Parse(args["-t"].(string))
		if err != nil {
			log.Fatalf(err.Error())
		}

		http.Handle("/", &webHookHandler{
			slackURL: args["<slack-url>"].(string),
			template: textTemplate,
			channel:  args["-c"],
			emoji:    args["-e"],
			username: args["-u"],
			debug:    args["-v"].(bool),
		})

		http.ListenAndServe(args["<address>"].(string), nil)
	}
}
