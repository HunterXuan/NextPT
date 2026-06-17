package sys

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"server/internal/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gogf/gf/v2/frame/g"
)

type sSysStorage struct {
	driver    string
	localPath string
	// S3 Client
	s3Client *s3.Client
	s3Bucket string
}

func init() {
	service.RegisterSysStorage(NewSysStorage())
}

func NewSysStorage() *sSysStorage {
	ctx := context.Background()

	// 默认 driver 为 local
	driver := g.Cfg().MustGet(ctx, "storage.driver").String()
	if driver == "" {
		driver = "local"
	}

	storage := &sSysStorage{
		driver: driver,
	}

	if driver == "s3" {
		storage.initS3(ctx)
	} else {
		storage.initLocal(ctx)
	}

	return storage
}

func (s *sSysStorage) initLocal(ctx context.Context) {
	s.localPath = g.Cfg().MustGet(ctx, "storage.local.path").String()
	if s.localPath == "" {
		s.localPath = "./storage/uploads"
	}
	// 确保根目录存在
	if err := os.MkdirAll(s.localPath, os.ModePerm); err != nil {
		g.Log().Warningf(ctx, "failed to create local storage directory: %v", err)
	}
}

func (s *sSysStorage) initS3(ctx context.Context) {
	s3Id := g.Cfg().MustGet(ctx, "storage.s3.id").String()
	s3Key := g.Cfg().MustGet(ctx, "storage.s3.key").String()
	region := g.Cfg().MustGet(ctx, "storage.s3.region").String()
	endpoint := g.Cfg().MustGet(ctx, "storage.s3.endpoint").String()
	bucket := g.Cfg().MustGet(ctx, "storage.s3.bucket").String()

	if bucket == "" || s3Id == "" || s3Key == "" {
		g.Log().Fatal(ctx, "S3 configuration is missing (bucket/id/key) while storage.driver is 's3'")
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Id, s3Key, "")),
	)
	if err != nil {
		g.Log().Fatalf(ctx, "failed to load s3 config: %+v", err)
	}

	if endpoint != "" {
		cfg.BaseEndpoint = aws.String(endpoint)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	s.s3Client = client
	s.s3Bucket = bucket
}

// Upload 将数据上传到存储引擎
func (s *sSysStorage) Upload(ctx context.Context, key string, data []byte, contentType string) error {
	if s.driver == "s3" {
		return s.uploadS3(ctx, key, data, contentType)
	}
	return s.uploadLocal(ctx, key, data, contentType)
}

// Download 从存储引擎下载文件并保存到本地指定路径
func (s *sSysStorage) Download(ctx context.Context, key string, destPath string) error {
	if s.driver == "s3" {
		return s.downloadS3(ctx, key, destPath)
	}
	return s.downloadLocal(ctx, key, destPath)
}

// Delete 从存储引擎删除文件
func (s *sSysStorage) Delete(ctx context.Context, key string) error {
	if s.driver == "s3" {
		return s.deleteS3(ctx, key)
	}
	return s.deleteLocal(ctx, key)
}

// GetLocalPath 获取文件的本地访问路径。
// 如果是 local 驱动，返回文件在磁盘的真实绝对路径；
// 如果是 s3 等远程驱动，返回本地缓存目录中的绝对路径，供调用方检查是否存在并按需调用 Download。
func (s *sSysStorage) GetLocalPath(ctx context.Context, key string) string {
	var baseDir string
	if s.driver == "local" {
		baseDir = s.localPath
	} else {
		// 远程驱动的本地缓存目录
		baseDir = "./storage/cache"
	}

	absPath, err := filepath.Abs(filepath.Join(baseDir, key))
	if err == nil {
		return absPath
	}
	return filepath.Join(baseDir, key)
}

func (s *sSysStorage) uploadLocal(ctx context.Context, key string, data []byte, contentType string) error {
	destPath := filepath.Join(s.localPath, key)
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write local file: %w", err)
	}
	return nil
}

func (s *sSysStorage) downloadLocal(ctx context.Context, key string, destPath string) error {
	srcPath := filepath.Join(s.localPath, key)
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open local storage file: %w", err)
	}
	defer srcFile.Close()

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create dest directory: %w", err)
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create dest file: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy to dest file: %w", err)
	}
	return nil
}

func (s *sSysStorage) uploadS3(ctx context.Context, key string, data []byte, contentType string) error {
	if s.s3Client == nil || s.s3Bucket == "" {
		return fmt.Errorf("S3 client is not initialized due to missing configuration")
	}
	body := bytes.NewReader(data)
	_, err := s.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return fmt.Errorf("failed to put object to s3: %w", err)
	}
	return nil
}

func (s *sSysStorage) downloadS3(ctx context.Context, key string, destPath string) error {
	if s.s3Client == nil || s.s3Bucket == "" {
		return fmt.Errorf("S3 client is not initialized due to missing configuration")
	}
	out, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to get object from s3: %w", err)
	}
	defer out.Body.Close()

	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create local directory: %w", err)
	}

	f, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, out.Body); err != nil {
		return fmt.Errorf("failed to write object to local file: %w", err)
	}
	return nil
}

func (s *sSysStorage) deleteLocal(ctx context.Context, key string) error {
	destPath := filepath.Join(s.localPath, key)
	if err := os.Remove(destPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete local file: %w", err)
	}
	return nil
}

func (s *sSysStorage) deleteS3(ctx context.Context, key string) error {
	if s.s3Client == nil || s.s3Bucket == "" {
		return fmt.Errorf("S3 client is not initialized due to missing configuration")
	}
	_, err := s.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from s3: %w", err)
	}
	return nil
}
