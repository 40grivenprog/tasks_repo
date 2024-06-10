package filewriter

import (
	"fmt"
	"os"
)

type FileWriter struct {
}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (fw *FileWriter) WriteResult(filePath string, result []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, value := range result {
		fmt.Fprintln(f, value) // print values to f, one per line
	}
	return nil
}
