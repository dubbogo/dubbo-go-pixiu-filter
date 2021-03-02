/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"sync"
)

var (
	CacheApi = sync.Map{}
	EmptyApi = &API{}
)

// API is api gateway concept, control request from browser、Mobile APP、third party people
type API struct {
	Name          string      `json:"name" yaml:"name"`
	ITypeStr      string      `json:"itype" yaml:"itype"`
	IType         ApiType     `json:"-" yaml:"-"`
	OTypeStr      string      `json:"otype" yaml:"otype"`
	OType         ApiType     `json:"-" yaml:"-"`
	Status        Status      `json:"status" yaml:"status"`
	Metadata      interface{} `json:"metadata" yaml:"metadata"`
	Method        string      `json:"method" yaml:"method"`
	RequestMethod `json:",omitempty" yaml:"-"`
}

// NewApi
func NewApi() *API {
	return &API{}
}

// FindApi find a api, if not exist, return false
func (a *API) FindApi(name string) (*API, bool) {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*API), true
	}

	return nil, false
}

// MatchMethod
func (a *API) MatchMethod(method string) bool {
	i := RequestMethodValue[method]
	if a.RequestMethod == RequestMethod(i) {
		return true
	}

	return false
}

// IsOk api status equals Up
func (a *API) IsOk(name string) bool {
	if v, ok := CacheApi.Load(name); ok {
		return v.(*API).Status == Up
	}

	return false
}

// Offline api offline
func (a *API) Offline(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*API).Status = Down
	}
}

// Online api online
func (a *API) Online(name string) {
	if v, ok := CacheApi.Load(name); ok {
		v.(*API).Status = Up
	}
}
