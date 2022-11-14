package main

import (
	"context"

	"aws-lambda/resize/pkg/resizer"
	"aws-lambda/resize/pkg/s3util"

	appcnf "aws-lambda/resize"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var s3man *s3util.S3utils
var rez *resizer.Resizer

func init() {
	appconf := appcnf.New()
	s3man = s3util.New(appconf)
	rez = resizer.New(appconf.LocalDir, []int{640, 160})
}
func HandleRequest(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		objKey := record.S3.Object.Key
		srcImgPath, err := s3man.DownloadFile(objKey)
		if err != nil {
			return err
		}
		resizedImgPaths, err := rez.GetResizedImagesPath(srcImgPath)
		if err != nil {
			return err
		}
		for _, resizedImgPath := range resizedImgPaths {
			err = s3man.UploadFile(resizedImgPath)
			if err != nil {
				return err
			}
		}
		err = s3man.DeleteFile(objKey)
		if err != nil {
			return err
		}
	}
	err := rez.RemoveImages()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
