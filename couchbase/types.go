package couchbase

type Document struct {
	ContainerId   string
	ContainerType string
	Location      ContainerLocation
	MetaData      CBMetaData
	Created       int64
	Updated       int64
}

type ContainerLocation struct {
	HostIp       string
	PortsMapping map[string]ContainerPorts
}

type ContainerPorts struct {
	HostPort     string
	HostPortType string
}

type CBMetaData struct {
	Data map[string]string
}
