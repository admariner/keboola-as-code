package model

import (
	"time"

	"github.com/c2h5oh/datasize"
)

type ImportCondition struct {
	Count int               `json:"count" validate:"min=1,max=10000000"`
	Size  datasize.ByteSize `json:"size" validate:"min=100,max=50000000"`
	Time  time.Duration     `json:"time" validate:"min=30s,max=24h"`
}

type Export struct {
	ID               string            `json:"exportId" validate:"required,min=1,max=48"`
	Name             string            `json:"name" validate:"required,min=1,max=40"`
	ImportConditions []ImportCondition `json:"importConditions" validate:"required"`
}

type MappedColumns []any

type Mapping struct {
	RevisionID  int           `json:"revisionId" validate:"required"`
	TableID     TableID       `json:"tableId" validate:"required,min=1,max=198"`
	Incremental bool          `json:"incremental" validate:"required"`
	Columns     MappedColumns `json:"columns" validate:"required,min=1,max=50"`
}

type Receiver struct {
	ID        string `json:"receiverId" validate:"required,min=1,max=48"`
	ProjectID int    `json:"projectId" validate:"required"`
	Name      string `json:"name" validate:"required,min=1,max=40"`
	Secret    string `json:"secret" validate:"required,len=48"`
}