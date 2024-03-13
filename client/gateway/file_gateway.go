package gateway

import (
	"log"
	"os"
)

type FileGateway struct {
	filePath string
}

func NewFileGateway(filePath string) *FileGateway {
	return &FileGateway{filePath: filePath}
}

func (f *FileGateway) WriteFile(content string) error {
	file, err := os.Create(f.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	n, err := file.WriteString(content)
	if err != nil {
		return err
	}
	log.Println("Wrote ", n, " bytes")
	return nil
}
