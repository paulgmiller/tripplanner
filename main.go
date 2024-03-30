package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/v60/github"
	openai "github.com/sashabaranov/go-openai"
)

const (
	owner  = "paulgmiller"
	repo   = "tripplanner"
	branch = "master"
)

func main() {
	githubpat := os.Getenv("GITHUB_PAT")
	openaiat := os.Getenv("OPENAI_AT")

	tp := &tripplanner{
		github.NewClient(nil).WithAuthToken(githubpat),
		openai.NewClient(openaiat),
	}

	http.Handle("/", tp)
	http.ListenAndServe("", nil)

}

type tripplanner struct {
	ghclient     *github.Client
	openaiclient *openai.Client
}

func (tp *tripplanner) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	ctx := context.Background()
	path := req.URL.Path
	fileContent, _, _, err := tp.ghclient.Repositories.GetContents(ctx, owner, repo, path, &github.RepositoryContentGetOptions{Ref: branch})
	if err != nil {
		log.Printf("can't find gh file %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	content, err := fileContent.GetContent()
	if err != nil {
		log.Printf("can't getch gh file %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	airesp, err := tp.openaiclient.CreateChatCompletion(
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
		log.Printf("ChatCompletion error: %v\n", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(airesp.Choices[0].Message.Content)

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Update file with new content"),
		Content: []byte(content),
		SHA:     github.String("DEADBEEFDEADBEEF"),
		Branch:  github.String(branch),
	}

	_, _, err = tp.ghclient.Repositories.UpdateFile(ctx, owner, repo, path, opts)
	if err != nil {
		log.Printf("can't update gh file %s", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

}
