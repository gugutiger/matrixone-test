// Copyright 2022 Matrix Origin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package multi

import (
	"github.com/matrixorigin/matrixone/pkg/container/nulls"
	"github.com/matrixorigin/matrixone/pkg/container/types"
	"github.com/matrixorigin/matrixone/pkg/container/vector"
	"github.com/matrixorigin/matrixone/pkg/vm/mheap"
	"github.com/matrixorigin/matrixone/pkg/vm/process"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLpadVarchar(t *testing.T) {
	cases := []struct {
		name      string
		vecs      []*vector.Vector
		proc      *process.Process
		wantBytes []byte
	}{
		{
			name:      "TEST01",
			vecs:      makeLpadVectors("hello", 1, "#", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("h"),
		},
		{
			name:      "TEST02",
			vecs:      makeLpadVectors("hello", 10, "#", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("#####hello"),
		},
		{
			name:      "TEST03",
			vecs:      makeLpadVectors("hello", 15, "#@&", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("#@&#@&#@&#hello"),
		},
		{
			name:      "TEST04",
			vecs:      makeLpadVectors("12345678", 10, "abcdefgh", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("ab12345678"),
		},
		{
			name:      "TEST05",
			vecs:      makeLpadVectors("hello", 0, "#@&", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(""),
		},
		{
			name:      "TEST06",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "Tx",
			vecs:      makeLpadVectors("hello", 1, "", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(""),
		},
		{
			name:      "Tx2",
			vecs:      makeLpadVectors("", 5, "x", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("xxxxx"),
		},
		{
			name:      "Tx3",
			vecs:      makeLpadVectors("??????", 10, "??????", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("??????????????????????????????"),
		},
		{
			name:      "tx4",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{0, 0, 0}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx5",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{0, 0, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{0, 1, 0}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{1, 0, 0}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{1, 1, 0}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{1, 0, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{0, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("hello", -1, "#@&", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte(nil),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("a???", 15, "???", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("???????????????????????????????????????a???"),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("a???a", 15, "???a", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("???a???a???a???a???a???aa???a"),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("a???aa", 15, "???aa", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("???aa???aa???aa???aa???aa"),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("a???aaa", 15, "???aaa", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("???aaa???aaa???aa???aaa"),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("a???aaaa", 15, "???aaaa", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("???aaaa???aaaa???aaaa"),
		},
		{
			name:      "tx6",
			vecs:      makeLpadVectors("aaaaaaaa", 4, "bbb", []int{1, 1, 1}),
			proc:      process.New(mheap.New(nil)),
			wantBytes: []byte("aaaa"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lpad, err := Lpad(c.vecs, c.proc)
			if err != nil {
				t.Fatal(err)
			}
			if c.wantBytes == nil {
				ret := nulls.Contains(lpad.Nsp, 0)
				require.Equal(t, ret, true)
			} else {
				require.Equal(t, c.wantBytes, lpad.Col.(*types.Bytes).Data)
			}

		})
	}

}

func makeLpadVectors(src string, length int64, pad string, nils []int) []*vector.Vector {
	vec := make([]*vector.Vector, 3)
	srcBytes := &types.Bytes{
		Data:    []byte(src),
		Offsets: []uint32{0},
		Lengths: []uint32{uint32(len(src))},
	}
	padBytes := &types.Bytes{
		Data:    []byte(pad),
		Offsets: []uint32{0},
		Lengths: []uint32{uint32(len(pad))},
	}

	vec[0] = &vector.Vector{
		Col:     srcBytes,
		Nsp:     &nulls.Nulls{},
		Typ:     types.Type{Oid: types.T_varchar, Size: 24},
		IsConst: true,
		Length:  1,
	}

	vec[1] = &vector.Vector{
		Col:     []int64{length},
		Nsp:     &nulls.Nulls{},
		Typ:     types.Type{Oid: types.T_int64},
		IsConst: true,
		Length:  1,
	}

	vec[2] = &vector.Vector{
		Col:     padBytes,
		Nsp:     &nulls.Nulls{},
		Typ:     types.Type{Oid: types.T_varchar, Size: 24},
		IsConst: true,
		Length:  1,
	}
	for i, n := range nils {
		if n == 0 {
			nulls.Add(vec[i].Nsp, uint64(i))
		}
	}
	return vec
}
