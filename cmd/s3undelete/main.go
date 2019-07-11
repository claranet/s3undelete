package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"log"
	"os"
	"s3undelete/pkg/s3undelete"
	"strings"
	"time"
)

func main() {
	var bucket = flag.String("bucket", "", "Target S3 bucket (required)")
	var keys = flag.Int64("keys", 1024, "Maximum number of keys per request")
	var age = flag.Duration("age", time.Hour, "Maximum time since deletion")

	flag.Parse()

	if strings.TrimSpace(*bucket) == "" {
		log.Print("bucket is a required parameter")
		flag.Usage()
		os.Exit(2)
	}

	if err := s3undelete.Undelete(*bucket, *keys, *age); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Fatalf("undelete bucket %s failed: %s", *bucket, awsErr.Message())
    	} else {
			log.Fatalf("undelete bucket %s failed: %v", *bucket, err)
		}
	}
}
