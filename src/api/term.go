// Copyright 2014 lessOS.com. All rights reserved.
//
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package api

type Term struct {
	TypeMeta `json:",inline"`
	Model    *TermModel `json:"model,omitempty"`
	ID       uint32     `json:"id,omitempty"`
	PID      uint32     `json:"pid,omitempty"`
	UID      string     `json:"uid,omitempty"`
	State    int16      `json:"state,omitempty"`
	UserID   string     `json:"userid,omitempty"`
	Title    string     `json:"title"`
	Weight   int32      `json:"weight,omitempty"`
	Created  string     `json:"created,omitempty"`
	Updated  string     `json:"updated,omitempty"`
}

type TermList struct {
	TypeMeta `json:",inline"`
	Metadata ListMeta   `json:"metadata,omitempty"`
	Model    *TermModel `json:"model,omitempty"`
	Items    []Term     `json:"items"`
}
