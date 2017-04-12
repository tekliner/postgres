package api

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

const (
	ResourceKindDeletedDatabase = "DeletedDatabase"
	ResourceNameDeletedDatabase = "deleted-database"
	ResourceTypeDeletedDatabase = "deleteddatabases"
)

type DeletedDatabase struct {
	unversioned.TypeMeta `json:",inline,omitempty"`
	api.ObjectMeta       `json:"metadata,omitempty"`
	Spec                 DeletedDatabaseSpec   `json:"spec,omitempty"`
	Status               DeletedDatabaseStatus `json:"status,omitempty"`
}

type DeletedDatabaseSpec struct {
	// Database authentication secret
	// +optional
	AuthSecret *api.SecretVolumeSource `json:"authSecret,omitempty"`
	// If true, invoke destroy operation
	// +optional
	Destroy bool `json:"destroy,omitempty"`
}

type DeletedDatabasePhase string

const (
	// used for Databases that are deleted
	PhaseDatabaseDeleted DeletedDatabasePhase = "Deleted"
	// used for Databases that are currently deleting
	PhaseDatabaseDeleting DeletedDatabasePhase = "Deleting"
	// used for Databases that are destroyed
	PhaseDatabaseDestroyed DeletedDatabasePhase = "Destroyed"
	// used for Databases that are currently destroying
	PhaseDatabaseDestroying DeletedDatabasePhase = "Destroying"
)

type DeletedDatabaseStatus struct {
	Created   *unversioned.Time    `json:"created,omitempty"`
	Deleted   *unversioned.Time    `json:"deleted,omitempty"`
	Destroyed *unversioned.Time    `json:"destroyed,omitempty"`
	Phase     DeletedDatabasePhase `json:"phase,omitempty"`
	Reason    string               `json:"reason,omitempty"`
}

type DeletedDatabaseList struct {
	unversioned.TypeMeta `json:",inline"`
	unversioned.ListMeta `json:"metadata,omitempty"`
	// Items is a list of DeletedDatabase TPR objects
	Items []DeletedDatabase `json:"items,omitempty"`
}