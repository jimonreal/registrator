package couchbase

import "github.com/gliderlabs/registrator/bridge"

type Document struct {
	ContainerId   string
	ContainerType string
	//	Location      ContainerLocation
	Location bridge.ServicePort
	MetaData map[string]string
	Created  int64
	Updated  int64
	LBConfig LoadBalancerConfig
	Enable   string
}

//type ContainerLocation struct {
//	HostIp       string
//	PortsMapping map[string]ServicePort
//}

//type CBMetaData struct {
//	Data map[string]string
//}

type LoadBalancerConfig struct {
	Enable  string
	Mode    string
	Balance string
	Param   string
}
