package s3

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/honestbee/drone-helm-repo/pkg/storage"
	"github.com/honestbee/drone-helm-repo/pkg/util"
)

const testdir = "../../../_testdata"

type uploadTest struct {
	S3URI     string
	Region    string
	Bucket    string
	Prefix    string
	AccessKey string
	SecretKey string
}

var testData = uploadTest{
	S3URI:     "s3://test-bucket/test/prefix",
	Region:    "ap-southeast-1",
	Bucket:    "test-bucket",
	Prefix:    "test/prefix",
	AccessKey: "TESTKEY",
	SecretKey: "TESTSECRET",
}

func TestStoreFileWithCredentialsFromEnv(t *testing.T) {
	os.Setenv("AWS_ACCESS_KEY_ID", testData.AccessKey)
	os.Setenv("AWS_SECRET_ACCESS_KEY", testData.SecretKey)
	s3ObjectStore, err := mockS3ObjectStore(&Config{
		AccessKey: "",
		SecretKey: "",
		Region:    testData.Region,
		S3URI:     testData.S3URI,
	})
	if err != nil {
		t.Errorf("%s\n", err)
	}
	logger, buf := getTestLogger()
	file := &util.FileStat{
		Path: filepath.Join(testdir, "samplechart/one-0.1.0.tgz"),
	}

	err = s3ObjectStore.StoreFile(file, logger)
	if err != nil {
		t.Errorf("%s\n", err)
		t.Errorf("%s\n", buf)
	}
}

func TestStoreFileWithStaticCredentials(t *testing.T) {
	s3ObjectStore, err := mockS3ObjectStore(&Config{
		AccessKey: testData.AccessKey,
		SecretKey: testData.SecretKey,
		Region:    testData.Region,
		S3URI:     testData.S3URI,
	})
	if err != nil {
		t.Errorf("%s\n", err)
	}
	logger, buf := getTestLogger()
	file := &util.FileStat{
		Path: filepath.Join(testdir, "samplechart/one-0.1.0.tgz"),
	}

	err = s3ObjectStore.StoreFile(file, logger)
	if err != nil {
		t.Errorf("%s\n", err)
		t.Errorf("%s\n", buf)
	}
}

// Create mocked s3 object store which uses the internal driver
// but mocks the aws s3manager uploader
func mockS3ObjectStore(conf *Config) (storage.ObjectStore, error) {
	s3URL, err := parseS3URI(conf.S3URI)
	if err != nil {
		return nil, err
	}
	awsConfig := initAwsConfig(conf)

	s3Uploader := mockedS3Uploader{config: awsConfig}
	var d storage.ObjectStore
	d = &driver{
		config:       *conf,
		s3Uploader:   s3Uploader,
		bucket:       s3URL.Host,
		bucketPrefix: strings.TrimPrefix(s3URL.Path, "/"),
	}
	return d, nil
}

func (uploader mockedS3Uploader) UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, opts ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	credentials, err := uploader.config.Credentials.Get()
	if err != nil {
		return nil, err
	}

	if credentials.AccessKeyID != testData.AccessKey {
		return nil, fmt.Errorf("AccessKey - Expected %q - got %q", testData.AccessKey, credentials.AccessKeyID)
	}
	if credentials.SecretAccessKey != testData.SecretKey {
		return nil, fmt.Errorf("SecretKey - Expected %q - got %q", testData.SecretKey, credentials.AccessKeyID)
	}
	if *uploader.config.Region != testData.Region {
		return nil, fmt.Errorf("Region - Expected %q - got %q", testData.Region, *uploader.config.Region)
	}
	if *input.Bucket != testData.Bucket {
		return nil, fmt.Errorf("Bucket - Expected %q - got %q", testData.Bucket, *input.Bucket)
	}
	if *input.Key != testData.Prefix {
		return nil, fmt.Errorf("Prefix - Expected %q - got %q", testData.Prefix, *input.Key)
	}
	return nil, nil
}

type mockedS3Uploader struct {
	config *aws.Config
}

func (uploader mockedS3Uploader) Upload(input *s3manager.UploadInput, options ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	return uploader.UploadWithContext(aws.BackgroundContext(), input, options...)
}

func getTestLogger() (*util.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	return &util.Logger{
		Out:   log.New(buf, "[Out] ", log.Lshortfile),
		Err:   log.New(buf, "[Err] ", log.Lshortfile),
		Debug: log.New(buf, "[DEBUG] ", log.Lshortfile),
	}, buf
}
