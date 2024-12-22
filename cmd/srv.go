/*
Copyright © 2024 superryanguo
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/superryanguo/ryai/docs"
	"github.com/superryanguo/ryai/llm"
	"github.com/superryanguo/ryai/llmapp"
	"github.com/superryanguo/ryai/ollama"
	"github.com/superryanguo/ryai/secret"
	"github.com/superryanguo/ryai/storage"
	//"github.com/superryanguo/ryai/utils"
)

// Ryai own the runtime instance
type Ryai struct {
	ctx       context.Context
	slog      *slog.Logger         // slog output to use
	slogLevel *slog.LevelVar       // slog level, for changing as needed
	http      *http.Client         // http client to use
	addr      string               // address to serve HTTP on
	db        storage.DB           // database to use
	vector    storage.VectorDB     // vector database to use
	secret    secret.DB            // secret database to use
	docs      *docs.Corpus         // document corpus to use
	embed     llm.Embedder         // LLM embedder to use
	llm       llm.ContentGenerator // LLM content generator to use
	llmapp    *llmapp.Client       // LLM client to use
}

func Run() {
	logger.Info("Srv is running...***...***...***...***...***...")

	g := &Ryai{
		ctx:       context.Background(),
		slog:      logger,
		slogLevel: loglevel,
		http:      http.DefaultClient,
		addr:      "localhost:4229",
	}

	var osrv string
	ai, err := ollama.NewClient(g.slog, g.http, osrv, ollama.DefaultEmbeddingModel)
	if err != nil {
		log.Fatal(err)
	}
	g.embed = ai
	//g.llm = ai
	//g.llmapp = llmapp.New(g.slog, ai, g.db)

	var docs = []llm.EmbedDoc{
		{Text: "for loops"},
		{Text: "for all time, always"},
	}

	vecs, err := g.embed.EmbedDocs(g.ctx, docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("vecs:%v\n", vecs)
	input := "how about Donald Trump?"
	gai, err := ollama.NewClient(g.slog, g.http, osrv, ollama.DefaultGenModel)
	if err != nil {
		log.Fatal(err)
	}
	rsp, err := gai.Prompt(g.ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Gai get rsp: %s\n", string(rsp))
	s, err := ollama.AssembleRsp(rsp)
	if err != nil {

	}

	fmt.Printf("AssembleRsp: %s\n", s)

	//utils.ShowJsonRsp(rsp)

	select {}
}

func Chat() {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		t := scanner.Text()
		if t == "exit" {
			break
		}
		chat(t)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
	}
}
func chat(input string) {
	ctx := context.Background()
	ai, err := ollama.NewClient(logger, http.DefaultClient, "", ollama.DefaultGenModel)
	if err != nil {
		log.Fatal(err)
	}
	rsp, err := ai.Prompt(ctx, input)
	if err != nil {
		log.Fatal(err)
	}

	s, err := ollama.AssembleRsp(rsp)
	if err != nil {
	}

	fmt.Printf("\033[32mAnswer: %s\033[0m\n", s)
}
