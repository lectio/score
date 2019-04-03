package score

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// LinkScoresStore defines a storage engine for link scores
type LinkScoresStore interface {
	Read(targetURLUniqueKey string) (LinkScores, error)
	Write(scores LinkScores) error
}

// LinkScoresJSONFileStore defines a simple JSON file storage engine for LinkScores
type LinkScoresJSONFileStore struct {
	validScoresStoragePath   string
	invalidScoresStoragePath string
}

// MakeLinkScoresJSONFileStore creates a new JSON file store in the given paths
func MakeLinkScoresJSONFileStore(validScoresStoragePath string, invalidScoresStoragePath string, createDestPaths bool) (*LinkScoresJSONFileStore, error) {
	result := new(LinkScoresJSONFileStore)
	result.validScoresStoragePath = validScoresStoragePath
	result.invalidScoresStoragePath = invalidScoresStoragePath

	if createDestPaths {
		if _, err := createDirIfNotExist(validScoresStoragePath); err != nil {
			return result, fmt.Errorf("unable to create valid scores path %q: %v", validScoresStoragePath, err)
		}
		if _, err := createDirIfNotExist(invalidScoresStoragePath); err != nil {
			return result, fmt.Errorf("unable to create valid scores path %q: %v", invalidScoresStoragePath, err)
		}
	}

	if _, err := os.Stat(validScoresStoragePath); os.IsNotExist(err) {
		return result, fmt.Errorf("valid scores path %q does not exist: %v", validScoresStoragePath, err)
	}
	if _, err := os.Stat(invalidScoresStoragePath); os.IsNotExist(err) {
		return result, fmt.Errorf("invalid scores path %q does not exist: %v", invalidScoresStoragePath, err)
	}

	return result, nil
}

// Path returns either validScoresStoragePath or invalidScoresStoragePath based on scores.IsValid()
func (f LinkScoresJSONFileStore) Path(scores LinkScores) string {
	if scores.IsValid() {
		return f.validScoresStoragePath
	}
	return f.invalidScoresStoragePath
}

// FileName creates the name of this file for file storage
func (f LinkScoresJSONFileStore) FileName(scores LinkScores) string {
	path := f.validScoresStoragePath
	suffix := scores.Identity().MachineName()
	if !scores.IsValid() {
		path = f.invalidScoresStoragePath
		suffix = suffix + "-error"
	}
	return fmt.Sprintf("%s.%s.json", filepath.Join(path, scores.TargetURLUniqueKey()), suffix)
}

func (f LinkScoresJSONFileStore) Read(targetURLUniqueKey string) (LinkScores, error) {
	panic("not implemented yet")
}

func (f LinkScoresJSONFileStore) Write(scores LinkScores) error {
	if scores == nil {
		return errors.New("unable to create scores JSON file, scores is nil")
	}
	fileName := f.FileName(scores)
	file, createErr := os.Create(fileName)
	if createErr != nil {
		return fmt.Errorf("unable to create scores JSON file %q: %v", fileName, createErr)
	}
	defer file.Close()

	frontMatter, fmErr := json.MarshalIndent(scores, "", "	")
	if fmErr != nil {
		return fmt.Errorf("unable to marshal scores into JSON %q: %v", fileName, fmErr)
	}

	_, writeErr := file.Write(frontMatter)
	if writeErr != nil {
		return fmt.Errorf("unable to write scores JSON file %q: %v", fileName, writeErr)
	}

	return nil
}

// createDirIfNotExist creates a path if it does not exist. It is similar to mkdir -p in shell command,
// which also creates parent directory if not exists.
func createDirIfNotExist(dir string) (bool, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		return true, err
	}
	return false, nil
}
