package s3

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/honestbee/drone-helm-repo/pkg/storage"
	"github.com/honestbee/drone-helm-repo/pkg/util"
)

type (
	//Config represents an s3 objectStore config
	Config struct {
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
		Region    string `json:"region"`
		S3URI     string `json:"s3_uri"`
	}
	//internal driver
	driver struct {
		config Config
		//s3Client s3iface.S3API
		s3Uploader s3manageriface.UploaderAPI
		//s3Service    *s3.S3
		bucket       string
		bucketPrefix string
	}
)

func (d *driver) StoreFile(fileStat *util.FileStat, logger *util.Logger) error {
	logger.Debug.Printf("Uploading %s to s3://%s/%s\n", fileStat.Path, d.bucket, d.bucketPrefix)

	file, err := os.Open(fileStat.Path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			logger.Err.Printf("Problem closing file %s: %v", fileStat.Path, err)
		}
	}()

	key := filepath.Join(d.bucketPrefix, fileStat.Name)
	params := &s3manager.UploadInput{
		Bucket: aws.String(d.bucket),
		Key:    aws.String(key),
		Body:   bufio.NewReader(file),
	}
	if _, err = d.s3Uploader.Upload(params); err != nil {
		return err
	}

	logger.Out.Printf("upload: %s to s3://%s\n", fileStat.Name, filepath.Join(d.bucket, key))
	return nil
}

// CreateS3ObjectStore initializes an S3 ObjectStore
func CreateS3ObjectStore(conf *Config) (storage.ObjectStore, error) {
	s3URL, err := parseS3URI(conf.S3URI)
	if err != nil {
		return nil, err
	}
	awsConfig := initAwsConfig(conf)

	//Init objectstore with s3 session
	s3Client := s3.New(session.New(awsConfig))
	var d storage.ObjectStore
	d = &driver{
		config:       *conf,
		s3Uploader:   s3manager.NewUploaderWithClient(s3Client),
		bucket:       s3URL.Host,
		bucketPrefix: strings.TrimPrefix(s3URL.Path, "/"),
	}
	return d, nil
}

func parseS3URI(s3URI string) (*url.URL, error) {
	s3URL, err := url.Parse(s3URI)
	if err != nil {
		return nil, fmt.Errorf("could not parse s3URI %q", s3URI)
	}
	if s3URL.Scheme != "s3" {
		return nil, fmt.Errorf("s3URI argument does not have valid protocol, should be 's3'")
	}
	if s3URL.Host == "" {
		return nil, fmt.Errorf("s3URI is missing bucket name")
	}
	return s3URL, nil
}

func initAwsConfig(conf *Config) *aws.Config {
	awsConfig := aws.NewConfig()
	creds := credentials.NewChainCredentials([]credentials.Provider{
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     conf.AccessKey,
				SecretAccessKey: conf.SecretKey,
			},
		},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	awsConfig.WithCredentials(creds)
	awsConfig.WithRegion(conf.Region)
	return awsConfig
}
