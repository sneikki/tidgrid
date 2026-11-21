package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestStarterLoadsAndRenders(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "src"), 0o755) // project.scan must exist
	p := filepath.Join(dir, "tidgrid.toml")
	if err := os.WriteFile(p, starter, 0o644); err != nil {
		t.Fatal(err)
	}
	v, errs := load(p)
	if len(errs) > 0 {
		t.Fatalf("starter does not load: %v", errs)
	}
	if out := v.render(); out == "" {
		t.Fatal("starter renders empty")
	}
}

func TestInitWritesStarter(t *testing.T) {
	p := filepath.Join(t.TempDir(), "tidgrid.toml")
	if code := runInit([]string{"-c", p}); code != 0 {
		t.Fatalf("exit %d, want 0", code)
	}
	got, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, starter) {
		t.Fatal("written file differs from starter")
	}
}

func TestInitRefusesOverwrite(t *testing.T) {
	p := filepath.Join(t.TempDir(), "tidgrid.toml")
	keep := []byte("# mine\n")
	if err := os.WriteFile(p, keep, 0o644); err != nil {
		t.Fatal(err)
	}
	if code := runInit([]string{"-c", p}); code != 2 {
		t.Fatalf("exit %d, want 2", code)
	}
	got, _ := os.ReadFile(p)
	if !bytes.Equal(got, keep) {
		t.Fatal("existing file was modified")
	}
}
