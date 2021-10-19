package s3

import (
	"errors"
	"log"
	"os"

	"github.com/minio/minio-go/v6"
)

// Conf represents a set of parameters to configure a s3 client.
type Conf struct {
	Key      string
	Secret   string
	Region   string
	Secure   bool
	Endpoint string
}

// Client represents a s3 client.
type Client struct {
	mClient *minio.Client
}

type Bucket struct {
	Name string
}

func NewClientWithBasicCredentials(conf Conf) (Client, error) {
	if conf.Key == "" {
		return Client{}, errors.New("the aws access key ID cannot be empty")
	}
	if conf.Secret == "" {
		return Client{}, errors.New("the aws access secret key cannot be empty")
	}
	if conf.Region == "" {
		return Client{}, errors.New("the aws region cannot be empty")
	}
	if conf.Endpoint == "" {
		conf.Endpoint = "s3.amazonaws.com"
	}
	c, err := minio.NewWithRegion(conf.Endpoint, conf.Key, conf.Secret, conf.Secure, conf.Region)
	return Client{mClient: c}, err
}

func GetClient() Client {
	client, err := NewClientWithBasicCredentials(Conf{
		Key:    os.Getenv("AWS_ACCESS_KEY_ID"),
		Secret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Region: os.Getenv("AWS_REGION"),
		Secure: true,
	})
	if err != nil {
		log.Fatalln("error connecting to s3", err)
	}
	return client
}

// Get Downloads and saves the object as a file in the local filesystem.
func (c Client) Get(bucketName, objectName, fileName string) error {
	return c.mClient.FGetObject(bucketName, objectName, fileName, minio.GetObjectOptions{})
}

func (c Client) GetAllObjectKeys(bucketName string) []string {
	result := make([]string, 0)
	// Create a done channel to control 'ListObjectsV2' go routine.
	doneCh := make(chan struct{})

	// Indicate to our routine to exit cleanly upon return.
	defer close(doneCh)

	// List all objects from a bucket-name with a matching prefix.
	for object := range c.mClient.ListObjectsV2(bucketName, "", false, doneCh) {
		if object.Err != nil {
			log.Println(object.Err)
			return result
		}
		result = append(result, object.Key)
	}
	return result
}
