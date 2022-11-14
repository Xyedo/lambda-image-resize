package s3util

import (
	appcnf "aws-lambda/resize"
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func New(conf appcnf.Config) *S3utils {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(conf.Region))
	if err != nil {
		log.Printf("s3manager::New: cant initiate aws-config:%s\n", err.Error())
		os.Exit(1)
		return nil
	}

	client := s3.NewFromConfig(cfg)
	return &S3utils{
		client:     client,
		downloader: manager.NewDownloader(client),
		uploader:   manager.NewUploader(client),
		config:     conf,
	}
}

type S3utils struct {
	client     *s3.Client
	downloader *manager.Downloader
	uploader   *manager.Uploader
	config     appcnf.Config
}

func (man *S3utils) DownloadFile(objectKey string) (string, error) {
	fileName := filepath.Base(objectKey)

	file, err := os.Create(man.config.LocalDir + "/" + fileName)
	if err != nil {
		log.Fatalf("s3manager::DownloadFile: cant create osDir:%s", err.Error())
		return "", err
	}
	defer file.Close()

	nB, err := man.downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(man.config.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Fatalf("s3manager::DownloadFile cant download file:%s", err.Error())
		return "", err
	}
	log.Println("s3manager::DownloadFile: Downloaded", file.Name(), nB, "bytes")
	return man.config.LocalDir + "/" + fileName, nil

}

func (man *S3utils) DeleteFile(objectKey string) error {
	_, err := man.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(man.config.BucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}
	return nil
}

func (man *S3utils) UploadFile(filename string) error {
	filePath := man.config.LocalDir + "/" + filename
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("s3manager::uploadFile: Unable to open file, %v", err)
		return err
	}
	defer file.Close()
	key := man.config.UploadDir + "/" + filename
	_, err = man.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(man.config.BucketName),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("s3manager::uploadFile: Unable to upload the file, %v", err)
		return err
	}
	log.Printf("s3manager::uploadFile: Successfully uploaded %q to %q\n", filename, man.config.BucketName)
	return nil
}
