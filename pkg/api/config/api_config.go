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

package config

import (
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
)

// HTTPVerb defines the restful api http verb
type HTTPVerb string

const (
	// MethodAny any method
	MethodAny HTTPVerb = "ANY"
	// MethodGet get
	MethodGet HTTPVerb = "GET"
	// MethodHead head
	MethodHead HTTPVerb = "HEAD"
	// MethodPost post
	MethodPost HTTPVerb = "POST"
	// MethodPut put
	MethodPut HTTPVerb = "PUT"
	// MethodPatch patch
	MethodPatch HTTPVerb = "PATCH" // RFC 5789
	// MethodDelete delete
	MethodDelete HTTPVerb = "DELETE"
	// MethodOptions options
	MethodOptions HTTPVerb = "OPTIONS"
)

// RequestType describes the type of the request. could be DUBBO/HTTP and others that we might implement in the future
type RequestType string

const (
	// DubboRequest represents the dubbo request
	DubboRequest RequestType = "dubbo"
	// HTTPRequest represents the http request
	HTTPRequest RequestType = "http"
)

// APIConfig defines the data structure of the api gateway configuration
type APIConfig struct {
	Name           string           `json:"name" yaml:"name"`
	Description    string           `json:"description" yaml:"description"`
	Resources      []Resource       `json:"resources" yaml:"resources"`
	Definitions    []Definition     `json:"definitions" yaml:"definitions"`
	PluginFilePath string           `json:"pluginFilePath" yaml:"pluginFilePath"`
	PluginsGroup   []PluginsGroup   `json:"pluginsGroup" yaml:"pluginsGroup"`
	RateLimit      ratelimit.Config `json:"rateLimit" yaml:"rateLimit"`
}

type Plugin struct {
	ID                 int64  `json:"id,inline,omitempty" yaml:"id,omitempty"`
	Name               string `json:"name" yaml:"name"`
	Version            string `json:"version" yaml:"version"`
	Priority           int    `json:"priority" yaml:"priority"`
	ExternalLookupName string `json:"externalLookupName" yaml:"externalLookupName"`
}

// PluginsGroup defines the plugins group info
type PluginsGroup struct {
	ID        int64    `json:"id,omitempty" yaml:"id,omitempty"`
	GroupName string   `json:"groupName" yaml:"groupName"`
	Plugins   []Plugin `json:"plugins" yaml:"plugins"`
}

//PluginsConfig defines the pre & post plugins
type PluginsConfig struct {
	PrePlugins  PluginsInUse `json:"pre" yaml:"pre"`
	PostPlugins PluginsInUse `json:"post" yaml:"post"`
}

type PluginsInUse struct {
	GroupNames  []string `json:"groupNames" yaml:"groupNames"`
	PluginNames []string `json:"pluginNames" yaml:"pluginNames"`
}

// Resource defines the API path
type Resource struct {
	ID          int               `json:"id,omitempty" yaml:"id,omitempty"`
	Type        string            `json:"type" yaml:"type"` // Restful, Dubbo
	Path        string            `json:"path" yaml:"path"`
	Timeout     time.Duration     `json:"timeout" yaml:"timeout"`
	Description string            `json:"description" yaml:"description"`
	Filters     []string          `json:"filters" yaml:"filters"`
	Plugins     PluginsConfig     `json:"plugins" yaml:"plugins"`
	Methods     []Method          `json:"methods" yaml:"methods"`
	Resources   []Resource        `json:"resources,omitempty" yaml:"resources,omitempty"`
	Headers     map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

const DefaultTimeoutStr = "1s"

// UnmarshalYAML Resource custom UnmarshalYAML
func (r *Resource) UnmarshalYAML(unmarshal func(interface{}) error) error {
	s := &struct {
		Timeout string `yaml:"timeout"`
	}{}
	type Alias Resource
	alias := (*Alias)(r)
	if err := unmarshal(alias); err != nil {
		return err
	}
	if err := unmarshal(s); err != nil {
		return err
	}
	// if timeout is empty must set a default value. if "" used to time.ParseDuration will err.
	if s.Timeout == "" {
		s.Timeout = DefaultTimeoutStr
	}
	d, err := time.ParseDuration(s.Timeout)
	if err != nil {
		return err
	}

	r.Timeout = d

	return nil
}

// Method defines the method of the api
type Method struct {
	ID                 int           `json:"id,omitempty" yaml:"id,omitempty"`
	ResourcePath       string        `json:"resourcePath" yaml:"resourcePath"`
	Enable             bool          `json:"enable" yaml:"enable"` // true means the method is up and false means method is down
	Timeout            time.Duration `json:"timeout" yaml:"timeout"`
	Mock               bool          `json:"mock" yaml:"mock"`
	Filters            []string      `json:"filters" yaml:"filters"`
	HTTPVerb           `json:"httpVerb" yaml:"httpVerb"`
	InboundRequest     `json:"inboundRequest" yaml:"inboundRequest"`
	IntegrationRequest `json:"integrationRequest" yaml:"integrationRequest"`
}

// UnmarshalYAML method custom UnmarshalYAML
func (m *Method) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type Alias Method
	alias := (*Alias)(m)
	if err := unmarshal(alias); err != nil {
		return err
	}
	s := &struct {
		Timeout string `yaml:"timeout"`
	}{}
	if err := unmarshal(s); err != nil {
		return err
	}
	// if timeout is empty must set a default value. if "" used to time.ParseDuration will err.
	if s.Timeout == "" {
		s.Timeout = DefaultTimeoutStr
	}
	d, err := time.ParseDuration(s.Timeout)
	if err != nil {
		return err
	}
	m.Timeout = d
	return nil
}

// InboundRequest defines the details of the inbound
type InboundRequest struct {
	RequestType  `json:"requestType" yaml:"requestType"` //http, TO-DO: dubbo
	Headers      []Params                                `json:"headers" yaml:"headers"`
	QueryStrings []Params                                `json:"queryStrings" yaml:"queryStrings"`
	RequestBody  []BodyDefinition                        `json:"requestBody" yaml:"requestBody"`
}

// Params defines the simple parameter definition
type Params struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Required bool   `json:"required" yaml:"required"`
}

// BodyDefinition connects the request body to the definitions
type BodyDefinition struct {
	DefinitionName string `json:"definitionName" yaml:"definitionName"`
}

// IntegrationRequest defines the backend request format and target
type IntegrationRequest struct {
	RequestType        `json:"requestType" yaml:"requestType"` // dubbo, TO-DO: http
	DubboBackendConfig `json:"dubboBackendConfig,inline,omitempty" yaml:"dubboBackendConfig,inline,omitempty"`
	HTTPBackendConfig  `json:"httpBackendConfig,inline,omitempty" yaml:"httpBackendConfig,inline,omitempty"`
	MappingParams      []MappingParam `json:"mappingParams,omitempty" yaml:"mappingParams,omitempty"`
}

// MappingParam defines the mapping rules of headers and queryStrings
type MappingParam struct {
	Name    string `json:"name,omitempty" yaml:"name"`
	MapTo   string `json:"mapTo,omitempty" yaml:"mapTo"`
	MapType string `json:"mapType,omitempty" yaml:"mapType"`
}

// DubboBackendConfig defines the basic dubbo backend config
type DubboBackendConfig struct {
	ClusterName     string `yaml:"clusterName" json:"clusterName"`
	ApplicationName string `yaml:"applicationName" json:"applicationName"`
	Protocol        string `yaml:"protocol" json:"protocol,omitempty" default:"dubbo"`
	Group           string `yaml:"group" json:"group"`
	Version         string `yaml:"version" json:"version"`
	Interface       string `yaml:"interface" json:"interface"`
	Method          string `yaml:"method" json:"method"`
	Retries         string `yaml:"retries" json:"retries,omitempty"`
}

// HTTPBackendConfig defines the basic dubbo backend config
type HTTPBackendConfig struct {
	URL string `yaml:"url" json:"url,omitempty"`
	// downstream host.
	Host string `yaml:"host" json:"host,omitempty"`
	// path to replace.
	Path string `yaml:"path" json:"path,omitempty"`
	// http protocol, http or https.
	Schema string `yaml:"schema" json:"scheme,omitempty"`
}

// Definition defines the complex json request body
type Definition struct {
	Name   string `json:"name" yaml:"name"`
	Schema string `json:"schema" yaml:"schema"` // use json schema
}
