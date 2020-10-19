package marshal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kick-project/kick/internal/utils/file"
	"gopkg.in/yaml.v2"
)

// MarshalFile marshals to a json or yaml.
func MarshalFile(v interface{}, path string) error {
	d, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("can not get absolute path: %w", err)
	}
	if _, err := os.Stat(d); os.IsNotExist(err) {
		return fmt.Errorf("parent directory of %s does not exists: %w", path, err)
	}

	f := file.NewAtomicWrite(path)
	defer f.Close()
	var out []byte
	st, err := suffixType(path)
	if err != nil {
		return err
	}
	if st == "json" {
		out, err = json.Marshal(v)
	} else if st == "yaml" {
		out, err = yaml.Marshal(v)
	}
	if err != nil {
		return fmt.Errorf("can not marshal: %w", err)
	}
	_, err = f.Write(out)
	if err != nil {
		return fmt.Errorf("can not write to file: %w", err)
	}
	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}
	return nil
}

// UnmarshalFile un-marshals from a json or yaml file.
func UnmarshalFile(v interface{}, p string) error {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("file %s does not exist: %w", p, err)
	}

	f, err := ioutil.ReadFile(p)
	if err != nil {
		return fmt.Errorf("can not read file %s: %w", p, err)
	}

	st, err := suffixType(p)
	if err != nil {
		return err
	}
	if st == "json" {
		err = json.Unmarshal([]byte(f), v)
	} else if st == "yaml" || st == "yml" {
		err = yaml.Unmarshal([]byte(f), v)
	}
	if err != nil {
		return fmt.Errorf("can not unmarshal file %s: %w", p, err)
	}
	return nil
}

// suffixType returns json or yaml based on file suffixType.
func suffixType(path string) (string, error) {
	if strings.HasSuffix(path, ".json") {
		return "json", nil
	}
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return "yaml", nil
	}

	err := fmt.Errorf(`file %s does not have a suffix of "*.json", "*.yaml" or "*.yml"`, path)
	return "", err
}
