/*
Copyright © 2024 superryanguo
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/superryanguo/ryai/config"
	"github.com/superryanguo/ryai/docs"
	"github.com/superryanguo/ryai/llm"
	"github.com/superryanguo/ryai/llmapp"
	"github.com/superryanguo/ryai/ollama"
	"github.com/superryanguo/ryai/secret"
	"github.com/superryanguo/ryai/storage"
	"log"
	"log/slog"
	"net/http"
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
	logger.Info("Srv is running...***...")
	cfg, err := config.ReadCfg()
	if err != nil {
		logger.Error("Readconfig", "error", err)
		return
	}

	fmt.Printf("RyaiConfig:%s", cfg)

	level := new(slog.LevelVar)
	if err = level.UnmarshalText([]byte(cfg.Log.Level)); err != nil {
		log.Fatal(err)
	}

	g := &Ryai{
		ctx: context.Background(),
		//slog:      slog.New(gcphandler.New(level)),
		slogLevel: level,
		http:      http.DefaultClient,
		addr:      "localhost:4229",
	}

	ai, err := ollama.NewClient(g.ctx, g.slog, g.http, g.addr, ollama.DefaultEmbeddingModel)
	if err != nil {
		log.Fatal(err)
	}
	g.embed = ai
	//g.llm = ai
	//g.llmapp = llmapp.New(g.slog, ai, g.db)

}
