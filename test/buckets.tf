terraform {
  required_version = ">=0.11.0"
}

provider "aws" {
  version = "~> 2.0"
}

resource "aws_s3_bucket" "unversioned" {
  bucket_prefix = "unver"
  force_destroy = true
}

resource "aws_s3_bucket" "versioned" {
  bucket_prefix = "ver"
  force_destroy = true

  versioning {
    enabled = true
  }
}

resource "aws_s3_bucket_object" "unverobjs" {
  count   = 5
  bucket  = "${aws_s3_bucket.unversioned.id}"
  key     = "somefile0${count.index}"
  content = "Test file ${count.index}"
}

resource "aws_s3_bucket_object" "verobjs" {
  count   = 5
  bucket  = "${aws_s3_bucket.versioned.id}"
  key     = "somefile0${count.index}"
  content = "Test file ${count.index}"
}
