package renderer

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fmotalleb/the-one/config"
	"github.com/fmotalleb/the-one/template"
)

func RenderTemplates(t *config.Template) error {
	emptyVars := map[string]any{}
	if !t.Enabled {
		return nil
	}

	sourceDir := t.SourceDirectory
	targetDir := t.TargetDirectory
	extension := t.Extension
	overwrite := t.OverWrite

	if extension != "" && !strings.HasPrefix(extension, ".") {
		extension = "." + extension
	}

	return filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(d.Name(), extension) {
			return nil
		}

		// Compute relative path

		var relPath string

		if relPath, err = filepath.Rel(sourceDir, path); err != nil {
			return err
		}

		// Remove extension from destination file
		destFile := strings.TrimSuffix(relPath, extension)
		destPath := filepath.Join(targetDir, destFile)

		// Ensure destination directory exists
		if err = os.MkdirAll(filepath.Dir(destPath), t.GetDirMode()); err != nil {
			return err
		}

		// Skip if file exists and overwrite is false
		if _, err = os.Stat(destPath); err == nil && !overwrite {
			return nil
		}

		// Read template content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Evaluate template
		result, err := template.EvaluateTemplate(string(content), emptyVars)
		if err != nil {
			return errors.New("failed to evaluate template: " + err.Error())
		}

		// Write output
		if err = os.WriteFile(destPath, []byte(result), t.GetFileMode()); err != nil {
			return err
		}

		return nil
	})
}
