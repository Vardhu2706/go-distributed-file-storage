package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "graduationday"
	pathKey := CASPathTransformFunc(key)
	
	expectedFilename := "c046a7dac5cfeddbd09ffca33b65a3d3638bf2e0"
	expectedPathName := "c046a/7dac5/cfedd/bd09f/fca33/b65a3/d3638/bf2e0"

	if pathKey.Pathname != expectedPathName {
		t.Errorf("Have %s want %s", pathKey.Pathname, expectedPathName)
	}

	if pathKey.Filename != expectedFilename {
		t.Errorf("Have %s want %s", pathKey.Filename, expectedFilename)
	}
}

func TestStore(t *testing.T) {

	s := newStore()
	id := generateID()

	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpeg bytes")
	
		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}
	
		if ok := s.Has(id, key); !ok {
			t.Errorf("Expected to have key %s", key)
		}
	
		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}
	
		b, _ := ioutil.ReadAll(r)
		
		if string(b)!= string(data) {
			t.Errorf("Want %s have %s", data, b)
		}
	
		if err := s.Delete(id, key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); ok {
			t.Errorf("Expected to NOT have key: %s\n", key)
		}
	}
}

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}