package couchbase

type Document struct {
	ContainerId   string
	ContainerType string
	Location      ContainerLocation
	MetaData      CBMetaData
	_created      int64
	_updated      int64
}

type ContainerLocation struct {
	HostIp       string
	PortsMapping []ContainerPorts
}

type ContainerPorts struct {
	GuestPort    string
	HostPort     string
	HostPortType string
}

type CBMetaData struct {
	data [][]string
}
