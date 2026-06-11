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
	if !cfg.Cleanup.Enabled || cfg.Cleanup.Interval == 0 {
		t.Fatal("expected cleanup job config")
	}
	got := cfg.Modules.ShortLink.DomainMappings["s.tool.sqlboy.me"]
	if got != "https://tool.sqlboy.me/short" {
		t.Fatalf("unexpected domain mapping: %q", got)
	}
}
