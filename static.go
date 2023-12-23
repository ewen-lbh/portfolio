package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func StaticallyRender(inside string, origin string, path string) error {
	outputPath := filepath.Join(inside, path, "index.html")
	url := fmt.Sprintf("%s/%s", origin, path)

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("while downloading %s: %w", url, err)
	}

	os.MkdirAll(filepath.Dir(outputPath), 0755)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("while creating %s: %w", outputPath, err)
	}

	defer file.Close()
	defer response.Body.Close()

	io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("while writing %s: %w", outputPath, err)
	}

	// fmt.Printf("[  ] Rendered %s to %s\n", url, outputPath)
	return nil
}
