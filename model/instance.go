package model

// Instance represents app instance object
type Instance struct {
	Name   string `json:"instanceName,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status,omitempty"`
}
