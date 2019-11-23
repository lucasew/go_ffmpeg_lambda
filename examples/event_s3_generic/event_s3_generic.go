package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/lucasew/go_ffmpeg_lambda/ffmpeg"
	"os"
)

type Request struct {
	Records           []events.S3EventRecord `json:"Records"`
	Params            []string               `json:"params"`
	Destination       string                 `json:"destination"`
	DestinationBucket string                 `json:"destination_bucket"`
}

func HandleRequest(ctx context.Context, ev Request) (string, error) {
	sess := session.Must(session.NewSession())
	downloader := s3manager.NewDownloader(sess)
	f, err := os.Open("dummy.dat")
	if err != nil {
		return err.Error(), err
	}
	size, err := downloader.Download(f, &s3.GetObjectInput{
		Bucket: &ev.Records[0].S3.Bucket.Name,
		Key:    &ev.Records[0].S3.Object.Key,
	})
	if err != nil {
		return fmt.Sprintf("Cant write %d bytes to disk\n", size), err
	}
	ffsess := ffmpeg.FFMpegSession{
		From:   f,
		To:     ev.Destination,
		Params: ev.Params,
	}
	err = ffsess.Run()
	if err != nil {
		return "erro ffmpeg", err
	}
	uploader := s3manager.NewUploader(sess)
	out, err := os.Open(ev.Destination)
	if err != nil {
		return "erro ao abrir arquivo de saida", err
	}
	_, err = uploader.Upload(&s3manager.UploadInput{
		Body:   out,
		Bucket: &ev.DestinationBucket,
		Key:    &ev.Destination,
	})
	if err != nil {
		return "erro upload", err
	}
	return "", nil
}

func main() {
	lambda.Start(HandleRequest)
}
