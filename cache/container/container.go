package container

import (
	"fmt"

	"github.com/x1nchen/portainer-cli/cache/internal"
	climodel "github.com/x1nchen/portainer-cli/model"
	"github.com/x1nchen/portainer/model"
	bolt "go.etcd.io/bbolt"
)

const (
	// BucketName represents the name of the bucket where this service stores data.
	BucketName = "container"
)

// Service represents a service for managing endpoint data.
type Service struct {
	db *bolt.DB
}

// NewService creates a new instance of a service.
func NewService(db *bolt.DB) (*Service, error) {
	err := internal.CreateBucket(db, BucketName)
	if err != nil {
		return nil, err
	}

	return &Service{
		db: db,
	}, nil
}

// Endpoint returns an endpoint by ID.
func (service *Service) GetContain(ID int) (*model.DockerContainer, error) {
	var mc model.DockerContainer
	identifier := internal.Itob(int(ID))

	err := internal.GetObject(service.db, BucketName, identifier, &mc)
	if err != nil {
		return nil, err
	}

	return &mc, nil
}

// UpdateEndpoint updates an endpoint.
func (service *Service) UpdateEndpoint(ID int, endpoint *model.Endpoint) error {
	identifier := internal.Itob(int(ID))
	return internal.UpdateObject(service.db, BucketName, identifier, endpoint)
}

// DeleteEndpoint deletes an endpoint.
func (service *Service) DeleteEndpoint(ID int) error {
	identifier := internal.Itob(int(ID))
	return internal.DeleteObject(service.db, BucketName, identifier)
}

// Endpoints return an array containing all the endpoints.
func (service *Service) Endpoints() ([]model.Endpoint, error) {
	var endpoints = make([]model.Endpoint, 0)

	err := service.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))

		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			var endpoint model.Endpoint
			err := internal.UnmarshalObjectWithJsoniter(v, &endpoint)
			if err != nil {
				return err
			}
			endpoints = append(endpoints, endpoint)
		}

		return nil
	})

	return endpoints, err
}

// GetNextIdentifier returns the next identifier for an endpoint.
func (service *Service) GetNextIdentifier() int {
	return internal.GetNextIdentifier(service.db, BucketName)
}

// 批量更新 container，key 格式是容器的 name#container_id
func (service *Service) BatchUpdateContainers(containers ...climodel.ContainerExtend) error {
	return service.db.Batch(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BucketName))
		for _, container := range containers {
			data, err := internal.MarshalObject(container)
			if err != nil {
				return err
			}

			var containerName string

			if len(container.Names) > 0 {
				if len(container.Names[0]) > 0 {
					// 注意：容器的名字有前缀 "/"，如 /node-api
					containerName = container.Names[0][1:]
				}
			}
			key := fmt.Sprintf("%s#%d", containerName, container.EndpointId)
			err = bucket.Put(internal.StringToBytes(key), data)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
