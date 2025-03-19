package main

import (
	"bytes"
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
		t.Errorf("Have %s want %s", pathKey.Filename, expectedPathName)
	}
}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "graduationday"
	data := []byte("some jpeg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "graduationday"
	data := []byte("some jpeg bytes")

	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if ok := s.Has(key); !ok {
		t.Errorf("Expected to have key %s", key)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := ioutil.ReadAll(r)
	
	if string(b)!= string(data) {
		t.Errorf("Want %s have %s", data, b)
	}

	s.Delete(key)
}