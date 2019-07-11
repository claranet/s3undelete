# s3undelete [![Documentation](https://godoc.org/github.com/claranet/s3undelete?status.svg)](http://godoc.org/github.com/claranet/s3undelete)

This utility simplifies the process of undeleting a deleted file in a
[versioned](https://docs.aws.amazon.com/AmazonS3/latest/dev/Versioning.html) S3 bucket. When versioning is enabled on a
bucket, deleting a file actually creates a delete marker which effectively masks the previous versions. To undelete a
file, the delete marker iteself must be deleted.

As detailed in [Deleting Object Versions](https://docs.aws.amazon.com/AmazonS3/latest/dev/DeletingObjectVersions.html),
you can use the console to see all versions, identify the delete maker and then remove it. Alternatively, you can make
a series of calls to the AWS API. These approaches are OK for a single file or possibly even a few, however when many
files need to be restored it would be too time consuming. `s3undelete` was built to perform bulk undeletes for objects
deleted within a configurable time range, by default the last hour.

### Installation

If you have [installed Go](http://golang.org/doc/install.html), you can simply run this command
to install `s3undelete`:

```bash
go get github.com/claranet/s3undelete/cmd/s3undelete
```

You can also download the [latest](https://github.com/claranet/s3undelete/releases/latest) x64 release.

### Usage

AWS access is achieved using the default credential provider chain as part of the AWS SDK. As detailed in the
[Specifying Credentials](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html) section of the
SDK documentation, credentials are sought in environment variables, the shared credentials file and finally the instance
profile if you are running within AWS. Please note that you will need to specify your region, for example with the
`AWS_REGION` environment variable.

##### Example

Assuming you have [got your AWS access keys](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/setting-up.html),
you can export the three required environment variables and call `s3undelete`. `s3undelete` requires the `-bucket`
parameter with which you name the bucket you wish to undelete files in.

```bash
export AWS_ACCESS_KEY_ID=*****
export AWS_SECRET_ACCESS_KEY=*****
export AWS_REGION=eu-west-1
s3undelete -bucket my-versioned-bucket 
```

##### Command line arguments

`s3undelete` accepts the following command line arguments:

- `-age` duration

  Maximum time since deletion, as a [duration specification](https://golang.org/pkg/time/#ParseDuration) with a default
  of an hour (`1h`).

- `-bucket` string **required**

  Target S3 bucket name.

- `-keys` int

  Maximum number of keys per request (default `1024`)

### Developing & Testing

Instead of using `go get`, you can clone this repository and use the `Makefile`. The following targets are available:

- `lint`

  Runs `golint` across the source reporting any style mistakes. If not already installed locally, you can run
  `go get -u golang.org/x/lint/golint` to install.
  
- `build`

  Runs `lint` and compiles the source to produce the `s3undelete` binary in the local directory.
  
- `test`

  Runs `build` then uses [Terraform](https://www.terraform.io/) to create two buckets with 5 objects each, one with
  versioning enabled and the other not. These objects are deleted and `s3undelete` is then tested. Once the tests have
  passed, the bucekts are destroyed. Terraform is configured in the same way as `s3undelete` but requires additional IAM
  permissions as detailed below.

- `install` **default**

  Runs `test` and copies the local `s3undelete` to the user's `$GOPATH/bin` folder.
  
- `clean`

  Removes the local `s3undelete` if present and runs `terraform destroy` to ensure the buckets have been removed.

### IAM Permissions

The following IAM policy documents detail the minimum permissions required to execute `s3undelete` and `terraform`.

##### Minimum required permissions for `s3undelete`

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:DeleteObjectVersion"
      ],
      "Resource": [
        "arn:aws:s3:::YOUR-BUCKET-NAME/*"
      ],
      "Effect": "Allow"
    },
    {
      "Action": [
        "s3:ListBucketVersions"
      ],
      "Resource": [
        "arn:aws:s3:::YOUR-BUCKET-NAME"
      ],
      "Effect": "Allow"
    }
  ]
}
```

##### Minimum required permissions for `terraform`

```
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "s3:CreateBucket",
        "s3:DeleteBucket",
        "s3:DeleteObject",
        "s3:DeleteObjectVersion",
        "s3:Get*",
        "s3:ListBucket",
        "s3:ListBucketVersions",
        "s3:PutBucketVersioning",
        "s3:PutObject"
      ],
      "Resource": [
        "*"
      ],
      "Effect": "Allow"
    }
  ]
}
``` 