package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"s3undelete/pkg/s3undelete"
)

func main() {
	var bucket = flag.String("bucket", "", "Target S3 bucket (required)")
	var keys = flag.Int64("keys", 1024, "Maximum number of keys per request")
	var age = flag.Duration("age", time.Hour, "Maximum time since deletion")
	var execute = flag.Bool("execute", false, "If true, executes the undelete operation. Otherwise, a dry run is performed.")
	var prefix = flag.String("prefix", "", "S3 Key Prefix")
	var filter = flag.String("filter", "", "A regex expression to filter the files by")

	flag.Parse()

	if strings.TrimSpace(*bucket) == "" {
		log.Print("bucket is a required parameter")
		flag.Usage()
		os.Exit(2)
	}

	if err := s3undelete.Undelete(*bucket, *keys, *age, *execute, *prefix, *filter); err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			log.Fatalf("undelete bucket %s failed: %s", *bucket, awsErr.Message())
		} else {
			log.Fatalf("undelete bucket %s failed: %v", *bucket, err)
		}
	}
}
