// Copyright 2015 lessOS.com, All rights reserved.
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
	"github.com/lessos/lessgo/types"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Spec struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	SrvName        string           `json:"srvname"`
	Status         int16            `json:"status,omitempty"`
	Title          string           `json:"title"`
	Comment        string           `json:"comment,omitempty"`
	NodeModels     []NodeModel      `json:"nodeModels,omitempty"`
	TermModels     []TermModel      `json:"termModels,omitempty"`
	Actions        []Action         `json:"actions,omitempty"`
	Views          []View           `json:"views,omitempty"`
	Router         Router           `json:"router,omitempty"`
}

type View struct {
	Path string `json:"path"`
}

type ViewList struct {
	types.TypeMeta `json:",inline"`
	Items          []View `json:"items,omitempty"`
}

type SpecList struct {
	types.TypeMeta `json:",inline"`
	Items          []Spec `json:"items,omitempty"`
}

type FieldModel struct {
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Length    string     `json:"length,omitempty"`
	Extra     []string   `json:"extra,omitempty"`
	Attrs     []KeyValue `json:"attrs,omitempty"`
	IndexType int        `json:"indexType,omitempty"`
	Title     string     `json:"title"`
	Comment   string     `json:"comment,omitempty"`
}

type NodeModel struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	ModName        string           `json:"modname,omitempty"`
	Status         int16            `json:"status,omitempty"`
	Title          string           `json:"title,omitempty"`
	Comment        string           `json:"comment,omitempty"`
	Fields         []FieldModel     `json:"fields,omitempty"`
	Terms          []TermModel      `json:"terms,omitempty"`
	Extensions     SpecExtensions   `json:"extensions,omitempty"`
}

type SpecExtensions struct {
	AccessCounter   bool   `json:"access_counter,omitempty"`
	CommentEnable   bool   `json:"comment_enable,omitempty"`
	CommentPerEntry bool   `json:"comment_perentry,omitempty"`
	Permalink       string `json:"permalink,omitempty"`
}

type NodeModelList struct {
	types.TypeMeta `json:",inline"`
	Items          []NodeModel `json:"items,omitempty"`
}

//
const (
	TermTag      = "tag"
	TermTaxonomy = "taxonomy"
)

type TermModel struct {
	types.TypeMeta `json:",inline"`
	Meta           types.ObjectMeta `json:"meta,omitempty"`
	ModName        string           `json:"modname,omitempty"`
	Status         int16            `json:"status,omitempty"`
	Type           string           `json:"type,omitempty"`
	Title          string           `json:"title,omitempty"`
	Comment        string           `json:"comment,omitempty"`
}

type TermModelList struct {
	types.TypeMeta `json:",inline"`
	Items          []TermModel `json:"items,omitempty"`
}

//
type Query struct {
	Spec   string   `json:"spec,omitempty"`
	Table  string   `json:"table,omitempty"`
	Fields string   `json:"fields,omitempty"`
	Order  string   `json:"order,omitempty"`
	Limit  int64    `json:"limit,omitempty"`
	Offset int64    `json:"offset,omitempty"`
	Filter []string `json:"filter,omitempty"`
}

//
type Action struct {
	types.TypeMeta `json:",inline"`
	Name           string       `json:"name"`
	ModName        string       `json:"modname,omitempty"`
	Datax          []ActionData `json:"datax,omitempty"`
}

type ActionData struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Pager    bool   `json:"pager,omitempty"`
	Query    Query  `json:"query,omitempty"`
	CacheTTL int64  `json:"cache_ttl,omitempty"` // cache time to live in milliseconds
}
