package model

import (
	"time"

	"github.com/keboola/go-client/pkg/keboola"

	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/key"
	"github.com/keboola/keboola-as-code/internal/pkg/service/buffer/store/slicestate"
	"github.com/keboola/keboola-as-code/internal/pkg/service/common/utctime"
)

const (
	SliceFilenameDateFormat = "20060102150405"
)

// Slice represent a file slice with records.
// A copy of the mapping is stored for retrieval optimization.
// A change in the mapping causes a new file and slice to be created so the mapping is immutable.
type Slice struct {
	key.SliceKey
	State           slicestate.State               `json:"state" validate:"required,oneof=active/opened/writing active/opened/closing active/closed/uploading active/closed/uploaded active/closed/failed archived/successful/imported"`
	Mapping         Mapping                        `json:"mapping" validate:"required,dive"`
	StorageResource *keboola.FileUploadCredentials `json:"storageResource" validate:"required"`
	Number          int                            `json:"sliceNumber" validate:"required"`
	ClosingAt       *utctime.UTCTime               `json:"closingAt,omitempty"`
	UploadingAt     *utctime.UTCTime               `json:"uploadingAt,omitempty"`
	UploadedAt      *utctime.UTCTime               `json:"uploadedAt,omitempty"`
	FailedAt        *utctime.UTCTime               `json:"failedAt,omitempty"`
	ImportedAt      *utctime.UTCTime               `json:"importedAt,omitempty"`
	LastError       string                         `json:"lastError,omitempty"`
	RetryAttempt    int                            `json:"retryAttempt,omitempty"`
	RetryAfter      *utctime.UTCTime               `json:"retryAfter,omitempty"`
	IsEmpty         bool                           `json:"isEmpty,omitempty"`
	// IDRange is assigned during the "slice close" operation, it defines the assigned auto-increment value.
	IDRange *SliceIDRange `json:"idRange,omitempty"`
}

type SliceIDRange struct {
	Start uint64 `json:"start" validate:"required"`
	Count uint64 `json:"count" validate:"required"`
}

func NewSlice(fileKey key.FileKey, now time.Time, mapping Mapping, number int, resource *keboola.FileUploadCredentials) Slice {
	return Slice{
		SliceKey:        key.SliceKey{FileKey: fileKey, SliceID: key.SliceID(now)},
		State:           slicestate.Writing,
		Mapping:         mapping,
		StorageResource: resource,
		Number:          number,
	}
}

func (v Slice) Filename() string {
	return v.OpenedAt().Format(SliceFilenameDateFormat) + ".gz"
}
