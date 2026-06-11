package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	root string
}

func NewLocal(root string) *Local {
	return &Local{root: root}
}

func (l *Local) Put(ctx context.Context, key string, body io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path := l.path(key)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, body)
	return err
}

func (l *Local) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return os.Open(l.path(key))
}

func (l *Local) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	err := os.Remove(l.path(key))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (l *Local) path(key string) string {
	return filepath.Join(l.root, filepath.Clean(key))
}
