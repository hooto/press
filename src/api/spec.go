// Copyright 2014 lessOS.com. All rights reserved.
//
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package api

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
type NodeModel struct {
	TypeMeta `json:",inline"`
	Metadata ObjectMeta  `json:"metadata,omitempty"`
	SpecID   string      `json:"specID,omitempty"`
	State    int16       `json:"state,omitempty"`
	Title    string      `json:"title,omitempty"`
	Comment  string      `json:"comment,omitempty"`
	Fields   []NodeField `json:"fields,omitempty"`
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
	Type     string     `json:"type,omitempty"`
	Title    string     `json:"title,omitempty"`
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
