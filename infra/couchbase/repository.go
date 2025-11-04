package couchbase

import (
	"context"
	"errors"
	"time"

	"github.com/couchbase/gocb/v2"
	"go.uber.org/zap"

	"microservicetest/domain"
)

type CouchbaseRepository struct {
	cluster *gocb.Cluster
	bucket  *gocb.Bucket
}

func NewCouchbaseRepository(couchbaseUrl string, username string, password string) *CouchbaseRepository {
	cluster, err := gocb.Connect(couchbaseUrl, gocb.ClusterOptions{
		TimeoutsConfig: gocb.TimeoutsConfig{
			ConnectTimeout: 3 * time.Second,
			KVTimeout:      3 * time.Second,
			QueryTimeout:   3 * time.Second,
		},
		Authenticator: gocb.PasswordAuthenticator{
			Username: username,
			Password: password,
		},
		Transcoder: gocb.NewJSONTranscoder(),
		Tracer:     tracer,
	})
	if err != nil {
		zap.L().Fatal("Failed to connect to couchbase", zap.Error(err))
	}

	bucket := cluster.Bucket("products")
	bucket.WaitUntilReady(3*time.Second, &gocb.WaitUntilReadyOptions{})

	return &CouchbaseRepository{
		cluster: cluster,
		bucket:  bucket,
	}
}

func (r *CouchbaseRepository) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	data, err := r.bucket.DefaultCollection().Get(id, &gocb.GetOptions{
		Timeout:    3 * time.Second,
		Context:    ctx,
	})
	if err != nil {
		if errors.Is(err, gocb.ErrDocumentNotFound) {
			return nil, errors.New("product not found")
		}

		zap.L().Error("Failed to get product", zap.Error(err))
		return nil, err
	}

	var product domain.Product
	if err := data.Content(&product); err != nil {
		zap.L().Error("Failed to unmarshal product", zap.Error(err))
		return nil, err
	}

	return &product, nil
}

func (r *CouchbaseRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	_, err := r.bucket.DefaultCollection().Insert(product.ID, product, &gocb.InsertOptions{
		Timeout: 3 * time.Second,
		Context: ctx,
	})
	if err != nil {
		zap.L().Error("Failed to create product", zap.Error(err))
		return err
	}

	return nil
}

func (r *CouchbaseRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	_, err := r.bucket.DefaultCollection().Replace(product.ID, product, &gocb.ReplaceOptions{
		Timeout:    3 * time.Second,
		Context:    ctx,
	})
	if err != nil {
		zap.L().Error("Failed to update product", zap.Error(err))
		return err
	}

	return nil
}
