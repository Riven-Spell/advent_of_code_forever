package inputs

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Riven-Spell/advent_of_code_forever/core"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type InputCache struct {
	cacheDir string
}

type Solution struct {
	A, B any
}

func (s *Solution) Empty() bool {
	return s == nil || (s.A == nil && s.B == nil)
}

var Cache = &InputCache{}

func (i *InputCache) GetCacheDir() (string, error) {
	if i.cacheDir != "" {
		return i.cacheDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	i.cacheDir = filepath.Join(home, ".aocf/")

	err = os.MkdirAll(i.cacheDir, 0755)

	return i.cacheDir, nil
}

func (i *InputCache) HasCachedInput(day, year uint) bool {
	cDir, err := i.GetCacheDir()
	if err != nil {
		return false
	}

	inputPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.txt", year, day))

	_, err = os.Stat(inputPath)
	return err == nil
}

func (i *InputCache) DeleteInput(day, year uint) error {
	cDir, err := i.GetCacheDir()
	if err != nil {
		return err
	}

	if year < 2015 {
		return errors.New("cannot delete non-existent AoC year (< 2015)")
	}

	inputPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.txt", year, day))
	solutionPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.solution.txt", year, day))

	if day > 0 && day <= 25 {
		_ = os.Remove(solutionPath) // it's OK if this fails because maybe it doesn't exist.

		return os.Remove(inputPath)
	} else {
		return os.RemoveAll(filepath.Dir(solutionPath))
	}
}

func (i *InputCache) PutSolution(day, year uint, solution Solution, replace bool) error {
	cDir, err := i.GetCacheDir()
	if err != nil {
		return err
	}

	solutionPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.solution.txt", year, day))
	err = os.MkdirAll(filepath.Dir(solutionPath), 0755)
	if err != nil {
		return err
	}

	if !replace {
		_, err = os.Stat(solutionPath)
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot put input: file either exists, or stat failed: %w", err)
		}
	}

	if solution.Empty() {
		err = os.Remove(solutionPath)
		return err
	}

	f, err := os.OpenFile(solutionPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	buf, _ := json.Marshal(solution)

	_, err = io.Copy(f, strings.NewReader(string(buf)))
	if err != nil {
		return err
	}

	return f.Close()
}

func (i *InputCache) GetSolution(day, year uint) (*Solution, error) {
	cDir, err := i.GetCacheDir()
	if err != nil {
		return nil, err
	}

	solutionPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.solution.txt", year, day))
	err = os.MkdirAll(filepath.Dir(solutionPath), 0755)
	if err != nil {
		return nil, err
	}

	f, err := os.OpenFile(solutionPath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var out Solution
	err = json.Unmarshal(buf, &out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func (i *InputCache) PutInput(day, year uint, input io.Reader, replace bool) error {
	cDir, err := i.GetCacheDir()
	if err != nil {
		return err
	}

	inputPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.txt", year, day))
	err = os.MkdirAll(filepath.Dir(inputPath), 0755)
	if err != nil {
		return err
	}

	if !replace {
		_, err = os.Stat(inputPath)
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot put input: file either exists, or stat failed: %w", err)
		}
	}

	f, err := os.OpenFile(inputPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, input)
	if err != nil {
		return err
	}

	return f.Close()
}

func (i *InputCache) DownloadInput(day, year uint, replace bool) error {
	targetURL := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day)

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return err
	}

	token, def := core.EEnvironmentVariable.AuthToken().Get()
	if def {
		return errors.New("auth token not specified, cannot download input.")
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: token,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return i.PutInput(day, year, resp.Body, replace)
}

func (i *InputCache) GetInput(day, year uint) (string, error) {
	cDir, err := i.GetCacheDir()
	if err != nil {
		return "", err
	}

	solutionPath := filepath.Join(cDir, fmt.Sprintf("%d/%d.txt", year, day))
	err = os.MkdirAll(filepath.Dir(solutionPath), 0755)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(solutionPath, os.O_RDONLY, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (i *InputCache) GetInputAndSolution(day, year uint) (string, *Solution, error) {
	input, err := i.GetInput(day, year)
	if err != nil {
		return "", nil, err
	}

	solution, _ := i.GetSolution(day, year)

	return input, solution, nil
}
