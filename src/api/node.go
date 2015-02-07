// Copyright 2014 lessOS.com. All rights reserved.
//
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package api

type Node struct {
	TypeMeta `json:",inline"`
	SpecID   string      `json:"specID,omitempty"`
	Model    *NodeModel  `json:"model,omitempty"`
	ID       string      `json:"id,omitempty"`
	State    int16       `json:"state,omitempty"`
	UserID   string      `json:"userid,omitempty"`
	Title    string      `json:"title,omitempty"`
	Created  string      `json:"created,omitempty"`
	Updated  string      `json:"updated,omitempty"`
	Fields   []NodeField `json:"fields,omitempty"`
	Terms    []NodeTerm  `json:"terms,omitempty"`
}

type NodeList struct {
	TypeMeta `json:",inline"`
	Model    *NodeModel `json:"model,omitempty"`
	Items    []Node     `json:"items,omitempty"`
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

type NodeField struct {
	Name  string     `json:"name"`
	Value string     `json:"value,omitempty"`
	Attrs []KeyValue `json:"attrs,omitempty"`
}

type NodeTerm struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
	Items []Term `json:"items,omitempty"`
}

// var rdoColumnTypes = map[string]string{
// 	"bool":            "bool",
// 	"string":          "varchar(%v)",
// 	"string-text":     "longtext",
// 	"date":            "date",
// 	"datetime":        "datetime",
// 	"int8":            "tinyint",
// 	"int16":           "smallint",
// 	"int32":           "integer",
// 	"int64":           "bigint",
// 	"uint8":           "tinyint unsigned",
// 	"uint16":          "smallint unsigned",
// 	"uint32":          "integer unsigned",
// 	"uint64":          "bigint unsigned",
// 	"float64":         "double precision",
// 	"float64-decimal": "numeric(%v, %v)",
// }
