// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.
//
package scaffold

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/k8s"
)

var (
	_apisixUpstreamsWithMTLSTemplate = `
apiVersion: apisix.apache.org/v1
kind: ApisixUpstream
metadata:
  name: %s
spec:
  scheme: https
  tlsSecret:
    name: %s
    namespace: %s
`
)

// NewApisixUpstreamsWithMTLS new a ApisixUpstreams CRD
func (s *Scaffold) NewApisixUpstreamsWithMTLS(name, secretName string) error {
	tls := fmt.Sprintf(_apisixUpstreamsWithMTLSTemplate, name, secretName, s.Namespace())
	if err := k8s.KubectlApplyFromStringE(s.t, s.kubectlOptions, tls); err != nil {
		return err
	}
	return nil
}
