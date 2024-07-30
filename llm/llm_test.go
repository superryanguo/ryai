// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package llm

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"testing"
)

func TestVector(t *testing.T) {
	v1 := Vector{1, 2, 3, 4}
	v2 := Vector{-200, -3000, 0, -10000}
	dot := v1.Dot(v2)
	if dot != -46200 {
		t.Errorf("%v.Dot(%v) = %v, want -46200", v1, v2, dot)
	}

	enc := v1.Encode()
	var v3 Vector
	v3.Decode(enc)
	if !slices.Equal(v3, v1) {
		t.Errorf("Decode(Encode(%v)) = %v, want %v", v1, v3, v1)
	}
}
func TestFloat(t *testing.T) {
	v := Vector{1, 2, 3, 4, 11, 15, 17, 19}
	enc := v.Encode()
	fmt.Printf("v enc =%02X\n", enc)
	for _, vi := range enc {
		fmt.Printf("%02X-", vi)
		//fmt.Printf("k =%d, vi =%02X\n", k, vi)
	}
	fmt.Println("")

	vp := make(Vector, len(enc)/4)
	//var p *Vector
	//p := &vp
	//*p = (*p)[:0]
	vp = vp[:0]
	//p := &vp

	for ; len(enc) >= 4; enc = enc[4:] {
		fmt.Printf("[]vector p=%v\n", vp)
		fmt.Printf("decode before enc =%02X\n", enc)
		vp = append(vp, math.Float32frombits(binary.BigEndian.Uint32(enc)))
		fmt.Printf("decode after  enc =%02X\n", enc)
	}
}
