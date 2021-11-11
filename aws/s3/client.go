package s3

import (
	"context"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/spf13/viper"
)

type Options struct {
	Endpoint     string                  `mapstructure:"endpoint_url"`
	Bucket       string                  `mapstructure:"bucket"`
	UsePathStyle bool                    `mapstructure:"use_path_style"`
	BucketACL    s3Types.BucketCannedACL `mapstructure:"bucket_acl"`
	ObjectACL    s3Types.ObjectCannedACL `mapstructure:"object_acl"`
}

type Client struct {
	*s3.Client
	defaultBucket string
	bucketAcl     s3Types.BucketCannedACL
	objectAcl     s3Types.ObjectCannedACL
}

func newClient(opts *Options, awsCfg *aws.Config) *Client {
	s3Opts := s3.Options{
		Region:        awsCfg.Region,
		HTTPClient:    awsCfg.HTTPClient,
		Credentials:   awsCfg.Credentials,
		APIOptions:    awsCfg.APIOptions,
		Logger:        awsCfg.Logger,
		ClientLogMode: awsCfg.ClientLogMode,
		UsePathStyle:  opts.UsePathStyle,
	}

	if opts.Endpoint != "" {
		s3Opts.EndpointResolver = s3.EndpointResolverFromURL(opts.Endpoint)
	}

	s3Client := s3.New(s3Opts)

	return &Client{
		Client:        s3Client,
		defaultBucket: opts.Bucket,
		bucketAcl:     opts.BucketACL,
		objectAcl:     opts.ObjectACL,
	}
}

func New(opts *Options, awsCfg *aws.Config) *Client {
	return newClient(opts, awsCfg)
}

// NewFromViper 从 viper 生成 Client 实例, 返回 Client 实例和 error
func NewFromViper(v *viper.Viper, awsCfg *aws.Config) (*Client, error) {
	return NewFromViperByKey(v, "s3", awsCfg)
}

func NewFromViperByKey(v *viper.Viper, key string, awsCfg *aws.Config) (*Client, error) {
	opts := new(Options)
	if err := v.UnmarshalKey(key, opts); err != nil {
		return nil, err
	}
	return New(opts, awsCfg), nil
}

// Open 打开默认 bucket 下的对象并返回一个 io.ReadCloser 和 error
func (c *Client) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	op, err := c.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.defaultBucket,
		Key:    &key,
	})
	if err != nil {
		return nil, err
	}
	return op.Body, nil
}

func (c *Client) Put(ctx context.Context, r io.Reader, key string) error {
	_, err := c.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.defaultBucket,
		Body:   r,
		Key:    &key,
		ACL:    c.objectAcl,
	})
	return err
}

func (c *Client) Download(ctx context.Context, key string, local string) error {
	reader, err := c.Open(ctx, key)
	if err != nil {
		return err
	}
	f, err := os.Create(local)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Delete(ctx context.Context, key string) error {
	_, err := c.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.defaultBucket,
		Key:    &key,
	})
	return err
}

func (c *Client) Listdir(ctx context.Context, dir string) ([]string, error) {
	op, err := c.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: &c.defaultBucket,
		Prefix: &dir,
	})
	if err != nil {
		return nil, err
	}
	list := make([]string, 0, len(op.Contents))
	for _, obj := range op.Contents {
		list = append(list, *obj.Key)
	}
	return list, nil
}

func (c *Client) Upload(ctx context.Context, local string, key string) error {
	f, err := os.Open(local)
	if err != nil {
		return err
	}
	return c.Put(ctx, f, key)
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &c.defaultBucket, Key: &key})
	if err != nil {
		return false, err
	}
	return true, nil
}
