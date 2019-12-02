package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

	sourceFolderPath := "/data"

	fileName := fmt.Sprintf("%d.zip", time.Now().Unix())
	RecursiveZip(sourceFolderPath, fileName)

	zipString := fmt.Sprintf("%s: %s zip created.", time.Now().String(), fileName)
	fmt.Println(zipString)

	if _, err := s3Client.FPutObject(os.Getenv("BUCKET_NAME"), fileName, fileName, minio.PutObjectOptions{
		ContentType: "application/zip",
	}); err != nil {
		log.Fatalln(err)
	}
	sucessLog := fmt.Sprintf("%s: %s uploaded successfully.", time.Now().String(), fileName)
	fmt.Println(sucessLog)
}

func main() {
	c := cron.New()
	if os.Getenv("ACCESS_KEY_ID") == "" {
		log.Fatal("ACCESS_KEY_ID can't be blank")
	}
	if os.Getenv("BUCKET_NAME") == "" {
		log.Fatal("BUCKET_NAME can't be blank")
	}
	if os.Getenv("CRON_SCHEDULE") == "" {
		log.Fatal("CRON_SCHEDULE can't be blank")
	}
	if os.Getenv("S3_URL") == "" {
		log.Fatal("S3_URL can't be blank")
	}
	if os.Getenv("SECRET_ACCESS_KEY") == "" {
		log.Fatal("SECRET_ACCESS_KEY can't be blank")
	}
	c.AddFunc("* * * * *", Upload)
	c.Start()
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
