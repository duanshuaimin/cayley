// Copyright 2014 The Cayley Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package iterator

import (
	"errors"
	"testing"
)

func TestMaterializeIteratorError(t *testing.T) {
	wantErr := errors.New("unique")
	errIt := newTestIterator(false, wantErr)

	// This tests that we properly return 0 results and the error when the
	// underlying iterator returns an error.
	mIt := NewMaterialize(errIt)

	if mIt.Next() != false {
		t.Errorf("Materialize iterator did not pass through underlying 'false'")
	}
	if mIt.Err() != wantErr {
		t.Errorf("Materialize iterator did not pass through underlying Err")
	}
}

func TestMaterializeIteratorErrorAbort(t *testing.T) {
	wantErr := errors.New("unique")
	errIt := newTestIterator(false, wantErr)

	// This tests that we properly return 0 results and the error when the
	// underlying iterator is larger than our 'abort at' value, and then
	// returns an error.
	or := NewOr()
	or.AddSubIterator(NewInt64(1, int64(abortMaterializeAt+1), true))
	or.AddSubIterator(errIt)

	mIt := NewMaterialize(or)

	// We should get all the underlying values...
	for i := 0; i < abortMaterializeAt+1; i++ {
		if !mIt.Next() {
			t.Errorf("Materialize iterator returned spurious 'false' on iteration %d", i)
			return
		}
		if mIt.Err() != nil {
			t.Errorf("Materialize iterator returned non-nil Err on iteration %d", i)
			return
		}
	}

	// ... and then the error value.
	if mIt.Next() != false {
		t.Errorf("Materialize iterator did not pass through underlying 'false'")
	}
	if mIt.Err() != wantErr {
		t.Errorf("Materialize iterator did not pass through underlying Err")
	}
}
