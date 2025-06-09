package lib

import (
	"bytes"
	"context"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var s3Client *s3.Client

type S3File struct {
	Key          string
	LastModified time.Time
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config: %v", err)
	}
	s3Client = s3.NewFromConfig(cfg)
}

func ListJSONFilesFromS3(bucket string, prefix string) ([]string, error) {
	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: &prefix,
	})
	if err != nil {
		return nil, err
	}

	files := []string{}
	for _, obj := range resp.Contents {
		if obj.Key != nil && *obj.Key != "" && (*obj.Key)[len(*obj.Key)-5:] == ".json" {
			files = append(files, *obj.Key)
		}
	}
	return files, nil
}

func ReadJSONFileFromS3(bucket, key string) ([]byte, error) {
	resp, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	return buf.Bytes(), err
}

func SaveToS3(bucket, key string, data []byte) error {
	reader := io.NopCloser(bytes.NewReader(data))
	contentLength := int64(len(data))
	_, err := s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(key),
		Body:          reader,
		ContentLength: aws.Int64(contentLength),
	})
	if err != nil {
		log.Printf("failed to upload to S3: %v\n", err)
	} else {
		log.Println("☁️  Uploaded to S3: ", key)
	}
	return err
}

func ListRecentJSONFilesFromS3(bucket string, within time.Duration, prefix string) ([]S3File, error) {
	resp, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: &prefix,
	})
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().Add(-within)
	var files []S3File

	for _, obj := range resp.Contents {
		if obj.Key != nil && obj.LastModified != nil &&
			obj.LastModified.After(cutoff) &&
			len(*obj.Key) > 5 && (*obj.Key)[len(*obj.Key)-5:] == ".json" {
			files = append(files, S3File{Key: *obj.Key, LastModified: *obj.LastModified})
		}
	}

	return files, nil
}
