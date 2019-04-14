package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/lectio/score"
)

const defaultFilePerm os.FileMode = 0644

// defaultValidScoresFileNameFormat creates filename of the format {path}/{uniqueId}_{machineName}.json
const defaultValidScoresFileNameFormat = "%[2]s_%[1]s.json"

// defaultInvalidScoresFileNameFormat creates filename of the format {path}/{uniqueId}.{machineName}-errors.json
const defaultInvalidScoresFileNameFormat = "%[2]s.%[1]s-error.json"

type fileCache struct {
	perm                        os.FileMode
	keys                        score.Keys
	validScoresFileNameFormat   string
	invalidScoresFileNameFormat string
	validScoresPath             string
	invalidScoresPath           string
	initialTotalCount           int
	simulate                    bool
	scorerMachineName           string
}

// MakeFileCache creates an instance of a cache, which stores links on disk, in a named path
func MakeFileCache(validScoresPath string, invalidScoresPath string, createPaths bool, keys score.Keys, initialTotalCount int, simulate bool) (Cache, error) {
	if createPaths {
		if err := os.MkdirAll(validScoresPath, defaultFilePerm); err != nil {
			return nil, err
		}
		if err := os.MkdirAll(invalidScoresPath, defaultFilePerm); err != nil {
			return nil, err
		}
	}

	if _, err := os.Stat(validScoresPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("lectio score fileCache validScoresPath does not exist: %q", validScoresPath)
	}

	if _, err := os.Stat(invalidScoresPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("lectio score fileCache invalidScoresPath does not exist: %q", invalidScoresPath)
	}

	cache := new(fileCache)
	cache.perm = defaultFilePerm
	cache.validScoresPath = validScoresPath
	cache.invalidScoresPath = invalidScoresPath
	cache.validScoresFileNameFormat = defaultValidScoresFileNameFormat
	cache.invalidScoresFileNameFormat = defaultInvalidScoresFileNameFormat
	cache.keys = keys
	cache.initialTotalCount = initialTotalCount
	cache.simulate = simulate
	cache.scorerMachineName = "aggregate"
	return cache, nil
}

func (c fileCache) Score(url *url.URL) (score.LinkScores, error) {
	als := score.GetAggregatedLinkScores(url, c.keys, c.initialTotalCount, c.simulate)
	return als, nil
}

func (c fileCache) Get(url *url.URL) (score.LinkScores, error) {
	link, found, expired, err := c.Find(url)
	if err != nil {
		return nil, err
	}

	if found && !expired {
		return link, err
	}

	link, err = c.Score(url)
	if err != nil {
		return nil, err
	}
	c.Save(link, 0)
	return link, nil
}

func (c fileCache) find(url *url.URL, fileName string) (score.LinkScores, bool, bool, error) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return nil, false, true, nil
	}

	file, openErr := os.Open(fileName)
	if openErr != nil {
		return nil, false, true, openErr
	}

	bytes, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return nil, false, true, readErr
	}

	var als score.AggregatedLinkScores
	parseErr := json.Unmarshal(bytes, &als)
	if parseErr != nil {
		return nil, false, true, parseErr
	}

	return &als, true, false, nil
}

func (c fileCache) Find(url *url.URL) (score.LinkScores, bool, bool, error) {
	// first look for the "valid" version; if found, use it
	fileName := c.urlFileName(url, true)
	scores, found, expired, err := c.find(url, fileName)
	if found && !expired && err == nil {
		return scores, found, expired, err
	}

	// looks like a "valid" version was not found, see if we have an invalid version
	fileName = c.urlFileName(url, false)
	return c.find(url, fileName)
}

func (c fileCache) Save(scores score.LinkScores, autoExpire time.Duration) error {
	if scores == nil {
		return errors.New("unable to create scores JSON file, scores is nil")
	}
	fileName := c.scoresFileName(scores)
	file, createErr := os.Create(fileName)
	if createErr != nil {
		return fmt.Errorf("unable to create scores JSON file %q: %v", fileName, createErr)
	}
	defer file.Close()

	scoresFile, fmErr := json.MarshalIndent(scores, "", "	")
	if fmErr != nil {
		return fmt.Errorf("unable to marshal scores into JSON %q: %v", fileName, fmErr)
	}

	_, writeErr := file.Write(scoresFile)
	if writeErr != nil {
		return fmt.Errorf("unable to write scores JSON file %q: %v", fileName, writeErr)
	}

	return nil
}

func (c fileCache) Close() error {
	return nil
}

// Path returns either validScoresStoragePath or invalidScoresStoragePath based on scores.IsValid()
func (c fileCache) path(scores score.LinkScores) string {
	if scores.IsValid() {
		return c.validScoresPath
	}
	return c.invalidScoresPath
}

// scoresfileName creates the name of this file for file storage
func (c fileCache) scoresFileName(scores score.LinkScores) string {
	path := c.validScoresPath
	format := c.validScoresFileNameFormat
	if !scores.IsValid() {
		path = c.invalidScoresPath
		format = c.invalidScoresFileNameFormat
	}
	return filepath.Join(path, fmt.Sprintf(format, scores.Scorer().ScorerMachineName(), scores.TargetURLUniqueKey()))
}

// urlfileName creates the name of this file for file storage
func (c fileCache) urlFileName(url *url.URL, valid bool) string {
	key := c.keys.ScoreKeyForURL(url)
	path := c.validScoresPath
	format := c.validScoresFileNameFormat
	if !valid {
		path = c.invalidScoresPath
		format = c.invalidScoresFileNameFormat
	}
	return filepath.Join(path, fmt.Sprintf(format, c.scorerMachineName, key))
}
