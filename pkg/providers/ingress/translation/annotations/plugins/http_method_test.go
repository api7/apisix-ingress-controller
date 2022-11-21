// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package plugins

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-ingress-controller/pkg/providers/ingress/translation/annotations"
	apisixv1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

// annotations:
//
//	k8s.apisix.apache.org/allow-http-method: GET,POST
func TestAnnotationsHttpAllowMethod(t *testing.T) {
	anno := map[string]string{
		annotations.AnnotationsHttpAllowMethod: "GET,POST",
	}
	p := NewHttpMethodHandler()
	out, err := p.Handle(annotations.NewExtractor(anno))
	assert.Nil(t, err, "checking given error")
	config := out.(*apisixv1.ResponseRewriteConfig)

	assert.Equal(t, 405, config.StatusCode)
	assert.Equal(t, []apisixv1.Expr{
		{ArrayVal: []apisixv1.Expr{
			{StringVal: "request_method"},
			{StringVal: "!"},
			{StringVal: "in"},
			{ArrayVal: []apisixv1.Expr{
				{StringVal: "GET"},
				{StringVal: "POST"},
			}},
		}},
	}, config.LuaRestyExpr)
}

// annotations:
//
//	k8s.apisix.apache.org/block-http-method: GET,PUT
func TestAnnotationsHttpBlockMethod(t *testing.T) {
	anno := map[string]string{
		annotations.AnnotationsHttpBlockMethod: "GET,PUT",
	}
	p := NewHttpMethodHandler()
	out, err := p.Handle(annotations.NewExtractor(anno))
	assert.Nil(t, err, "checking given error")
	config := out.(*apisixv1.ResponseRewriteConfig)

	assert.Equal(t, 405, config.StatusCode)
	assert.Equal(t, []apisixv1.Expr{
		{ArrayVal: []apisixv1.Expr{
			{StringVal: "request_method"},
			{StringVal: "in"},
			{ArrayVal: []apisixv1.Expr{
				{StringVal: "GET"},
				{StringVal: "PUT"},
			}},
		}},
	}, config.LuaRestyExpr)
}

// annotations:
//
//	k8s.apisix.apache.org/allow-http-method: GET
//	k8s.apisix.apache.org/block-http-method: POST,PUT
//
// Only allow methods would be accepted, block methods would be ignored.
func TestAnnotationsHttpBothMethod(t *testing.T) {
	anno := map[string]string{
		annotations.AnnotationsHttpAllowMethod: "GET",
		annotations.AnnotationsHttpBlockMethod: "POST,PUT",
	}
	p := NewHttpMethodHandler()
	out, err := p.Handle(annotations.NewExtractor(anno))
	assert.Nil(t, err, "checking given error")
	config := out.(*apisixv1.ResponseRewriteConfig)

	assert.Equal(t, 405, config.StatusCode)
	assert.Equal(t, []apisixv1.Expr{
		{ArrayVal: []apisixv1.Expr{
			{StringVal: "request_method"},
			{StringVal: "!"},
			{StringVal: "in"},
			{ArrayVal: []apisixv1.Expr{
				{StringVal: "GET"},
			}},
		}},
	}, config.LuaRestyExpr)
}
