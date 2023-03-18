// Code generated by ent extension "primarykey", DO NOT EDIT.

package key

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/keboola/go-client/pkg/keboola"
)

type ConfigurationRowKey struct {
	BranchID    keboola.BranchID
	ComponentID keboola.ComponentID
	ConfigID    keboola.ConfigID
	RowID       keboola.RowID
}

// String converts the key to string, it is required for ent validators.
func (v ConfigurationRowKey) String() string {
	var builder strings.Builder
	{
		builder.WriteString("BranchID:")
		builder.WriteString(strconv.Itoa(int(v.BranchID)))
	}
	{
		builder.WriteString("/")
		builder.WriteString("ComponentID:")
		builder.WriteString(string(v.ComponentID))
	}
	{
		builder.WriteString("/")
		builder.WriteString("ConfigID:")
		builder.WriteString(string(v.ConfigID))
	}
	{
		builder.WriteString("/")
		builder.WriteString("RowID:")
		builder.WriteString(string(v.RowID))
	}
	return builder.String()
}

// Value converts Go struct to database string value.
func (v ConfigurationRowKey) Value() (driver.Value, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return v.String(), nil
}

// Value converts string value from database to Go struct.
func (v *ConfigurationRowKey) Scan(in any) error {
	str, ok := in.(string)
	if !ok {
		return fmt.Errorf(`value "%#v" is not string`, in)
	}

	for _, item := range strings.Split(str, "/") {
		key, value, _ := strings.Cut(item, ":")
		switch key {
		case "BranchID":
			valueInt, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			v.BranchID = keboola.BranchID(valueInt)
		case "ComponentID":
			v.ComponentID = keboola.ComponentID(value)
		case "ConfigID":
			v.ConfigID = keboola.ConfigID(value)
		case "RowID":
			v.RowID = keboola.RowID(value)
		default:
			return fmt.Errorf(`unexpected key part "%s"`, key)
		}
	}

	return v.Validate()
}

// Validate check that all parts of the key are not empty.
func (v ConfigurationRowKey) Validate() error {
	if v.BranchID == 0 {
		return fmt.Errorf(`key part "BranchID" is not set`)
	}
	if v.ComponentID == "" {
		return fmt.Errorf(`key part "ComponentID" is not set`)
	}
	if v.ConfigID == "" {
		return fmt.Errorf(`key part "ConfigID" is not set`)
	}
	if v.RowID == "" {
		return fmt.Errorf(`key part "RowID" is not set`)
	}
	return nil
}