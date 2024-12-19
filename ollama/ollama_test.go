// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ollama

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/superryanguo/ryai/httprr"
	"github.com/superryanguo/ryai/llm"
	"github.com/superryanguo/ryai/testutil"
)

var docs = []llm.EmbedDoc{
	{Text: "for loops"},
	{Text: "for all time, always"},
	{Text: "break statements"},
	{Text: "breakdancing"},
	{Text: "forever could never be long enough for me"},
	{Text: "the macarena"},
}

var matches = map[string]string{
	"for loops":            "break statements",
	"for all time, always": "forever could never be long enough for me",
	"breakdancing":         "the macarena",
}

func init() {
	for k, v := range matches {
		matches[v] = k
	}
}

func newTestClient(t *testing.T, rrfile string) *Client {
	check := testutil.Checker(t)
	lg := testutil.Slogger(t)

	rr, err := httprr.Open(rrfile, http.DefaultTransport)
	check(err)

	c, err := NewClient(lg, rr.Client(), "", "mxbai-embed-large")
	check(err)

	return c
}

func TestEmbedBatch(t *testing.T) {
	ctx := context.Background()
	check := testutil.Checker(t)
	c := newTestClient(t, "testdata/embedbatch.httprr")
	vecs, err := c.EmbedDocs(ctx, docs)
	check(err)
	if len(vecs) != len(docs) {
		t.Fatalf("len(vecs) = %d, but len(docs) = %d", len(vecs), len(docs))
	}

	var buf bytes.Buffer
	for i := range docs {
		for j := range docs {
			fmt.Fprintf(&buf, " %.4f", vecs[i].Dot(vecs[j]))
		}
		fmt.Fprintf(&buf, "\n")
	}

	for i, d := range docs {
		best := ""
		bestDot := 0.0
		for j := range docs {
			if dot := vecs[i].Dot(vecs[j]); i != j && dot > bestDot {
				best, bestDot = docs[j].Text, dot
			}
		}
		if best != matches[d.Text] {
			if buf.Len() > 0 {
				t.Errorf("dot matrix:\n%s", buf.String())
				buf.Reset()
			}
			t.Errorf("%q: best=%q, want %q", d.Text, best, matches[d.Text])
		}
	}
}

func TestBigBatch(t *testing.T) {
	ctx := context.Background()
	check := testutil.Checker(t)
	c := newTestClient(t, "testdata/bigbatch.httprr")
	var docs []llm.EmbedDoc

	for i := range 1025 {
		docs = append(docs, llm.EmbedDoc{Text: fmt.Sprintf("word%d", i)})
	}
	docs = docs[:1025]
	vecs, err := c.EmbedDocs(ctx, docs)
	check(err)
	if len(vecs) != len(docs) {
		t.Fatalf("len(vecs) = %d, but len(docs) = %d", len(vecs), len(docs))
	}
}
