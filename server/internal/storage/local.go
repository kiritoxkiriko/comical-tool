package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Local struct {
	root string
}

func NewLocal(root string) *Local {
	return &Local{root: root}
}

func (l *Local) Ensure(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return os.MkdirAll(l.root, 0o755)
}

func (l *Local) Put(ctx context.Context, key string, body io.Reader) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := l.path(key)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, body); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func (l *Local) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	path, err := l.path(key)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (l *Local) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	path, err := l.path(key)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (l *Local) Head(ctx context.Context, key string) (ObjectInfo, error) {
	if err := ctx.Err(); err != nil {
		return ObjectInfo{}, err
	}
	path, err := l.path(key)
	if err != nil {
		return ObjectInfo{}, err
	}
	info, err := os.Stat(path)
	if err != nil {
		return ObjectInfo{}, err
	}
	return ObjectInfo{Size: info.Size()}, nil
}

func (l *Local) path(key string) (string, error) {
	clean := filepath.Clean(key)
	if clean == "." || filepath.IsAbs(clean) || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || clean == ".." {
		return "", errors.New("invalid object key")
	}
	return filepath.Join(l.root, clean), nil
}
