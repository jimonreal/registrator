package couchbase

import (
	"log"
	"net"
	"net/url"
	"os"
	"strconv"

	couchbase "github.com/couchbaselabs/gocb"
	"github.com/gliderlabs/registrator/bridge"
)

func init() {
	bridge.Register(new(Factory), "couchbase")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
	BUCKET_NAME := os.Getenv("REGISTRATOR_BUCKET")
	PASSWORD := os.Getenv("REGISTRATOR_BUCKET_PASSWORD")

	client, err := couchbase.Connect("couchbase://"+uri.Host, BUCKET_NAME, PASSWORD)

	if err != nil {
		log.Fatalf("couchbase: error connecting to %v: %v", uri, err)
	}

	bucket, err := pool.GetBucket(BUCKET_NAME)

	if err != nil {
		log.Fatalf("couchbase: error can't get bucket [%v]: [%v]", BUCKET_NAME, err)
	}

	return &CouchbaseAdapter{client: client, pool: pool, bucket: bucket, path: uri.Path}
}

type CouchbaseAdapter struct {
	client couchbase.Client
	bucket *couchbase.Bucket

	path string
}

func (r *CouchbaseAdapter) Ping() error {
	err := r.bucket.Refresh()
	if err != nil {
		return err
	}
	return nil
}

func (r *CouchbaseAdapter) Register(service *bridge.Service) error {
	path := r.path + "/" + service.Name + "/" + service.ID
	port := strconv.Itoa(service.Port)
	addr := net.JoinHostPort(service.IP, port)

	var err error

	err = r.bucket.Set(path, service.TTL, map[string]interface{}{addr: 1})

	if err != nil {
		log.Printf("couchbase: failed to register service: %v\n", err)
	}
	return err
}

func (r *CouchbaseAdapter) Deregister(service *bridge.Service) error {
	path := r.path + "/" + service.Name + "/" + service.ID

	var err error
	if r.bucket != nil {
		err = r.bucket.Delete(path)
	}

	if err != nil {
		log.Printf("couchbase: failed to deregister service: %v\n", err)
	}
	return err
}

func (r *CouchbaseAdapter) Refresh(service *bridge.Service) error {
	return r.Register(service)
}
