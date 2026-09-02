package main

import "github.com/gravitational/trace"

type terraformAttributeMeta struct {
	terraformType             string
	attributeType             string
	validatorPlanModifierBase string
	useStateForUnknownPackage string
}

var terraformAttributeMetas = []terraformAttributeMeta{
	{
		terraformType:             Types + ".StringType",
		attributeType:             "StringAttribute",
		validatorPlanModifierBase: ".String",
		useStateForUnknownPackage: ResourceSchema + "/stringplanmodifier",
	},
	{
		terraformType:             Types + ".BoolType",
		attributeType:             "BoolAttribute",
		validatorPlanModifierBase: ".Bool",
		useStateForUnknownPackage: ResourceSchema + "/boolplanmodifier",
	},
	{
		terraformType:             Types + ".Int64Type",
		attributeType:             "Int64Attribute",
		validatorPlanModifierBase: ".Int64",
		useStateForUnknownPackage: ResourceSchema + "/int64planmodifier",
	},
	{
		terraformType:             Types + ".Float64Type",
		attributeType:             "Float64Attribute",
		validatorPlanModifierBase: ".Float64",
		useStateForUnknownPackage: ResourceSchema + "/float64planmodifier",
	},
	{
		terraformType:             Types + ".ListType",
		attributeType:             "ListAttribute",
		validatorPlanModifierBase: ".List",
		useStateForUnknownPackage: ResourceSchema + "/listplanmodifier",
	},
	{
		terraformType:             Types + ".MapType",
		attributeType:             "MapAttribute",
		validatorPlanModifierBase: ".Map",
		useStateForUnknownPackage: ResourceSchema + "/mapplanmodifier",
	},
	{
		terraformType:             Types + ".ObjectType",
		attributeType:             "ObjectAttribute",
		validatorPlanModifierBase: ".Object",
		useStateForUnknownPackage: ResourceSchema + "/objectplanmodifier",
	},
}

var (
	terraformAttributeMetaByType      = map[string]terraformAttributeMeta{}
	terraformAttributeMetaByAttribute = map[string]terraformAttributeMeta{}
)

func init() {
	for _, meta := range terraformAttributeMetas {
		terraformAttributeMetaByType[meta.terraformType] = meta
		terraformAttributeMetaByAttribute[meta.attributeType] = meta
	}
}

func attributeTypeForTerraformType(t string) (string, error) {
	meta, ok := terraformAttributeMetaByType[t]
	if !ok {
		return "", trace.BadParameter("unexpected type %q", t)
	}

	return meta.attributeType, nil
}

func baseTypeForAttributeType(t string) (string, error) {
	meta, ok := terraformAttributeMetaByAttribute[t]
	if !ok {
		return "", trace.BadParameter("unexpected attribute type %q", t)
	}

	return meta.validatorPlanModifierBase, nil
}

func useStateForUnknownForType(t string) (string, error) {
	meta, ok := terraformAttributeMetaByAttribute[t]
	if !ok {
		return "", trace.BadParameter("unexpected attribute type %q", t)
	}

	return meta.useStateForUnknownPackage + ".UseStateForUnknown()", nil
}
