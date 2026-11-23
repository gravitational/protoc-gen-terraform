/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"testing"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stretchr/testify/require"
)

// issue #59 schema:
//
//	message Objects {
//	  Nested embedded = 13 [(gogoproto.nullable) = false, (gogoproto.embed) = true];
//	  Leaf embedded_nullable = 14 [(gogoproto.nullable) = true, (gogoproto.embed) = true];
//	}
//	message Nested { Leaf leaf = 1; }
//	message Leaf { string value = 1; }
//
// gogo lays that out as:
//
//	type Objects struct {
//	  Nested
//	  *Leaf
//	}
//
// obj.Leaf is the embedded *Leaf (depth 1). Nested.Leaf is only reached as
// obj.Nested.Leaf. Writing the nested attribute through obj.Leaf stores the
// nullable embed's value in .leaf.value.
func issueObjectsMessage() *Message {
	leafValue := &Field{
		Name:          "Value",
		NameSnake:     "value",
		Kind:          PrimitiveKind,
		Path:          "Leaf.value",
		TerraformType: stringType,
		ProtobufType: ProtobufType{
			GoType:             "string",
			GoElemType:         "string",
			GoElemTypeIndirect: "string",
		},
	}
	nestedLeaf := &Field{
		Name:          "Leaf",
		NameSnake:     "leaf",
		Kind:          ObjectKind,
		IsNullable:    true,
		Path:          "Nested.leaf",
		TerraformType: objectType,
		ProtobufType: ProtobufType{
			GoType:             "*Leaf",
			GoElemType:         "*Leaf",
			GoElemTypeIndirect: "Leaf",
		},
		Message: &Message{
			Name:   "Leaf",
			GoType: "Leaf",
			Fields: []*Field{leafValue},
		},
	}
	topValue := &Field{
		Name:          "Value",
		NameSnake:     "value",
		Kind:          PrimitiveKind,
		Path:          "Objects.value",
		TerraformType: stringType,
		ProtobufType: ProtobufType{
			GoType:             "string",
			GoElemType:         "string",
			GoElemTypeIndirect: "string",
		},
	}

	// Non-nullable embed of Nested, nullable embed of *Leaf. Same shape gogo emits.
	fields := flattenEmbeddedFields([]*Field{nestedLeaf}, "Nested", false)
	fields = append(fields, flattenEmbeddedFields([]*Field{topValue}, "*pkg.Leaf", true)...)

	return &Message{
		Name:   "Objects",
		GoType: "Objects",
		IsRoot: true,
		Fields: fields,
	}
}

func generateCopyMethods(t *testing.T, msg *Message) (string, string) {
	t.Helper()
	gen := generator.New()
	imps := NewImports(generator.NewPluginImports(gen), nil)

	var from, to bytes.Buffer
	_, err := NewMessageCopyFromGenerator(msg, &imps).Generate(&from)
	require.NoError(t, err)
	copier := NewMessageCopyToGenerator(msg, &imps)
	_, err = copier.Generate(&to)
	require.NoError(t, err)
	_, err = copier.GeneratePreserveUnknown(&to)
	require.NoError(t, err)
	return from.String(), to.String()
}

func TestEmbeddedAttributeCollisionUsesQualifiedGoPath(t *testing.T) {
	original := &Field{Name: "Leaf"}
	promoted := flattenEmbeddedFields([]*Field{original}, "github.com/example/types.Nested", false)
	require.Equal(t, "Nested.Leaf", promoted[0].Name)
	require.Equal(t, "Leaf", original.Name)

	nullable := flattenEmbeddedFields([]*Field{{Name: "Value"}}, "*pkg.Leaf", true)
	require.Equal(t, "Leaf.Value", nullable[0].Name)
	require.True(t, nullable[0].ParentIsOptionalEmbed)
	require.Equal(t, "Leaf", nullable[0].ParentIsOptionalEmbedFieldName)
	require.Equal(t, "pkg.Leaf", nullable[0].ParentIsOptionalEmbedFullType)
	require.Nil(t, flattenEmbeddedFields(nil, "Nested", false))

	// A nullable embed nested inside a non-nullable embed keeps one pointer:
	// the nil-check and the write both go through Mid.Leaf.
	nested := flattenEmbeddedFields(nullable, "Mid", false)
	require.Equal(t, "Mid.Leaf.Value", nested[0].Name)
	require.Equal(t, "Mid.Leaf", nested[0].ParentIsOptionalEmbedFieldName)
	require.True(t, nested[0].ParentIsOptionalEmbed)
	require.Equal(t, "Leaf.Value", nullable[0].Name)

	from, to := generateCopyMethods(t, issueObjectsMessage())

	// Nested.leaf is addressed through the embedded Nested value. The inner
	// obj.Value write is on that rebound pointer, not on the parent.
	require.Contains(t, from, "obj.Nested.Leaf = &Leaf{}")
	require.Contains(t, from, "obj := obj.Nested.Leaf")
	require.Contains(t, to, "obj.Nested.Leaf")
	// Nullable embed's value is stored on that embed.
	require.Contains(t, from, "obj.Leaf.Value = t")
	require.Contains(t, to, "obj.Leaf.Value")
}

// TestEmbedPromotionShadowing locks the Go rule behind issue #59.
// The shallower *Leaf embed wins obj.Leaf; only obj.Nested.Leaf is the nested attribute.
func TestEmbedPromotionShadowing(t *testing.T) {
	type Leaf struct{ Value string }
	type Nested struct{ Leaf *Leaf }
	type Objects struct {
		Nested
		*Leaf
	}

	obj := &Objects{
		Nested: Nested{Leaf: &Leaf{Value: "embedded-leaf"}},
		Leaf:   &Leaf{Value: "embedded-nullable-value"},
	}

	require.Equal(t, "embedded-nullable-value", obj.Leaf.Value)
	require.Equal(t, "embedded-leaf", obj.Nested.Leaf.Value)

	// What the generator used to emit for the nested attribute.
	obj.Leaf = &Leaf{}
	obj.Leaf.Value = "embedded-leaf"
	obj.Leaf.Value = "embedded-nullable-value"
	require.Equal(t, "embedded-nullable-value", obj.Leaf.Value)
	require.Equal(t, "embedded-leaf", obj.Nested.Leaf.Value)

	// Qualified selectors keep the two attributes apart.
	obj.Nested.Leaf = &Leaf{Value: "embedded-leaf"}
	obj.Leaf = &Leaf{Value: "embedded-nullable-value"}
	require.Equal(t, "embedded-leaf", obj.Nested.Leaf.Value)
	require.Equal(t, "embedded-nullable-value", obj.Leaf.Value)
}
