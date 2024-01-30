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
package annotations

import (
	"fmt"
	"net/http"
	"os"
	"time"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"

	"github.com/apache/apisix-ingress-controller/test/e2e/scaffold"
)

var _ = ginkgo.Describe("suite-annotations: allow-http-methods annotations", func() {
	s := scaffold.NewDefaultScaffold()

	ginkgo.It("enable in ingress networking/v1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/http-allow-methods: POST,PUT
  name: ingress-v1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /*
        pathType: Prefix
        backend:
          service:
            name: %s
            port:
              number: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		respGet := s.NewAPISIXClient().GET("/get").WithHeader("Host", "httpbin.org").Expect()
		respGet.Status(http.StatusMethodNotAllowed)

		respPost := s.NewAPISIXClient().POST("/post").WithHeader("Host", "httpbin.org").Expect()
		respPost.Status(http.StatusOK)

		respPut := s.NewAPISIXClient().PUT("/put").WithHeader("Host", "httpbin.org").Expect()
		respPut.Status(http.StatusOK)
	})

	ginkgo.It("enable in ingress networking/v1beta1", func() {
		if os.Getenv("K8s_Version") == "v1.24.0" {
			return
		}
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/http-allow-methods: POST,PUT
  name: ingress-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /*
        pathType: Prefix
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		respGet := s.NewAPISIXClient().GET("/get").WithHeader("Host", "httpbin.org").Expect()
		respGet.Status(http.StatusMethodNotAllowed)

		respPost := s.NewAPISIXClient().POST("/post").WithHeader("Host", "httpbin.org").Expect()
		respPost.Status(http.StatusOK)

		respPut := s.NewAPISIXClient().PUT("/put").WithHeader("Host", "httpbin.org").Expect()
		respPut.Status(http.StatusOK)
	})
})

var _ = ginkgo.Describe("suite-annotations: blocklist-http-methods annotations", func() {
	s := scaffold.NewDefaultScaffold()

	ginkgo.It("enable in ingress networking/v1", func() {
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/http-block-methods: GET
  name: ingress-v1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /*
        pathType: Prefix
        backend:
          service:
            name: %s
            port:
              number: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		respGet := s.NewAPISIXClient().GET("/get").WithHeader("Host", "httpbin.org").Expect()
		respGet.Status(http.StatusMethodNotAllowed)

		respPost := s.NewAPISIXClient().POST("/post").WithHeader("Host", "httpbin.org").Expect()
		respPost.Status(http.StatusOK)

		respPut := s.NewAPISIXClient().PUT("/put").WithHeader("Host", "httpbin.org").Expect()
		respPut.Status(http.StatusOK)
	})

	ginkgo.It("enable in ingress networking/v1beta1", func() {
		if os.Getenv("K8s_Version") == "v1.24.0" {
			return
		}
		backendSvc, backendPort := s.DefaultHTTPBackend()
		ing := fmt.Sprintf(`
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: apisix
    k8s.apisix.apache.org/http-block-methods: GET
  name: ingress-v1beta1
spec:
  rules:
  - host: httpbin.org
    http:
      paths:
      - path: /*
        pathType: Prefix
        backend:
          serviceName: %s
          servicePort: %d
`, backendSvc, backendPort[0])
		err := s.CreateResourceFromString(ing)
		assert.Nil(ginkgo.GinkgoT(), err, "creating ingress")
		time.Sleep(5 * time.Second)

		respGet := s.NewAPISIXClient().GET("/get").WithHeader("Host", "httpbin.org").Expect()
		respGet.Status(http.StatusMethodNotAllowed)

		respPost := s.NewAPISIXClient().POST("/post").WithHeader("Host", "httpbin.org").Expect()
		respPost.Status(http.StatusOK)

		respPut := s.NewAPISIXClient().PUT("/put").WithHeader("Host", "httpbin.org").Expect()
		respPut.Status(http.StatusOK)
	})
})
