// Copyright 2014 lessOS.com. All rights reserved.
//
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.
package api

type Router struct {
	Routes []Route `json:"routes,omitempty"`
	// DefaultPagelet string  `json:"defaultPagelet,omitempty"` // e.g. index.tpl
}

type Route struct {
	Path       string            `json:"path,omitempty"` // e.g. /app/:id
	DataAction string            `json:"dataAction,omitempty"`
	Template   string            `json:"template,omitempty"` // e.g. index.tpl
	Params     map[string]string `json:"params,omitempty"`
	Tree       []string          `json:",omitempty"`
}
