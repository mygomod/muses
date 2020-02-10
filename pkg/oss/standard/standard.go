package standard

import (
	"io"
)

type Oss interface {
	PutObject(dstPath string, reader io.Reader, options ...Option) error
	PutObjectFromFile(dstPath, srcPath string, options ...Option) error
	GetObject(dstPath string, options ...Option) (io.ReadCloser, error)
	GetObjectToFile(dstPath, srcPath string, options ...Option) error
	DeleteObject(dstPath string) error
	DeleteObjects(dstPaths []string, options ...Option) (DeleteObjectsResult, error)
	IsObjectExist(dstPath string) (bool, error)
	ListObjects(options ...Option) (ListObjectsResult, error)
	SignURL(dstPath string, method string, expiredInSec int64, options ...Option) (string, error)
}
