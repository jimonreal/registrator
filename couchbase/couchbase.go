package couchbase

import (
	"encoding/json"
	"fmt"
	couchbase "github.com/couchbaselabs/gocb"
	"github.com/gliderlabs/registrator/bridge"
	"log"
	"net/url"
	"os"
)

func init() {
	bridge.Register(new(Factory), "couchbase")
}

type Factory struct{}

func (f *Factory) New(uri *url.URL) bridge.RegistryAdapter {
	BUCKET_NAME := os.Getenv("CB_BUCKET")
	BUCKET_PASSWORD := os.Getenv("CB_BUCKET_PASSWORD")

	fmt.Println("Uri Host:", uri.Host)

	client, err := couchbase.Connect("http://" + uri.Host)

	if err != nil {
		log.Fatalf("couchbase: error connecting to %v: %v", uri, err)
	}

	bucket, err := client.OpenBucket(BUCKET_NAME, BUCKET_PASSWORD)

	if err != nil {
		log.Fatalf("couchbase: error can't get bucket [%v]: [%v]", BUCKET_NAME, err)
	}

	return &CouchbaseAdapter{client: client, bucket: bucket, path: uri.Path}
}

type CouchbaseAdapter struct {
	client *couchbase.Cluster
	bucket *couchbase.Bucket

	path string
}

func (r *CouchbaseAdapter) Ping() error {
	_, err := r.bucket.Get("version", nil) //check for another approach
	if err != nil {
		return err
	}
	return nil
}

func (r *CouchbaseAdapter) Register(service *bridge.Service) error {
	var err error
	log.Println("Register...")

	lbc := LoadBalancerConfig{
		Enable:  mapDefault(service.Attrs, "SERVICE_ENABLE", "true"),
		Mode:    mapDefault(service.Attrs, "SERVICE_MODE", "http"),
		Balance: mapDefault(service.Attrs, "SERVICE_BALANCER_ALGORITHM", "roundrobin"),
		Param:   mapDefault(service.Attrs, "SERVICE_BALANCER_PARAMS", ""),
	}
	doc := Document{
		ContainerId:   service.ID,
		ContainerType: "Test",
		Location:      service.Origin,
		MetaData:      service.Attrs,
		Enable:        mapDefault(service.Attrs, "SERVICE_ENABLE", "true"),
		LBConfig:      lbc,
	}

	jsonDoc, err := json.Marshal(doc)
	if err != nil {
		log.Fatal("couchbase: failed to marshal document:", err)
	}

	_, err = r.bucket.Upsert(service.ID, jsonDoc, uint32(service.TTL))
	if err != nil {
		log.Fatal("couchbase: failed to Insert document:", err)
	}

	return err
}

func (r *CouchbaseAdapter) Deregister(service *bridge.Service) error {
	var err error
	log.Println("Deregister:", service.ID)
	var tmp []byte

	cas, err := r.bucket.Get(service.ID, &tmp)
	if err != nil {
		log.Fatal("couchbase: Deregister could not get key:", err)
	}

	_, err = r.bucket.Remove(service.ID, cas)
	if err != nil {
		log.Printf("couchbase: failed to deregister service: %v\n", err)
	}
	return err
}

func (r *CouchbaseAdapter) Refresh(service *bridge.Service) error {
	return r.Register(service)
}
