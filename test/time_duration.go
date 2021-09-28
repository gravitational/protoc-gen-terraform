package test

import (
	"context"
	fmt "fmt"
	time "time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	tftypes "github.com/hashicorp/terraform-plugin-go/tftypes"
)

const (
	timeFormat    = time.RFC3339
	timeThreshold = time.Nanosecond
)

// TimeType represents time.Time Terraform type which is stored in RFC3339 format, nanoseconds truncated
type TimeType struct {
	attr.Type
}

// ApplyTerraform5AttributePathStep is not implemented for TimeType
func (t TimeType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// String returns string representation of TimeType
func (t TimeType) String() string {
	return "TimeType"
}

// Equal returns type equality
func (t TimeType) Equal(o attr.Type) bool {
	other, ok := o.(TimeType)
	if !ok {
		return false
	}
	return t == other
}

// TerraformType returns type which is used in Terraform status (time is stored as string)
func (t TimeType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// ValueFromTerraform decodes terraform value and returns it as TimeType
func (t TimeType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return TimeValue{Unknown: true}, nil
	}
	if in.IsNull() {
		return TimeValue{Null: true}, nil
	}
	var raw string
	err := in.As(&raw)
	if err != nil {
		return nil, err
	}

	// Error is deliberately silenced here. If a value is corrupted, this would be caught in Validate() method which
	// for some reason is called after ValueFromTerraform().
	current, err := time.Parse(timeFormat, raw)
	if err != nil {
		return nil, err
	}

	return TimeValue{Value: current}, nil
}

// TimeValue represents Terraform value of type TimeType
type TimeValue struct {
	// Unknown will be true if the value is not yet known.
	Unknown bool
	// Null will be true if the value was not set, or was explicitly set to
	// null.
	Null bool
	// Value contains the set value, as long as Unknown and Null are both
	// false.
	Value time.Time
}

// Type returns value type
func (t TimeValue) Type(_ context.Context) attr.Type {
	return TimeType{}
}

// ToTerraformValue returns the data contained in the *String as a string. If
// Unknown is true, it returns a tftypes.UnknownValue. If Null is true, it
// returns nil.
func (t TimeValue) ToTerraformValue(_ context.Context) (interface{}, error) {
	if t.Null {
		return nil, nil
	}
	if t.Unknown {
		return tftypes.UnknownValue, nil
	}
	return t.Value.Truncate(timeThreshold).Format(timeFormat), nil
}

// Equal returns true if `other` is a *String and has the same value as `s`.
func (t TimeValue) Equal(other attr.Value) bool {
	o, ok := other.(TimeValue)
	if !ok {
		return false
	}
	if t.Unknown != o.Unknown {
		return false
	}
	if t.Null != o.Null {
		return false
	}
	return t.Value == o.Value
}

// DurationType represents time.Time Terraform type which is stored in RFC3339 format, nanoseconds truncated
type DurationType struct {
	attr.Type
}

// ApplyTerraform5AttributePathStep is not implemented for TimeType
func (t DurationType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// String returns string representation of TimeType
func (t DurationType) String() string {
	return "DurationType"
}

// Equal returns type equality
func (t DurationType) Equal(o attr.Type) bool {
	other, ok := o.(DurationType)
	if !ok {
		return false
	}
	return t == other
}

// DurationType returns type which is used in Terraform status (time is stored as string)
func (t DurationType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// ValueFromTerraform decodes terraform value and returns it as TimeType
func (t DurationType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if !in.IsKnown() {
		return DurationValue{Unknown: true}, nil
	}
	if in.IsNull() {
		return DurationValue{Null: true}, nil
	}
	var raw string
	err := in.As(&raw)
	if err != nil {
		return nil, err
	}

	// Error is deliberately silenced here. If a value is corrupted, this would be caught in Validate() method which
	// for some reason is called after ValueFromTerraform().
	current, err := time.ParseDuration(raw)
	if err != nil {
		return nil, err
	}

	return DurationValue{Value: current}, nil
}

// DurationValue represents Terraform value of type TimeType
type DurationValue struct {
	// Unknown will be true if the value is not yet known.
	Unknown bool
	// Null will be true if the value was not set, or was explicitly set to
	// null.
	Null bool
	// Value contains the set value, as long as Unknown and Null are both
	// false.
	Value time.Duration
}

// Type returns value type
func (t DurationValue) Type(_ context.Context) attr.Type {
	return TimeType{}
}

// ToTerraformValue returns the data contained in the *String as a string. If
// Unknown is true, it returns a tftypes.UnknownValue. If Null is true, it
// returns nil.
func (t DurationValue) ToTerraformValue(_ context.Context) (interface{}, error) {
	if t.Null {
		return nil, nil
	}
	if t.Unknown {
		return tftypes.UnknownValue, nil
	}
	return t.Value.String(), nil
}

// Equal returns true if `other` is a *String and has the same value as `s`.
func (t DurationValue) Equal(other attr.Value) bool {
	o, ok := other.(DurationValue)
	if !ok {
		return false
	}
	if t.Unknown != o.Unknown {
		return false
	}
	if t.Null != o.Null {
		return false
	}
	return t.Value == o.Value
}
