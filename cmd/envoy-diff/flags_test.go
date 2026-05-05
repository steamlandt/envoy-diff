package main

import (
	"strings"
	"testing"
)

func TestParseFlags_MissingArgs(t *testing.T) {
	_, err := parseFlags([]string{})
	if err == nil {
		t.Fatal("expected error for missing arguments")
	}
}

func TestParseFlags_TooFewArgs(t *testing.T) {
	_, err := parseFlags([]string{"file.env"})
	if err == nil {
		t.Fatal("expected error for only one argument")
	}
}

func TestParseFlags_BasicSources(t *testing.T) {
	cfg, err := parseFlags([]string{"a.env", "b.env"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.sourceA != "a.env" {
		t.Errorf("sourceA = %q, want %q", cfg.sourceA, "a.env")
	}
	if cfg.sourceB != "b.env" {
		t.Errorf("sourceB = %q, want %q", cfg.sourceB, "b.env")
	}
}

func TestParseFlags_Flags(t *testing.T) {
	cfg, err := parseFlags([]string{"-keys", "-exit-code", "-color=false", "self", "pid:42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.onlyKeys {
		t.Error("expected onlyKeys to be true")
	}
	if !cfg.exitCode {
		t.Error("expected exitCode to be true")
	}
	if cfg.color {
		t.Error("expected color to be false")
	}
	if cfg.sourceA != "self" {
		t.Errorf("sourceA = %q, want %q", cfg.sourceA, "self")
	}
	if cfg.sourceB != "pid:42" {
		t.Errorf("sourceB = %q, want %q", cfg.sourceB, "pid:42")
	}
}

func TestParseFlags_UnknownFlag(t *testing.T) {
	_, err := parseFlags([]string{"-unknown", "a.env", "b.env"})
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
	if !strings.Contains(err.Error(), "unknown") && !strings.Contains(err.Error(), "flag") {
		t.Errorf("unexpected error message: %v", err)
	}
}
