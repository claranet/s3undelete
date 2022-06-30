//Package s3undelete provides a single function for undeleting files in S3
package s3undelete

import (
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

//Undelete looks for delete markers in a S3 bucket that are no older
//than maxage. Paged calls to ListObjectVersions are made requesting
//maxkeys per page. Matching delete markers are deleted which causes
//the most recent version of the corresponding key to be "undeleted".
//This has no effect on a bucket without versioning and returns error
func Undelete(bucket string, maxkeys int64, maxage time.Duration, execute bool, prefix, filter string) error {

	watermark := time.Now().Add(-maxage)
	s3session := s3.New(session.New(), aws.NewConfig())

	var r *regexp.Regexp
	if filter != "" {
		r = regexp.MustCompile(filter)
	}
	var keyPrefix *string
	if prefix != "" {
		keyPrefix = &prefix
	}

	listObjVers := &s3.ListObjectVersionsInput{
		Bucket:  aws.String(bucket),
		MaxKeys: aws.Int64(maxkeys),
		Prefix:  keyPrefix,
	}

	keys := int64(0)
	err := s3session.ListObjectVersionsPages(listObjVers,
		func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
			keys += maxkeys
			if keys%(maxkeys*10) == 0 {
				log.Printf("PROCESSED: %v object versions", keys)
			}

			for _, entry := range page.DeleteMarkers {
				if aws.BoolValue(entry.IsLatest) {
					if aws.TimeValue(entry.LastModified).After(watermark) && (r == nil || r.MatchString(*entry.Key)) {
						delObj := &s3.DeleteObjectInput{
							Bucket:    aws.String(bucket),
							Key:       entry.Key,
							VersionId: entry.VersionId,
						}

						if execute {
							_, err := s3session.DeleteObject(delObj)
							if err != nil {
								log.Printf("delete delete marker for key '%s': %v", aws.StringValue(entry.Key), err)
							} else {
								log.Printf(" RESTORED: %s (%s)", aws.StringValue(entry.Key), aws.StringValue(entry.VersionId))
							}
						} else {
							log.Printf("DRY RUN - Going to restore  %s (%s)", aws.StringValue(entry.Key), aws.StringValue(entry.VersionId))
						}
					} else {
						log.Printf("  IGNORED: %s (%s)", aws.StringValue(entry.Key), aws.StringValue(entry.VersionId))
					}
				}
			}

			return aws.BoolValue(page.IsTruncated)
		})

	return err
}
