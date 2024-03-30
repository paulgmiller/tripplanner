package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v60/github"
	openai "github.com/sashabaranov/go-openai"
)

const (
	owner = "paulgmiller"
	repo  = "tripplanner"
)

func main() {
	githubpat := os.Getenv("GITHUB_PAT")
	openaiat := os.Getenv("OPENAI_AT")

	//prometheus.http
	ghclient := github.NewClient(nil).WithAuthToken(githubpat)
	ctx := context.Background()
	fileContent, _, _, err := ghclient.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: "master"})
	if err != nil {
		// handle error
	}

	content, err := fileContent.GetContent()
	if err != nil {
		// handle error
	}

	http.NewServeMux().Handle("/")

	openaiclient := openai.NewClient(openaiat)
	resp, err := openaiclient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}

type tripplanner struct {
}

//
