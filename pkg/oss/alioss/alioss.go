package alioss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/goecology/muses/pkg/oss/standard"
	"io"
	"io/ioutil"
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

func (c *Client) GetObject(dstPath string, options ...standard.Option) (ouput []byte, err error) {
	var reader io.ReadCloser
	reader, err = c.b.GetObject(dstPath)
	return ioutil.ReadAll(reader)
}

func (c *Client) GetObjectToFile(dstPath, srcPath string, options ...standard.Option) error {
	panic("implement me")
}

func (c *Client) DeleteObject(dstPath string) (err error) {
	err = c.b.DeleteObject(dstPath)
	if err != nil {
		return
	}
	return
}

func (c *Client) DeleteObjects(dstPaths []string, options ...standard.Option) (output standard.DeleteObjectsResult, err error) {
	var resp oss.DeleteObjectsResult
	resp, err = c.b.DeleteObjects(dstPaths)
	if err != nil {
		return
	}
	output = standard.DeleteObjectsResult{
		Local:          resp.XMLName.Local,
		Space:          resp.XMLName.Space,
		DeletedObjects: resp.DeletedObjects,
	}
	return
}

func (c *Client) IsObjectExist(dstPath string) (bool, error) {
	panic("implement me")
}

func (c *Client) ListObjects(options ...standard.Option) (standard.ListObjectsResult, error) {
	panic("implement me")
}
