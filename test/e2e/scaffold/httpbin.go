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
package scaffold

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/k8s"
	corev1 "k8s.io/api/core/v1"
)

var (
	_httpbinDeployment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpbin-deployment-e2e-test
spec:
  replicas: %d
  selector:
    matchLabels:
      app: httpbin-deployment-e2e-test
  strategy:
    rollingUpdate:
      maxSurge: 50%%
      maxUnavailable: 1
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: httpbin-deployment-e2e-test
    spec:
      terminationGracePeriodSeconds: 0
      containers:
        - livenessProbe:
            failureThreshold: 3
            initialDelaySeconds: 2
            periodSeconds: 5
            successThreshold: 1
            tcpSocket:
              port: 80
            timeoutSeconds: 2
          readinessProbe:
            failureThreshold: 3
            initialDelaySeconds: 2
            periodSeconds: 5
            successThreshold: 1
            tcpSocket:
              port: 80
            timeoutSeconds: 2
          image: "kennethreitz/httpbin"
          imagePullPolicy: IfNotPresent
          name: httpbin-deployment-e2e-test
          ports:
            - containerPort: 80
              name: "http"
              protocol: "TCP"
`
	_httpService = `
apiVersion: v1
kind: Service
metadata:
  name: httpbin-service-e2e-test
spec:
  selector:
    app: httpbin-deployment-e2e-test
  ports:
    - name: http
      port: 80
      protocol: TCP
      targetPort: 80
  type: ClusterIP
`
)

func (s *Scaffold) newHTTPBIN() (*corev1.Service, error) {
	httpbinDeployment := fmt.Sprintf(_httpbinDeployment, 1)
	if err := k8s.KubectlApplyFromStringE(s.t, s.kubectlOptions, httpbinDeployment); err != nil {
		return nil, err
	}
	if err := k8s.KubectlApplyFromStringE(s.t, s.kubectlOptions, _httpService); err != nil {
		return nil, err
	}
	svc, err := k8s.GetServiceE(s.t, s.kubectlOptions, "httpbin-service-e2e-test")
	if err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *Scaffold) ScaleHTTPBIN(num int) error {
	httpbinDeployment := fmt.Sprintf(_httpbinDeployment, num)
	if err := k8s.KubectlApplyFromStringE(s.t, s.kubectlOptions, httpbinDeployment); err != nil {
		return err
	}
	return nil
}
