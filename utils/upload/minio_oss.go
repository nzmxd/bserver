package upload

import (
	"bytes"
	"context"
	"errors"
	"git.zingfront.cn/liubin/bserver/global"
	"git.zingfront.cn/liubin/bserver/utils"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

var MinioClient *Minio // 优化性能，但是不支持动态配置

type Minio struct {
	Client *minio.Client
	bucket string
}

func GetMinio(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*Minio, error) {
	if MinioClient != nil {
		return MinioClient, nil
	}
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL, // Set to true if using https
	})
	if err != nil {
		return nil, err
	}
	// 尝试创建bucket
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			// log.Printf("We already own %s\n", bucketName)
		} else {
			return nil, err
		}
	}
	MinioClient = &Minio{Client: minioClient, bucket: bucketName}
	return MinioClient, nil
}

func (m *Minio) UploadFile(file *multipart.FileHeader) (filePathres, key string, uploadErr error) {
	f, openError := file.Open()
	// mutipart.File to os.File
	if openError != nil {
		global.LOG.Error("function file.Open() Failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}

	filecontent := bytes.Buffer{}
	_, err := io.Copy(&filecontent, f)
	if err != nil {
		global.LOG.Error("读取文件失败", zap.Any("err", err.Error()))
		return "", "", errors.New("读取文件失败, err:" + err.Error())
	}
	f.Close() // 创建文件 defer 关闭

	// 对文件名进行加密存储
	ext := filepath.Ext(file.Filename)
	filename := utils.MD5V([]byte(strings.TrimSuffix(file.Filename, ext))) + ext
	if global.CONFIG.Minio.BasePath == "" {
		filePathres = "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + filename
	} else {
		filePathres = global.CONFIG.Minio.BasePath + "/" + time.Now().Format("2006-01-02") + "/" + filename
	}

	// 设置超时10分钟
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	// Upload the file with PutObject   大文件自动切换为分片上传
	info, err := m.Client.PutObject(ctx, global.CONFIG.Minio.BucketName, filePathres, &filecontent, file.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		global.LOG.Error("上传文件到minio失败", zap.Any("err", err.Error()))
		return "", "", errors.New("上传文件到minio失败, err:" + err.Error())
	}
	return global.CONFIG.Minio.BucketUrl + "/" + info.Key, filePathres, nil
}

func (m *Minio) UploadLocalFile(localFilePath string, objectNamePrefix string) (string, error) {
	file, err := os.Open(localFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	// 自动检测 Content-Type
	buf := make([]byte, 512)
	_, _ = file.Read(buf)
	contentType := http.DetectContentType(buf)
	// 重置文件指针
	_, _ = file.Seek(0, io.SeekStart)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()
	datePath := time.Now().Format("2006-01-02")
	parts := []string{global.CONFIG.Minio.BasePath}
	if objectNamePrefix != "" {
		parts = append(parts, objectNamePrefix)
	}
	parts = append(parts, datePath, fileInfo.Name())
	objectName := path.Join(parts...)
	_, err = m.Client.PutObject(ctx, m.bucket, objectName, file, fileInfo.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return global.CONFIG.Minio.BucketUrl + "/" + objectName, err
}

func (m *Minio) DeleteFile(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Delete the object from MinIO
	err := m.Client.RemoveObject(ctx, m.bucket, key, minio.RemoveObjectOptions{})
	return err
}

func (m *Minio) GetObjectToFile(objectName string, filePath string) error {
	return m.Client.FGetObject(context.Background(), m.bucket, objectName, filePath, minio.GetObjectOptions{})
}

func (m *Minio) GetPresignedDownloadURL(objectName string) (string, error) {
	presignedURL, err := m.Client.PresignedGetObject(
		context.Background(), m.bucket, objectName, time.Hour, nil,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
