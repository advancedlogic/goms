package modules

import (
	"github.com/advancedlogic/goms/pkg/models"
	"github.com/minio/minio-go"
	"io/ioutil"
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
		WithAccessKey(mb.GetStringOrDefault("store.accessKey", "")).
		WithSecretKey(mb.GetStringOrDefault("store.secretKey", ""))
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

func (mb *MinioBuilder) WithSecretKey(secretKey string) *MinioBuilder {
	mb.secretKey = secretKey
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

	reader, err := client.GetObject(m.bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	if value, err := ioutil.ReadAll(reader); err == nil {
		return string(value), nil
	} else {
		return nil, err
	}
}

func (m *Minio) Update(key string, data interface{}) error {
	return m.Create(key, data)
}

func (m *Minio) Delete(key string) error {
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return err
	}

	return client.RemoveObject(m.bucket, key)
}

func (m *Minio) List() ([]interface{}, error) {
	client, err := minio.New(m.endpoint, m.accessKey, m.secretKey, false)
	if err != nil {
		return nil, err
	}
	doneCh := make(chan struct{})
	defer close(doneCh)
	values := make([]interface{}, 0)
	for value := range client.ListObjectsV2(m.bucket, "", true, doneCh) {
		values = append(values, value)
	}
	return values, nil
}
