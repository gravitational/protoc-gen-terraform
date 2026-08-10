package v1

import "github.com/gravitational/protoc-gen-terraform/v5/examples/types"

type DurationType = types.DurationType
type DurationValue = types.DurationValue

type TimeType = types.TimeType
type TimeValue = types.TimeValue

var (
	UseRFC3339Time = types.UseRFC3339Time

	NewDuration     = types.NewDuration
	NullDuration    = types.NullDuration
	UnknownDuration = types.UnknownDuration

	NewTime     = types.NewTime
	NullTime    = types.NullTime
	UnknownTime = types.UnknownTime
)
