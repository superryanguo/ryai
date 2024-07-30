package llm

import (
	"encoding/binary"
	"math"
)

type Embedder interface {
	EmbedDocs(docs []EmbedDoc) ([]Vector, error)
}

// An EmbedDoc is a single document to be embedded.
type EmbedDoc struct {
	Title string // title of document
	Text  string // text of document
}

// A Vector is an embedding vector, typically a high-dimensional unit vector.
type Vector []float32

// Dot returns the dot product of v and w.
func (v Vector) Dot(w Vector) float64 {
	v = v[:min(len(v), len(w))]
	w = w[:len(v)] // make "i in range for v" imply "i in range for w" to remove bounds check in loop
	t := float64(0)
	for i := range v {
		t += float64(v[i]) * float64(w[i])
	}
	return t
}

// Encode returns a byte encoding of the vector v,
// suitable for storing in a database.
func (v Vector) Encode() []byte {
	val := make([]byte, 4*len(v))
	for i, f := range v {
		binary.BigEndian.PutUint32(val[4*i:], math.Float32bits(f))
	}
	return val
}

// Decode decodes the byte encoding enc into the vector v.
// Enc should be a multiple of 4 bytes; any trailing bytes are ignored.
func (v *Vector) Decode(enc []byte) {
	if len(*v) < len(enc)/4 {
		*v = make(Vector, len(enc)/4)
	}
	*v = (*v)[:0]
	for ; len(enc) >= 4; enc = enc[4:] {
		*v = append(*v, math.Float32frombits(binary.BigEndian.Uint32(enc)))
	}
}
