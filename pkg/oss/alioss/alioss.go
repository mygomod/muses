package alioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/goecology/muses/pkg/oss/standard"
	"io"
	"os"
)

type Client struct {
	b        *oss.Bucket
	isDelete bool
}

func NewOss(endpoints, accessKeyId, accessKeySecret, bucketName string, isDelete bool) (client standard.Oss, err error) {
	c, e := oss.New(
		endpoints, accessKeyId, accessKeySecret,
	)
	if e != nil {
		return
	}
	b, e := c.Bucket(bucketName)
	if e != nil {
		return
	}
	client = &Client{
		b:        b,
		isDelete: isDelete,
	}
	return
}

func (c *Client) PutObject(dstPath string, reader io.Reader, options ...standard.Option) error {
	return c.b.PutObject(dstPath, reader)
}

func (c *Client) SignURL(dstPath string, method string, expiredInSec int64, options ...standard.Option) (string, error) {
	return c.b.SignURL(dstPath, oss.HTTPMethod(method), 120)
}

func (c *Client) PutObjectFromFile(dstPath, srcPath string, options ...standard.Option) (err error) {
	err = c.b.PutObjectFromFile(dstPath, srcPath)
	if err != nil {
		return
	}
	if c.isDelete {
		err = os.Remove(srcPath)
	}
	return

}

func (c *Client) GetObject(dstPath string, options ...standard.Option) (io.ReadCloser, error) {
	panic("implement me")
}

func (c *Client) GetObjectToFile(dstPath, srcPath string, options ...standard.Option) error {
	panic("implement me")
}

func (c *Client) DeleteObject(dstPath string) error {
	panic("implement me")
}

func (c *Client) DeleteObjects(dstPaths []string, options ...standard.Option) (standard.DeleteObjectsResult, error) {
	panic("implement me")
}

func (c *Client) IsObjectExist(dstPath string) (bool, error) {
	panic("implement me")
}

func (c *Client) ListObjects(options ...standard.Option) (standard.ListObjectsResult, error) {
	panic("implement me")
}
