Missing link between Jira and Slack Integration.

This little program will listen for requests from WebHook from Jira and
translate them into Slack API Requests, which will show Bot messages in
specified topic.

See: https://api.slack.com/incoming-webhooks

# Installation

## Archlinux

Package available via AUR [package jira-to-slack](https://aur4.archlinux.org/packages/jira-to-slack).

## go get

```
go get github.com/seletskiy/jira-to-slack
```

# Usage

Run program as daemon via your favorite manager:

```
jira-to-slack $slackUrl \
    -L :12345 \
    -t '<http://$jiraBaseUrl/{{.issue.key}}|{{.issue.key}}>: {{.issue.fields.summary}}'
    -u jira
    -c '#devops'
    -v
```

Replace `$slackUrl` with slack integration URL
(https://api.slack.com/docs/oauth).

Setup WebHook in your Jira instance that will pointing to address where
`jira-to-slack` is hosting.
