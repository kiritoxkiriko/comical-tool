package config

import "testing"

func TestLoadExampleConfig(t *testing.T) {
	cfg, err := Load("../../../deploy/config.example.toml")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Modules.ShortLink.DefaultTTL == 0 {
		t.Fatal("expected short link default ttl")
	}
	got := cfg.Modules.ShortLink.DomainMappings["s.example.com"]
	if got != "https://myapp.example.com/short" {
		t.Fatalf("unexpected domain mapping: %q", got)
	}
}
