package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	testDir, err := os.MkdirTemp("", "terraform-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	nestedDir := filepath.Join(testDir, "modules", "networking")
	if err := os.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create nested directories: %v", err)
	}

	testFiles := map[string]string{
		filepath.Join(testDir, "main.tf"): `
provider "aws" {
  region = "us-west-2"
}

resource "aws_instance" "web_server" {
  ami           = "ami-0c55b159cbfafe1f0"
  instance_type = "t2.micro"
  
  tags = {
    Name = "WebServer"
  }
}

import {
  to = aws_instance.web_server
  id = "i-abcd1234"
}

resource "aws_s3_bucket" "data" {
  bucket = "my-data-bucket"
}

import {
  to = aws_s3_bucket.data
  id = "my-data-bucket"
}`,

		filepath.Join(testDir, "modules", "main.tf"): `
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

module "vpc" {
  source = "./networking"
}

import {
  to = module.vpc
  id = "vpc-12345"
}`,

		filepath.Join(nestedDir, "vpc.tf"): `
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"
  
  tags = {
    Name = "MainVPC"
  }
}

import {
  to = aws_vpc.main
  id = "vpc-67890"
}

resource "aws_subnet" "public" {
  vpc_id     = aws_vpc.main.id
  cidr_block = "10.0.1.0/24"
}

import {
  to = aws_subnet.public
  id = "subnet-12345"
}`,
	}

	for path, content := range testFiles {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", path, err)
		}
	}

	stats := Stats{}
	
	files, err := findTerraformFiles(testDir)
	if err != nil {
		t.Fatalf("findTerraformFiles failed: %v", err)
	}
	
	for _, file := range files {
		err := processFile(file, &stats)
		if err != nil {
			t.Fatalf("processFile failed for %s: %v", file, err)
		}
	}

	expectedStats := Stats{
		FilesProcessed:      3,
		FilesModified:       3,
		ImportBlocksRemoved: 5,
	}
	
	if stats.FilesProcessed != expectedStats.FilesProcessed {
		t.Errorf("Expected FilesProcessed to be %d, but got %d", expectedStats.FilesProcessed, stats.FilesProcessed)
	}
	if stats.FilesModified != expectedStats.FilesModified {
		t.Errorf("Expected FilesModified to be %d, but got %d", expectedStats.FilesModified, stats.FilesModified)
	}
	if stats.ImportBlocksRemoved != expectedStats.ImportBlocksRemoved {
		t.Errorf("Expected ImportBlocksRemoved to be %d, but got %d", expectedStats.ImportBlocksRemoved, stats.ImportBlocksRemoved)
	}

	for path := range testFiles {
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read modified file %s: %v", path, err)
		}
		
		t.Logf("File %s content after processing: %s", path, string(content))
		
		if string(content) == "" {
			t.Errorf("File %s should not be empty after processing", path)
		}
		
		if string(content) == testFiles[path] {
			t.Errorf("File %s was not modified", path)
		}
		
		if strings.Contains(string(content), "import {") {
			t.Errorf("File %s still contains import blocks after processing", path)
		}
	}
}

func TestEdgeCases(t *testing.T) {
	testDir, err := os.MkdirTemp("", "terraform-edge-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(testDir)

	emptyFile := filepath.Join(testDir, "empty.tf")
	if err := os.WriteFile(emptyFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}

	onlyImportFile := filepath.Join(testDir, "only_import.tf")
	onlyImportContent := `
import {
  to = aws_instance.example
  id = "i-abcd1234"
}

import {
  to = aws_s3_bucket.example
  id = "my-bucket"
}
`
	if err := os.WriteFile(onlyImportFile, []byte(onlyImportContent), 0644); err != nil {
		t.Fatalf("Failed to write only import file: %v", err)
	}

	commentedFile := filepath.Join(testDir, "commented.tf")
	commentedContent := `
resource "aws_instance" "web" {
  ami           = "ami-123456"
  instance_type = "t2.micro"
}

# This is a commented import block
# import {
#   to = aws_instance.web
#   id = "i-abcd1234"
# }
`
	if err := os.WriteFile(commentedFile, []byte(commentedContent), 0644); err != nil {
		t.Fatalf("Failed to write commented file: %v", err)
	}

	stats := Stats{}
	
	files, err := findTerraformFiles(testDir)
	if err != nil {
		t.Fatalf("findTerraformFiles failed: %v", err)
	}
	
	for _, file := range files {
		err := processFile(file, &stats)
		if err != nil {
			t.Fatalf("processFile failed for %s: %v", file, err)
		}
	}

	if stats.FilesProcessed != 3 {
		t.Errorf("Expected FilesProcessed to be 3, but got %d", stats.FilesProcessed)
	}
	if stats.FilesModified != 1 {
		t.Errorf("Expected FilesModified to be 1, but got %d", stats.FilesModified)
	}
	if stats.ImportBlocksRemoved != 2 {
		t.Errorf("Expected ImportBlocksRemoved to be 2, but got %d", stats.ImportBlocksRemoved)
	}

	content, err := os.ReadFile(onlyImportFile)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}
	
	if len(content) > 0 && len(string(content)) > 0 {
		contentStr := string(content)
		for _, c := range contentStr {
			if c != ' ' && c != '\n' && c != '\t' && c != '\r' {
				t.Errorf("Expected only_import.tf to be empty or whitespace, but got: %s", contentStr)
				break
			}
		}
	}
}
