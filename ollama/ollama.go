// Package ollama implements access to offline Ollama model.
//
// [Client] implements [llm.Embedder]. Use [NewClient] to connect.
package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/superryanguo/ryai/llm"
)

const (
	DefaultEmbeddingModel = "mxbai-embed-large"
	EmbedUrl              = "/api/embed"
	DefaultGenModel       = "llama3.2:3b"
	GenUrl                = "/api/generate"
	DefaultGenModel2      = "llama3"
	maxBatch              = 512 // default physical batch size in ollama
)

// A Client represents a connection to Ollama.
type Client struct {
	slog  *slog.Logger
	hc    *http.Client
	url   *url.URL // url of the ollama server
	model string
}

type Response struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func AssembleRsp(d []byte) (string, error) {
	var s string
	var err error

	lines := strings.Split(string(d), "\n")
	for _, line := range lines {
		var resp Response
		err = json.Unmarshal([]byte(line), &resp)
		if err != nil {
			return s, err
		}
		s += resp.Response
	}
	return s, nil
}

// NewClient returns a connection to Ollama server. If empty, the
// server is assumed to be hosted at http://127.0.0.1:11434.
// The model is the model name to use for embedding.
// A typical model for embedding is "mxbai-embed-large".
func NewClient(lg *slog.Logger, hc *http.Client, server string, model string) (*Client, error) {
	if server == "" {
		host := os.Getenv("OLLAMA_HOST")
		if host == "" {
			host = "127.0.0.1"
		}
		server = "http://" + host + ":11434"
	}
	u, err := url.Parse(server)
	if err != nil {
		return nil, err
	}
	return &Client{slog: lg, hc: hc, url: u, model: model}, nil
}

// EmbedDocs returns the vector embeddings for the docs,
// implementing [llm.Embedder].
func (c *Client) EmbedDocs(ctx context.Context, docs []llm.EmbedDoc) ([]llm.Vector, error) {
	embedURL := c.url.JoinPath(EmbedUrl) // ollama embed endpoint
	var vecs []llm.Vector
	for docs := range slices.Chunk(docs, maxBatch) {
		var inputs []string
		for _, doc := range docs {
			// ollama does not support adding content with title
			input := doc.Title + "\n\n" + doc.Text
			inputs = append(inputs, input)
		}
		vs, err := embed(ctx, c.hc, embedURL, inputs, c.model)
		if err != nil {
			return nil, err
		}
		vecs = append(vecs, vs...)
	}
	return vecs, nil
}

func (c *Client) Prompt(ctx context.Context, input string) ([]byte, error) {
	u := c.url.JoinPath(GenUrl)
	rsp, err := prompt(ctx, c.hc, u, input, c.model)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func prompt(ctx context.Context, hc *http.Client, u *url.URL, inputs string, model string) ([]byte, error) {
	ask := struct {
		Model  string `json:"model"`
		Input  string `json:"prompt"`
		Stream bool   `json:"stream"`
	}{
		Model:  model,
		Input:  inputs,
		Stream: true,
	}

	erj, err := json.Marshal(ask)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewReader(erj))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := hc.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	rsp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

func embed(ctx context.Context, hc *http.Client, embedURL *url.URL, inputs []string, model string) ([]llm.Vector, error) {
	embReq := struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}{
		Model: model,
		Input: inputs,
	}
	erj, err := json.Marshal(embReq)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, embedURL.String(), bytes.NewReader(erj))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := hc.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	embResp, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if err := embedError(response, embResp); err != nil {
		return nil, err
	}
	return embeddings(embResp)
}

// embedError extracts error from ollama's response, if any.
func embedError(resp *http.Response, embResp []byte) error {
	if resp.StatusCode == 200 {
		return nil
	}
	if resp.StatusCode == 400 {
		var e struct {
			Error string `json:"error"`
		}
		// ollama returns JSON with error field set for bad requests.
		if err := json.Unmarshal(embResp, &e); err != nil {
			return err
		}
		return fmt.Errorf("ollama response error: %s", e.Error)
	}
	return fmt.Errorf("ollama response error: %s", resp.Status)
}

func embeddings(embResp []byte) ([]llm.Vector, error) {
	// In case there are no errors, ollama returns
	// a JSON with Embeddings field set.
	var e struct {
		Embeddings []llm.Vector `json:"embeddings"`
	}
	if err := json.Unmarshal(embResp, &e); err != nil {
		return nil, err
	}
	return e.Embeddings, nil
}
