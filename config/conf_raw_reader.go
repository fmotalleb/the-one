package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func deepMerge(dst, src map[string]any) (map[string]any, error) {
	err := mergo.Merge(&dst, src, mergo.WithAppendSlice)
	return dst, err
}

func ReadAndMergeConfig(pattern string) (map[string]any, error) {
	log().Info("Starting config load", zap.String("pattern", pattern))
	return mergeFromPattern(pattern, map[string]bool{})
}

func mergeFromPattern(pattern string, visited map[string]bool) (map[string]any, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		log().Error("Invalid glob pattern", zap.String("pattern", pattern), zap.Error(err))
		return nil, fmt.Errorf("invalid glob pattern: %w", err)
	}
	if len(files) == 0 {
		log().Warn("No config files matched pattern", zap.String("pattern", pattern))
	}

	result := make(map[string]any)
	for _, file := range files {
		log().Debug("Merging config file", zap.String("file", file))
		conf, err := readAndResolveIncludes(file, visited)
		if err != nil {
			log().Error("Failed to read and merge includes", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		result, err = deepMerge(result, conf)
		if err != nil {
			log().Error("Deep merge failed", zap.String("file", file), zap.Error(err))
			return nil, err
		}
		log().Debug("Merged successfully", zap.String("file", file))
	}
	return result, nil
}

func readAndResolveIncludes(path string, visited map[string]bool) (map[string]any, error) {
	absPath := path
	if u, err := url.Parse(path); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		absPath = u.String()
	} else if p, err := filepath.Abs(path); err == nil {
		absPath = p
	}

	if visited[absPath] {
		log().Warn("Circular include detected", zap.String("path", absPath))
		return make(map[string]any), nil
	}
	log().Info("Reading config", zap.String("path", absPath))
	visited[absPath] = true

	reader, ext, err := readFrom(path)
	if err != nil {
		return nil, err
	}

	raw, err := parseConfig(ext, reader, path)
	if err != nil {
		return nil, err
	}
	log().Debug("Parsed config", zap.String("path", path))

	raw, err = readIncludedFiles(raw, path, visited)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func readIncludedFiles(raw map[string]any, path string, visited map[string]bool) (map[string]any, error) {
	if includes, ok := raw["include"].([]any); ok {
		for _, inc := range includes {
			if pattern, ok := inc.(string); ok {
				log().Info("Processing include", zap.String("from", path), zap.String("pattern", pattern))
				included, err := mergeFromPattern(pattern, visited)
				if err != nil {
					log().Error("Failed to process include", zap.String("pattern", pattern), zap.Error(err))
					return nil, err
				}
				raw, err = deepMerge(included, raw)
				if err != nil {
					log().Error("Deep merge failed during include", zap.String("pattern", pattern), zap.Error(err))
					return nil, err
				}
				log().Debug("Include merged", zap.String("pattern", pattern))
			}
		}
	}
	return raw, nil
}

func parseConfig(ext string, reader io.Reader, path string) (map[string]any, error) {
	v := viper.New()
	v.SetConfigType(ext)
	if err := v.ReadConfig(reader); err != nil {
		log().Error("Failed to read config", zap.String("path", path), zap.Error(err))
		return nil, err
	}

	raw := make(map[string]any)
	if err := v.Unmarshal(&raw); err != nil {
		log().Error("Failed to unmarshal config", zap.String("path", path), zap.Error(err))
		return nil, err
	}
	return raw, nil
}

func readFrom(path string) (io.Reader, string, error) {
	var reader io.Reader
	var ext string
	var err error

	u, err := url.Parse(path)
	if err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		reader, ext, err = readRemote(path, u)
	} else {
		reader, ext, err = readFile(path)
	}
	if err != nil {
		return nil, "", err
	}
	return reader, ext, nil
}

func readFile(path string) (io.ReadWriter, string, error) {
	file, err := os.Open(path)
	if err != nil {
		log().Error("Failed to open file", zap.String("path", path), zap.Error(err))
		return nil, "", err
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, file); err != nil {
		log().Error("Failed to read file into buffer", zap.String("path", path), zap.Error(err))
		return nil, "", err
	}

	ext := strings.TrimPrefix(filepath.Ext(path), ".")
	return buf, ext, nil
}

func readRemote(path string, u *url.URL) (io.Reader, string, error) {
	log().Info("Fetching remote config", zap.String("url", path))

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		path,
		nil,
	)
	if err != nil {
		log().Error("Failed to create request", zap.String("url", path), zap.Error(err))
		return nil, "", err
	}
	if u.User != nil {
		pass, _ := u.User.Password()
		req.SetBasicAuth(u.User.Username(), pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log().Error("HTTP request failed", zap.String("url", path), zap.Error(err))
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log().Error("Non-200 response", zap.String("url", path), zap.Int("status", resp.StatusCode))
		return nil, "", fmt.Errorf("http error: %s", resp.Status)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		log().Error("Failed to read response body", zap.String("url", path), zap.Error(err))
		return nil, "", err
	}

	ext := strings.TrimPrefix(filepath.Ext(u.Path), ".")
	return buf, ext, nil
}
