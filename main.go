package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/machinebox/graphql"
	"github.com/peterbourgon/ff"
	"github.com/peterbourgon/ff/ffcli"
)

type graphqlResponse struct {
	User struct {
		Todos []struct {
			Body         string
			Completed_At string
			Created_At   string
		}
	}
}

const (
	graphQLAPIAddress = "https://wip.chat/graphql"
)

var (
	wipUsername string
	apiKey      string
	taskMessage string
)

func main() {
	// cli setup
	globalFlags := flag.NewFlagSet("wip-auto-done", flag.ExitOnError)

	globalFlags.StringVar(&wipUsername, "wip-user", "", "your username on wip.chat")
	globalFlags.StringVar(&apiKey, "api-key", "", "your wip.chat api key (https://wip.chat/api)")
	globalFlags.StringVar(&taskMessage, "message", "#tryhardinglife", "the message body of the completed todo")

	root := &ffcli.Command{
		Usage:   "wip-auto-done [global flags]",
		FlagSet: globalFlags,
		Options: []ff.Option{ff.WithEnvVarNoPrefix()},

		Exec: func(args []string) error {
			// create client and retrieve latest completed task
			graphqlClient := graphql.NewClient(graphQLAPIAddress)
			graphqlRequest := graphql.NewRequest(`
			{
				user(username: "` + wipUsername + `") {
					todos(completed: true, limit: 1) {
						body
						completed_at
						created_at
					}
				}
			}
			`)
			var response graphqlResponse
			err := graphqlClient.Run(context.Background(), graphqlRequest, &response)
			if err != nil {
				return err
			}

			// setup time comparison
			timeFormatLayout := "2006-01-02T15:04:05Z"
			now := time.Now().UTC()
			todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

			// check if latest completed task was today
			for _, todo := range response.User.Todos {
				if todo.Completed_At != "" {
					completedTime, err := time.Parse(timeFormatLayout, todo.Completed_At)
					if err != nil {
						return err
					}
					if completedTime.After(todayStart) {
						// return if we find that the latest task was completed today
						return nil
					}
				}
			}

			// if execution reached this point then it means that task a task needs to be created and completed
			graphqlRequest = graphql.NewRequest(`
				mutation {
					createTodo(input: {body: "` + taskMessage + `", completed_at: "` + now.Format(timeFormatLayout) + `"}) {
						id
					}
				}
				`)
			var resp interface{}
			graphqlRequest.Header.Add("Authorization", fmt.Sprintf("bearer %s", apiKey))
			err = graphqlClient.Run(context.Background(), graphqlRequest, &resp)
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("Couldn't create the todo, check that your api key is valid")
			}
			fmt.Println("Added done task with following message: " + taskMessage)
			return nil
		},
	}

	// cli stuff, basically print error in console and display help if needed
	if err := root.Run(os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatalf("fatal: %+v", err)
	}
}
