package modules

import (
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/minio/minio-go"
	"io"
	"log"
	"strings"
)

type Minio struct {
	location  string
	bucket    string
	endpoint  string
	accessKey string
	secretKey string
}

type MinioBuilder struct {
	*models.Environment
	*Minio
}

func NewMinioBuilder(environment *models.Environment) *MinioBuilder {
	mb := &MinioBuilder{
		Environment: environment,
		Minio:       &Minio{},
	}
	return mb.
		WithLocation(mb.GetStringOrDefault("store.location", "default")).
		WithBucket(mb.GetStringOrDefault("store.bucket", "default")).
		WithEndpoint(mb.GetStringOrDefault("store.endpoint", "localhost:9000")).
		WithAccessKey(mb.GetString("store.accessKey")).
		WithSecretKey(mb.GetString("store.secretKey"))
}

func (mb *MinioBuilder) WithLocation(name string) *MinioBuilder {
	mb.location = name
	return mb
}

func (mb *MinioBuilder) WithBucket(name string) *MinioBuilder {
	mb.bucket = name
	return mb
}

func (mb *MinioBuilder) WithEndpoint(endpoint string) *MinioBuilder {
	mb.endpoint = endpoint
	return mb
}

func (mb *MinioBuilder) WithAccessKey(accessKey string) *MinioBuilder {
	mb.accessKey = accessKey
	return mb
}

func (mb *MinioBuilder) WithSecretKey(accessKey string) *MinioBuilder {
	mb.accessKey = accessKey
	return mb
}

func (mb *MinioBuilder) Build() (*Minio, error) {
	client, err := minio.New(mb.endpoint, mb.accessKey, mb.secretKey, false)
	if err != nil {
		return nil, err
	}
	err = client.MakeBucket(mb.bucket, mb.location)
	if err != nil {
		if exists, err := client.BucketExists(mb.bucket); exists && err == nil {
			mb.Warnf("Bucket %s already exists", mb.bucket)
		} else {
			log.Fatal(err)
		}
		mb.Error(err)
	} else {
		mb.Infof("Successfully created bucket %s", mb.bucket)
	}

	return mb.Minio, nil
}

func (m *Minio) Create(key string, data interface{}) error {
	reader := strings.NewReader(data.(string))
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return err
	}
	_, err = client.PutObject(m.bucket, key, reader, -1, minio.PutObjectOptions{
		ContentType: "plain/txt",
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Minio) Read(key string) (interface{}, error) {
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return "", err
	}

	obj, err := client.GetObject(m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer obj.Close()
	var b io.ByteWriter
	bobj, err := io.Copy(b, obj)

}
