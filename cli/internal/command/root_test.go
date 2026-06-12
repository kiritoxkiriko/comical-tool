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

func TestFileCommandContainsDownload(t *testing.T) {
	cmd := NewRoot()
	fileCmd, _, err := cmd.Find([]string{"file"})
	if err != nil {
		t.Fatal(err)
	}
	names := map[string]bool{}
	for _, child := range fileCmd.Commands() {
		names[child.Name()] = true
	}
	for _, name := range []string{"upload", "list", "download", "delete"} {
		if !names[name] {
			t.Fatalf("expected file %s command", name)
		}
	}
}

func TestFileCommandAccessPolicyFlags(t *testing.T) {
	cmd := NewRoot()
	uploadCmd, _, err := cmd.Find([]string{"file", "upload"})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"password", "max-visits"} {
		if uploadCmd.Flags().Lookup(name) == nil {
			t.Fatalf("expected file upload --%s flag", name)
		}
	}
	downloadCmd, _, err := cmd.Find([]string{"file", "download"})
	if err != nil {
		t.Fatal(err)
	}
	if downloadCmd.Flags().Lookup("password") == nil {
		t.Fatal("expected file download --password flag")
	}
}
