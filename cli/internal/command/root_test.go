package command

import "testing"

func TestRootCommandContainsModules(t *testing.T) {
	cmd := NewRoot()
	names := map[string]bool{}
	for _, child := range cmd.Commands() {
		names[child.Name()] = true
	}
	for _, name := range []string{"config", "short", "image", "clip", "file", "admin"} {
		if !names[name] {
			t.Fatalf("expected %s command", name)
		}
	}
}
