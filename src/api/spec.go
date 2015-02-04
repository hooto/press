// Copyright 2014 lessOS.com. All rights reserved.
//
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package api

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Instance struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta `json:"metadata,omitempty"`
	Spec     Spec       `json:"spec,omitempty"`
}

type Spec struct {
	TypeMeta   `json:",inline"`
	Metadata   ObjectMeta  `json:"metadata,omitempty"`
	State      int16       `json:"state,omitempty"`
	Title      string      `json:"title,omitempty"`
	Comment    string      `json:"comment,omitempty"`
	NodeModels []NodeModel `json:"nodeModels,omitempty"`
	TermModels []TermModel `json:"termModels,omitempty"`
	Actions    []Action    `json:"actions,omitempty"`
	Router     Router      `json:"router,omitempty"`
}

type SpecList struct {
	TypeMeta `json:",inline"`
	Items    []Spec `json:"items,omitempty"`
}

//

type FieldModel struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Length    string     `json:"length,omitempty"`
	Extra     []string   `json:"extra,omitempty"`
	Attrs     []KeyValue `json:"attrs,omitempty"`
	IndexType string     `json:"indexType,omitempty"`
	Title     string     `json:"title"`
	Comment   string     `json:"comment,omitempty"`
}

type NodeModel struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta   `json:"metadata,omitempty"`
	SpecID   string       `json:"specID,omitempty"`
	State    int16        `json:"state,omitempty"`
	Title    string       `json:"title,omitempty"`
	Comment  string       `json:"comment,omitempty"`
	Fields   []FieldModel `json:"fields,omitempty"`
	Terms    []TermModel  `json:"terms,omitempty"`
}

type NodeModelList struct {
	TypeMeta `json:",inline"`
	Items    []NodeModel `json:"items,omitempty"`
}

//
const (
	TermTag      = "tag"
	TermTaxonomy = "taxonomy"
)

type TermModel struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta `json:"metadata,omitempty"`
	SpecID   string     `json:"specID,omitempty"`
	State    int16      `json:"state,omitempty"`
	Type     string     `json:"type,omitempty"`
	Title    string     `json:"title,omitempty"`
	Comment  string     `json:"comment,omitempty"`
}

type TermModelList struct {
	TypeMeta `json:",inline"`
	Items    []TermModel `json:"items,omitempty"`
}

//
type Query struct {
	Spec   string   `json:"spec,omitempty"`
	Table  string   `json:"table,omitempty"`
	Fields string   `json:"fields,omitempty"`
	Order  []string `json:"order,omitempty"`
	Limit  int64    `json:"limit,omitempty"`
	Offset int64    `json:"offset,omitempty"`
	Filter []string `json:"filter,omitempty"`
}

//
type Action struct {
	Name  string       `json:"name,omitempty"`
	Datax []ActionData `json:"datax,omitempty"`
}

type ActionData struct {
	Name  string `json:"name,omitempty"`
	Type  string `json:"type,omitempty"`
	Query Query  `json:"query,omitempty"`
}
