package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go /path/to/xml/files")
	}

	xmlDir := os.Args[1]

	// Find XSD in parent directory
	parentDir := filepath.Dir(xmlDir)
	var xsdPath string

	err := filepath.Walk(parentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".xsd") {
			xsdPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error searching for XSD: %v", err)
	}

	if xsdPath == "" {
		log.Fatal("❌ No .xsd file found in the parent directory.")
	}

	log.Printf("📄 Using XSD file: %s\n", xsdPath)

	// Validate all XML files
	err = filepath.Walk(xmlDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing file %s: %v\n", path, err)
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".xml") {
			log.Printf("🔍 Validating: %s\n", path)
			cmd := exec.Command("xmllint", "--noout", "--schema", xsdPath, path)
			output, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("❌ INVALID XML: %s\nError: %v\nDetails:\n%s\n", path, err, string(output))
			} else {
				log.Printf("✅ VALID XML: %s\n", path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking XML directory: %v", err)
	}
}
