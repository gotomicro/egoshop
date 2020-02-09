package service

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/goecology/egoshop/appgo/model/constx"
	"github.com/goecology/muses/pkg/logger"
	"github.com/satori/uuid"
	"go.uber.org/zap"
	"strings"
	"time"
)

var Oss *alioss

type alioss struct {
	*oss.Client
	buckets map[string]*bucket
	logger  *logger.Client
}

type bucket struct {
	*oss.Bucket
}

// https://help.aliyun.com/document_detail/31857.html
// https://github.com/aliyun/aliyun-oss-go-sdk/blob/master/README-CN.md
func InitOssCli(endpoints, accessKeyId, accessKeySecret, bucketName string, logger *logger.Client) {
	c, e := oss.New(
		endpoints, accessKeyId, accessKeySecret,
	)
	if e != nil {
		logger.Warn("oss fail", zap.Error(e))
		return
	}
	Oss = &alioss{
		Client:  c,
		logger:  logger,
		buckets: make(map[string]*bucket),
	}

	Oss.newBucket(bucketName, bucketName) // 公共读，通过bucketUrl访问ojb，无过期时间
}

func (o *alioss) newBucket(name string, bucketName string) {
	b, e := Oss.Bucket(bucketName)
	if e != nil {
		o.logger.Warn("oss fail", zap.Error(e))
		return
	}
	o.buckets[name] = &bucket{b}
}

func (o *alioss) Buckets() []oss.BucketProperties {
	lsRes, e := Oss.ListBuckets()
	if e != nil {
		o.logger.Warn("oss fail", zap.Error(e))
		return nil
	}
	return lsRes.Buckets
}

// Buck 返回bucket，pubr为true则返回公共读bucket，否原则返回私有读写bucket
func (o *alioss) Buck(bucketName string) *bucket {
	return Oss.buckets[bucketName]
}

// PutObj 使用pri bucket上传对象
func (o *alioss) PutObj(bucketName string, key string, filePath string) error {
	return o.Buck(bucketName).PutObjectFromFile(key, filePath)
}

func (b *bucket) PutObj(key string, filePath string) error {
	return b.PutObjectFromFile(key, filePath)
}

// GetObj 使用pri bucket直接返回对象
func (o *alioss) GetObj(bucketName string, key string, filePath string) error {
	return o.Buck(bucketName).GetObjectToFile(key, filePath)
}

func (b *bucket) GetObj(key string, filePath string) error {
	return b.GetObjectToFile(key, filePath)
}

func (o *alioss) GetObjURL(bucketName string, key string) (string, error) {
	return o.Buck(bucketName).GetObjURL(key)
}

// GetObjURL 使用pri bucket返回对象url
func (o *alioss) GetObjURLBySign(bucketName string, key string) (string, error) {
	return o.Buck(bucketName).SignURL(key, oss.HTTPGet, 120)
}

func (b *bucket) GetObjURL(key string) (string, error) {
	return b.SignURL(key, oss.HTTPGet, 120)
}

func (o *alioss) OssPrefix(ossType int) string {
	resp, _ := constx.OssMap[ossType]
	return resp
}

// Key 随机生成一个key
func (o *alioss) Key(prefix string) string {
	month := time.Now().Format("200601")
	// Reossv2上传报错：Thespecifiedobjectisnotvalid.
	//	object路径开头不能与“/”
	return prefix + "/" + month + "/" + strings.ReplaceAll(uuid.NewV4().String(), "-", "") + ".jpg"
}
