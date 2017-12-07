// Copyright 2015 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"regexp"

	"github.com/lessos/lessgo/types"
)

type Node struct {
	types.TypeMeta     `json:",inline"`
	SelfLink           string       `json:"self_link,omitempty"`
	Model              *NodeModel   `json:"model,omitempty"`
	ID                 string       `json:"id,omitempty"`
	PID                string       `json:"pid,omitempty"`
	Status             int16        `json:"status,omitempty"`
	UserID             string       `json:"userid,omitempty"`
	Title              string       `json:"title,omitempty"`
	Created            string       `json:"created,omitempty"`
	Updated            string       `json:"updated,omitempty"`
	Fields             []*NodeField `json:"fields,omitempty"`
	Terms              []NodeTerm   `json:"terms,omitempty"`
	ExtAccessCounter   uint32       `json:"ext_access_counter,omitempty"`
	ExtCommentEnable   bool         `json:"ext_comment_enable,omitempty"`
	ExtCommentPerEntry bool         `json:"ext_comment_perentry,omitempty"`
	ExtPermalinkName   string       `json:"ext_permalink_name,omitempty"`
	ExtNodeRefer       string       `json:"ext_node_refer,omitempty"`
}

var (
	NodeExtNodeReferReg = regexp.MustCompile("^[0-9a-f]{12,16}$")
)

type NodeList struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ListMeta `json:"meta,omitempty"`
	Model          *NodeModel     `json:"model,omitempty"`
	Items          []Node         `json:"items,omitempty"`
}

type NodeFieldType string

const (
	NodeFieldBool     NodeFieldType = "bool"
	NodeFieldString   NodeFieldType = "string"
	NodeFieldText     NodeFieldType = "text"
	NodeFieldDate     NodeFieldType = "date"
	NodeFieldDateTime NodeFieldType = "datetime"
	NodeFieldInt8     NodeFieldType = "int8"
	NodeFieldInt16    NodeFieldType = "int16"
	NodeFieldInt32    NodeFieldType = "int32"
	NodeFieldInt64    NodeFieldType = "int64"
	NodeFieldUint8    NodeFieldType = "uint8"
	NodeFieldUint16   NodeFieldType = "uint16"
	NodeFieldUint32   NodeFieldType = "uint32"
	NodeFieldUint64   NodeFieldType = "uint64"
	NodeFieldFloat    NodeFieldType = "float"
	NodeFieldDecimal  NodeFieldType = "decimal"
)

var (
	NodeFieldTypes = []string{
		"bool",
		"string",
		"text",
		"date",
		"datetime",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float",
		"decimal",
	}
)

type NodeField struct {
	Name        string            `json:"name"`
	Value       string            `json:"value,omitempty"`
	ValueCaches map[string]string `json:"value_caches,omitempty"`
	Attrs       []KeyValue        `json:"attrs,omitempty"`
}

type NodeTerm struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Type  string `json:"type,omitempty"`
	Items []Term `json:"items,omitempty"`
}
