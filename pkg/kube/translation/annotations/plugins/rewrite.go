// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package plugins

import (
	"regexp"

	"github.com/apache/apisix-ingress-controller/pkg/kube/translation/annotations"
	apisixv1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

const (
	_rewriteTarget              = annotations.AnnotationsPrefix + "rewrite-target"
	_rewriteTargetRegex         = annotations.AnnotationsPrefix + "rewrite-target-regex"
	_rewriteTargetRegexTemplate = annotations.AnnotationsPrefix + "rewrite-target-regex-template"
)

type rewrite struct{}

// NewRewriteHandler creates a handler to convert
// annotations about request rewrite control to APISIX proxy-rewrite plugin.
func NewRewriteHandler() PluginHandler {
	return &rewrite{}
}

func (i *rewrite) PluginName() string {
	return "proxy-rewrite"
}

func (i *rewrite) Handle(ing *annotations.Ingress) (interface{}, error) {
	var plugin apisixv1.RewriteConfig
	rewriteTarget := annotations.GetStringAnnotation(_rewriteTarget, ing)
	rewriteTargetRegex := annotations.GetStringAnnotation(_rewriteTargetRegex, ing)
	rewriteTemplate := annotations.GetStringAnnotation(_rewriteTargetRegexTemplate, ing)
	if rewriteTarget != "" || rewriteTargetRegex != "" || rewriteTemplate != "" {
		plugin.RewriteTarget = rewriteTarget
		if rewriteTargetRegex != "" && rewriteTemplate != "" {
			_, err := regexp.Compile(rewriteTargetRegex)
			if err != nil {
				return nil, err
			}
			plugin.RewriteTargetRegex = []string{rewriteTargetRegex, rewriteTemplate}
		}
		return &plugin, nil
	}
	return nil, nil
}
