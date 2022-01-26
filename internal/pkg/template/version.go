package template

import (
	"fmt"

	"github.com/Masterminds/semver"

	"github.com/keboola/keboola-as-code/internal/pkg/json"
)

type value = semver.Version

// Version is wrapper around semver.Version - for better error message in UnmarshalJSON.
type Version struct {
	value
}

func NewVersion(str string) (Version, error) {
	v, err := semver.NewVersion(str)
	if err != nil {
		return Version{}, err
	}
	return Version{value: *v}, nil
}

func ZeroVersion() Version {
	v, err := NewVersion(`0.0.1`)
	if err != nil {
		panic(err)
	}
	return v
}

func (v Version) Value() *semver.Version {
	value := v.value
	return &value
}

// IncMajor increments major version, for example 1.2.3 -> 2.0.0.
func (v Version) IncMajor() Version {
	return Version{value: v.Value().IncMajor()}
}

// UnmarshalJSON returns human-readable error message, if semantic version is invalid.
func (v *Version) UnmarshalJSON(b []byte) (err error) {
	var versionStr string
	if err := json.Decode(b, &versionStr); err != nil {
		return err
	}

	value, err := semver.NewVersion(versionStr)
	if err != nil {
		return fmt.Errorf(`invalid semantic version "%s"`, versionStr)
	}
	v.value = *value
	return nil
}