package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go"
	"github.com/robfig/cron"
)

// RecursiveZip is function to create zip for given directory
func RecursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

// Upload is function to upload zip to Spaces
func Upload() {
	s3Client, err := minio.New(os.Getenv("S3_URL"), os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), true)
	if err != nil {
		log.Fatalln(err)
	}

	RecursiveZip("data", "./data.zip")

	if _, err := s3Client.FPutObject(os.Getenv("BUCKET_NAME"), "data.zip", "data.zip", minio.PutObjectOptions{
		ContentType: "application/zip",
	}); err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully uploaded")
}

func main() {
	c := cron.New()
	c.AddFunc(os.Getenv("CRON_SCHEDULE"), Upload)
	c.Start()
}
