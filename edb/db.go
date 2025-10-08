package edb

import (
	"fmt"
	"os"
	"path/filepath"
)

type Database struct {
	base string
}

func (d *Database) Use(table string) (*Table, error) {
	b := d.base
	filePath := filepath.Join(b, fmt.Sprintf("%s.db", table))
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return &Table{file: f}, nil
}

func (d *Database) Create(table string) error {
	b := d.base
	filePath := filepath.Join(b, fmt.Sprintf("%s.db", table))
	f, err := os.Create(filePath)
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()

	if err != nil {
		return err
	}

	return nil
}
