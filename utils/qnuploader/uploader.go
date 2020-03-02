package qnuploader

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/spf13/viper"
	"mime/multipart"
	"strings"
)

type Uploader struct {
	AccessKey string // 七牛accesskey
	SecretKey string // 七牛secretkey

	VideoBucket  string // 七牛视频资源存储的空间
	AudioBucket  string // 七牛音频资源存储的空间
	StaticBucket string // 七牛其他静态资源存储的空间

	VideoURI  string // 访问视频资源的uri
	AudioURI  string // 访问音频资源的uri
	StaticURI string // 访问静态资源的uri

	CallbackURI   string // 文件上传之后回调的uri
	bucketManager *storage.BucketManager
	mac           *qbox.Mac
	cfg           *storage.Config
}

type UploadCbBody struct {
	Key      string `json:"key"`
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize,string"`
	MimeType string `json:"mimeType"`
	ExtName  string `json:"extName"`
	Name     string `json:"fname"`
	Uploader int64  `json:"uploader, string"`
	Platform int8   `json:"platform, string"`
}

type UploadImageCbBody struct {
	UploadCbBody
	Width  int64 `json:"width, string"`
	Height int64 `json:"height, string"`
}

type UploadVideoCbBody struct {
	UploadCbBody
	Duration float64 `json:"duration, string"`
}

//HandleNullString:处理"null"转成null
func (b UploadCbBody) HandleNullString(buf []byte) []byte {
	return []byte(strings.Replace(string(buf), `"null"`, `null`, -1))
}

// readConf: 读取七牛上传的配置
func (u *Uploader) readConf() {
	u.AccessKey = viper.GetString("qiniu.access_key")
	u.SecretKey = viper.GetString("qiniu.secret_key")

	u.VideoBucket = viper.GetString("qiniu.video_bucket")
	u.AudioBucket = viper.GetString("qiniu.audio_bucket")
	u.StaticBucket = viper.GetString("qiniu.static_bucket")

	u.VideoURI = viper.GetString("qiniu.video_uri")
	u.AudioURI = viper.GetString("qiniu.audio_uri")
	u.StaticURI = viper.GetString("qiniu.static_uri")

	u.CallbackURI = viper.GetString("qiniu.callback_uri")
}

func NewUploader(zone *storage.Zone) *Uploader {
	u := &Uploader{}
	u.readConf()

	if zone == nil {
		zone = &storage.ZoneHuabei
	}

	u.mac = qbox.NewMac(u.AccessKey, u.SecretKey)
	u.cfg = &storage.Config{Zone: zone, UseHTTPS: false}
	u.bucketManager = storage.NewBucketManager(u.mac, u.cfg)
	return u
}

// GetUptoken: 获取七牛上传uptoken
func (u *Uploader) GetUptoken(bucket, cbUrl string, cbMap map[string]string) string {
	body := map[string]string{
		"key":      "$(key)",
		"hash":     "$(etag)",
		"extName":  "$(ext)",
		"fsize":    "$(fsize)",
		"mimeType": "$(mimeType)",
		"fname":    "$(x:filename)",
		"uploader": "$(x:uploader)",
		"platform": "$(x:platform)",
	}

	for k, v := range cbMap {
		body[k] = v
	}

	buf, _ := json.Marshal(body)
	var putPolicy storage.PutPolicy

	if cbUrl == "" {
		putPolicy = storage.PutPolicy{
			Scope:      bucket,
			ReturnBody: string(buf),
		}
	} else {
		putPolicy = storage.PutPolicy{
			Scope:            bucket,
			CallbackBody:     string(buf),
			CallbackBodyType: "application/json",
			CallbackURL:      cbUrl,
		}
	}

	putPolicy.Expires = 12 * 3600
	return putPolicy.UploadToken(u.mac)
}

// GetStaticUptoken: 获取七牛静态资源的上传uptoken
func (u *Uploader) GetStaticUptoken(cbUrl string, cbMap map[string]string) string {
	if cbMap == nil {
		cbMap = make(map[string]string)
	}

	return u.GetUptoken(u.StaticBucket, cbUrl, cbMap)
}

// GetImageUptoken: 获取七牛图片资源的上传uptoken
func (u *Uploader) GetImageUptoken(cbUrl string, cbMap map[string]string) string {
	if cbMap == nil {
		cbMap = make(map[string]string)
	}

	cbMap["width"] = "$(imageInfo.width)"
	cbMap["height"] = "$(imageInfo.height)"
	return u.GetUptoken(u.StaticBucket, cbUrl, cbMap)
}

// GetVideoUptoken: 获取七牛视频上传的uptoken
func (u *Uploader) GetVideoUptoken(cbUrl string, cbMap map[string]string) string {
	if cbMap == nil {
		cbMap = make(map[string]string)
	}

	cbMap["duration"] = "$(avinfo.video.duration)"
	return u.GetUptoken(u.VideoBucket, cbUrl, cbMap)
}

// FetchNetResToBucket: 获取网络上的资源到空间中
func (u *Uploader) FetchNetResToBucket(resUrl, bucket, newKey string) (storage.FetchRet, error) {
	return u.bucketManager.Fetch(resUrl, bucket, newKey)
}

// DeleteBucketRes: 删除空间中的资源, resType资源类型: static, video, audio
func (u *Uploader) DeleteBucketRes(resType, key string) error {
	var bucket string
	switch resType {
	case "static":
		bucket = u.StaticBucket
	case "video":
		bucket = u.VideoBucket
	case "audio":
		bucket = u.AudioBucket
	default:
		return errors.New("resType not exists, please input resType: static, video, audio")
	}

	return u.bucketManager.Delete(bucket, key)
}

// GetFileInfo: 获取文件信息
// 如果没有文件的话，会返回err: no such file or directory
// resType资源类型: static, video, audio
func (u *Uploader) GetBucketResInfo(resType, key string) (storage.FileInfo, error) {
	var bucket string
	switch resType {
	case "static":
		bucket = u.StaticBucket
	case "video":
		bucket = u.VideoBucket
	case "audio":
		bucket = u.AudioBucket
	default:
		return storage.FileInfo{}, errors.New("resType not exists, please input resType: static, video, audio")
	}

	return u.bucketManager.Stat(bucket, key)
}

func (u *Uploader) UploadFormFile(body interface{}, params map[string]string, uptoken, key string, file *multipart.FileHeader) (err error) {
	formUploader := storage.NewFormUploader(u.cfg)
	f, _ := file.Open()
	defer func() {
		if e := f.Close(); e != nil {
			err = e
		}
	}()

	putExtra := storage.PutExtra{Params: params}
	err = formUploader.Put(context.Background(), body, uptoken, key, f, file.Size, &putExtra)
	return
}
