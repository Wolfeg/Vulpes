package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/go-playground/webhooks.v5/gitlab"
)

var statusIcon = map[string]string{
	"created":              "🗒️",
	"waiting_for_resource": "🗒️",
	"preparing":            "🗒️",
	"pending":              "🗒️",
	"running":              "🔄",
	"success":              "✅",
	"failed":               "❌",
	"canceled":             "🚫",
	"skipped":              "↩️",
	"manual":               "⚙️",
	"scheduled":            "⏱️",
}

var watchedStatuses []string = []string{"running", "success", "failed", "canceled"}

var jiraUri = os.Getenv("JIRAURI")
var jiraProjectCode = os.Getenv("JIRAPROJECTCODE")
var gitlabUri = os.Getenv("GITLABURI")
var chatID, _ = strconv.ParseInt(os.Getenv("GITLABTGCHATID"), 10, 64)
var botToken = os.Getenv("GITLABTGTOKEN")
var listenPort = os.Getenv("GITLABTGLISTENPORT")
var gitlabSecretToken = os.Getenv("GITLABTGSECRET")

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func main() {
	http.HandleFunc("/", GitlabServer)
	http.ListenAndServe(fmt.Sprintf(":%v", listenPort), nil)
}

func GitlabServer(w http.ResponseWriter, r *http.Request) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	checkErr(err)
	fmt.Printf("Logged in as %v\n", bot.Self.UserName)
	checkErr(err)
	hook, err := gitlab.New(gitlab.Options.Secret(gitlabSecretToken))
	checkErr(err)
	payload, err := hook.Parse(r, gitlab.PipelineEvents)
	checkErr(err)
	switch pipeline := payload.(type) {
	case gitlab.PipelineEventPayload:
		if contains(watchedStatuses, pipeline.ObjectAttributes.Status) {
			pipelineJSON, err := json.MarshalIndent(pipeline, "", "  ")
			checkErr(err)
			fmt.Printf("Got new pipeline pipeline:\n%s\n", string(pipelineJSON))
			textMsg := fmt.Sprintf("New pipeline [#%v](%v) event\nProject: [%v](%v)\nBranch: [%v](%v)\nAuthor: [%v](%v)\nStatus: `%v %v`",
				pipeline.ObjectAttributes.ID,
				fmt.Sprintf("%v/-/pipelines/%v", pipeline.Project.WebURL, pipeline.ObjectAttributes.ID),
				pipeline.Project.Name,
				pipeline.Project.WebURL,
				pipeline.ObjectAttributes.Ref,
				fmt.Sprintf("%v/-/tree/%v", pipeline.Project.WebURL, pipeline.ObjectAttributes.Ref),
				pipeline.User.Name,
				fmt.Sprintf("%v/%v", gitlabUri, pipeline.User.UserName),
				statusIcon[pipeline.ObjectAttributes.Status],
				pipeline.ObjectAttributes.Status)
			re := regexp.MustCompile(jiraProjectCode + `-\d+`)
			jiraIssue := re.FindString(pipeline.Commit.Message)
			if jiraIssue != "" {
				textMsg = textMsg + fmt.Sprintf("\nJira task: [%v](%v/projects/%v/issues/%v)", jiraIssue, jiraUri, jiraProjectCode, jiraIssue)
			}
			firstMsg := tgbotapi.NewMessage(chatID, textMsg)
			firstMsg.ParseMode = "Markdown"
			firstMsg.DisableWebPagePreview = true
			_, err = bot.Send(firstMsg)
			checkErr(err)
		}
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}
