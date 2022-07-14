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
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/apache/apisix-ingress-controller/pkg/config"
	"github.com/apache/apisix-ingress-controller/pkg/kube"
	"github.com/gavv/httpexpect/v2"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/testing"
	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type Options struct {
	Name                       string
	Kubeconfig                 string
	APISIXConfigPath           string
	IngressAPISIXReplicas      int
	HTTPBinServicePort         int
	APISIXRouteVersion         string
	APISIXTlsVersion           string
	APISIXConsumerVersion      string
	ApisixPluginConfigVersion  string
	APISIXClusterConfigVersion string
	APISIXAdminAPIKey          string
	EnableWebhooks             bool
	APISIXPublishAddress       string
	disableNamespaceSelector   bool
	ApisixResourceSyncInterval string
	ApisixResourceVersion      string
}

type Scaffold struct {
	opts               *Options
	kubectlOptions     *k8s.KubectlOptions
	namespace          string
	t                  testing.TestingT
	nodes              []corev1.Node
	etcdService        *corev1.Service
	apisixService      *corev1.Service
	httpbinService     *corev1.Service
	testBackendService *corev1.Service
	finializers        []func()

	apisixAdminTunnel      *k8s.Tunnel
	apisixHttpTunnel       *k8s.Tunnel
	apisixHttpsTunnel      *k8s.Tunnel
	apisixTCPTunnel        *k8s.Tunnel
	apisixTLSOverTCPTunnel *k8s.Tunnel
	apisixUDPTunnel        *k8s.Tunnel
	apisixControlTunnel    *k8s.Tunnel

	// Used for template rendering.
	EtcdServiceFQDN string
}

type apisixResourceVersionInfo struct {
	V2      string
	V2beta3 string
}

var (
	apisixResourceVersion = &apisixResourceVersionInfo{
		V2:      "apisix.apache.org/v2",
		V2beta3: "apisix.apache.org/v2beta3",
	}
)

// GetKubeconfig returns the kubeconfig file path.
// Order:
// env KUBECONFIG;
// ~/.kube/config;
// "" (in case in-cluster configuration will be used).
func GetKubeconfig() string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		kubeconfig = filepath.Join(u.HomeDir, ".kube", "config")
		if _, err := os.Stat(kubeconfig); err != nil && !os.IsNotExist(err) {
			kubeconfig = ""
		}
	}
	return kubeconfig
}

// NewScaffold creates an e2e test scaffold.
func NewScaffold(o *Options) *Scaffold {
	if o.ApisixResourceVersion == ApisixResourceVersion().V2 {
		if o.APISIXRouteVersion == "" {
			o.APISIXRouteVersion = kube.ApisixRouteV2
		}
		if o.APISIXTlsVersion == "" {
			o.APISIXTlsVersion = config.ApisixV2
		}
		if o.APISIXConsumerVersion == "" {
			o.APISIXConsumerVersion = config.ApisixV2
		}
		if o.ApisixPluginConfigVersion == "" {
			o.ApisixPluginConfigVersion = config.ApisixV2
		}
		if o.APISIXClusterConfigVersion == "" {
			o.APISIXClusterConfigVersion = config.ApisixV2
		}
	} else {
		if o.APISIXRouteVersion == "" {
			o.APISIXRouteVersion = kube.ApisixRouteV2beta3
		}
		if o.APISIXTlsVersion == "" {
			o.APISIXTlsVersion = config.ApisixV2beta3
		}
		if o.APISIXConsumerVersion == "" {
			o.APISIXConsumerVersion = config.ApisixV2beta3
		}
		if o.ApisixPluginConfigVersion == "" {
			o.ApisixPluginConfigVersion = config.ApisixV2beta3
		}
		if o.APISIXClusterConfigVersion == "" {
			o.APISIXClusterConfigVersion = config.ApisixV2beta3
		}
	}
	if o.APISIXAdminAPIKey == "" {
		o.APISIXAdminAPIKey = "edd1c9f034335f136f87ad84b625c8f1"
	}
	if o.ApisixResourceSyncInterval == "" {
		o.ApisixResourceSyncInterval = "300s"
	}
	if o.Kubeconfig == "" {
		o.Kubeconfig = GetKubeconfig()
	}
	if o.APISIXConfigPath == "" {
		o.APISIXConfigPath = "testdata/apisix-gw-config.yaml"
	}
	if o.HTTPBinServicePort == 0 {
		o.HTTPBinServicePort = 80
	}
	defer ginkgo.GinkgoRecover()

	s := &Scaffold{
		opts: o,
		t:    ginkgo.GinkgoT(),
	}

	ginkgo.BeforeEach(s.beforeEach)
	ginkgo.AfterEach(s.afterEach)

	return s
}

// NewDefaultScaffold creates a scaffold with some default options.
func NewDefaultScaffold() *Scaffold {
	opts := &Options{
		Name:                  "default",
		IngressAPISIXReplicas: 1,
		ApisixResourceVersion: ApisixResourceVersion().V2beta3,
	}
	return NewScaffold(opts)
}

// NewDefaultV2Scaffold creates a scaffold with some default options.
func NewDefaultV2Scaffold() *Scaffold {
	opts := &Options{
		Name:                  "default",
		IngressAPISIXReplicas: 1,
		ApisixResourceVersion: ApisixResourceVersion().V2,
	}
	return NewScaffold(opts)
}

// KillPod kill the pod which name is podName.
func (s *Scaffold) KillPod(podName string) error {
	cli, err := k8s.GetKubernetesClientE(s.t)
	if err != nil {
		return err
	}
	return cli.CoreV1().Pods(s.namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
}

// DefaultHTTPBackend returns the service name and service ports
// of the default http backend.
func (s *Scaffold) DefaultHTTPBackend() (string, []int32) {
	var ports []int32
	for _, p := range s.httpbinService.Spec.Ports {
		ports = append(ports, p.Port)
	}
	return s.httpbinService.Name, ports
}

// ApisixAdminServiceAndPort returns the apisix service name and
// it's admin port.
func (s *Scaffold) ApisixAdminServiceAndPort() (string, int32) {
	return "apisix-service-e2e-test", 9180
}

// NewAPISIXClient creates the default HTTP client.
func (s *Scaffold) NewAPISIXClient() *httpexpect.Expect {
	u := url.URL{
		Scheme: "http",
		Host:   s.apisixHttpTunnel.Endpoint(),
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: u.String(),
		Client: &http.Client{
			Transport: &http.Transport{},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(
			httpexpect.NewAssertReporter(ginkgo.GinkgoT()),
		),
	})
}

// GetAPISIXHTTPSEndpoint get apisix https endpoint from tunnel map
func (s *Scaffold) GetAPISIXHTTPSEndpoint() string {
	return s.apisixHttpsTunnel.Endpoint()
}

// NewAPISIXClientWithTCPProxy creates the HTTP client but with the TCP proxy of APISIX.
func (s *Scaffold) NewAPISIXClientWithTCPProxy() *httpexpect.Expect {
	u := url.URL{
		Scheme: "http",
		Host:   s.apisixTCPTunnel.Endpoint(),
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: u.String(),
		Client: &http.Client{
			Transport: &http.Transport{},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(
			httpexpect.NewAssertReporter(ginkgo.GinkgoT()),
		),
	})
}

// NewAPISIXClientWithTLSOverTCP creates a TSL over TCP client
func (s *Scaffold) NewAPISIXClientWithTLSOverTCP(host string) *httpexpect.Expect {
	u := url.URL{
		Scheme: "https",
		Host:   s.apisixTLSOverTCPTunnel.Endpoint(),
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: u.String(),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// accept any certificate; for testing only!
					InsecureSkipVerify: true,
					ServerName:         host,
				},
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Reporter: httpexpect.NewAssertReporter(
			httpexpect.NewAssertReporter(ginkgo.GinkgoT()),
		),
	})
}

func (s *Scaffold) DNSResolver() *net.Resolver {
	return &net.Resolver{
		PreferGo: false,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			return d.DialContext(ctx, "udp", s.apisixUDPTunnel.Endpoint())
		},
	}
}

func (s *Scaffold) UpdateNamespace(ns string) {
	s.kubectlOptions.Namespace = ns
}

// NewAPISIXHttpsClient creates the default HTTPS client.
func (s *Scaffold) NewAPISIXHttpsClient(host string) *httpexpect.Expect {
	u := url.URL{
		Scheme: "https",
		Host:   s.apisixHttpsTunnel.Endpoint(),
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: u.String(),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// accept any certificate; for testing only!
					InsecureSkipVerify: true,
					ServerName:         host,
				},
			},
		},
		Reporter: httpexpect.NewAssertReporter(
			httpexpect.NewAssertReporter(ginkgo.GinkgoT()),
		),
	})
}

// NewAPISIXHttpsClientWithCertificates creates the default HTTPS client with giving trusted CA and client certs.
func (s *Scaffold) NewAPISIXHttpsClientWithCertificates(host string, insecure bool, ca *x509.CertPool, certs []tls.Certificate) *httpexpect.Expect {
	u := url.URL{
		Scheme: "https",
		Host:   s.apisixHttpsTunnel.Endpoint(),
	}
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: u.String(),
		Client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: insecure,
					ServerName:         host,
					RootCAs:            ca,
					Certificates:       certs,
				},
			},
		},
		Reporter: httpexpect.NewAssertReporter(
			httpexpect.NewAssertReporter(ginkgo.GinkgoT()),
		),
	})
}

// APISIXGatewayServiceEndpoint returns the apisix http gateway endpoint.
func (s *Scaffold) APISIXGatewayServiceEndpoint() string {
	return s.apisixHttpTunnel.Endpoint()
}

// RestartAPISIXDeploy delete apisix pod and wait new pod be ready
func (s *Scaffold) RestartAPISIXDeploy() {
	s.shutdownApisixTunnel()
	pods, err := k8s.ListPodsE(s.t, s.kubectlOptions, metav1.ListOptions{
		LabelSelector: "app=apisix-deployment-e2e-test",
	})
	assert.NoError(s.t, err, "list apisix pod")
	for _, pod := range pods {
		err = s.KillPod(pod.Name)
		assert.NoError(s.t, err, "killing apisix pod")
	}
	err = s.waitAllAPISIXPodsAvailable()
	assert.NoError(s.t, err, "waiting for new apisix instance ready")
	err = s.newAPISIXTunnels()
	assert.NoError(s.t, err, "renew apisix tunnels")
}

func (s *Scaffold) beforeEach() {
	var err error
	s.namespace = fmt.Sprintf("ingress-apisix-e2e-tests-%s-%d", s.opts.Name, time.Now().Nanosecond())
	s.kubectlOptions = &k8s.KubectlOptions{
		ConfigPath: s.opts.Kubeconfig,
		Namespace:  s.namespace,
	}

	s.finializers = nil
	labels := make(map[string]string)
	labels["apisix.ingress.watch"] = s.namespace
	k8s.CreateNamespaceWithMetadata(s.t, s.kubectlOptions, metav1.ObjectMeta{Name: s.namespace, Labels: labels})

	s.nodes, err = k8s.GetReadyNodesE(s.t, s.kubectlOptions)
	assert.Nil(s.t, err, "querying ready nodes")

	s.etcdService, err = s.newEtcd()
	assert.Nil(s.t, err, "initializing etcd")

	err = s.waitAllEtcdPodsAvailable()
	assert.Nil(s.t, err, "waiting for etcd ready")

	s.apisixService, err = s.newAPISIX()
	assert.Nil(s.t, err, "initializing Apache APISIX")

	err = s.waitAllAPISIXPodsAvailable()
	assert.Nil(s.t, err, "waiting for apisix ready")

	err = s.newAPISIXTunnels()
	assert.Nil(s.t, err, "creating apisix tunnels")

	s.httpbinService, err = s.newHTTPBIN()
	assert.Nil(s.t, err, "initializing httpbin")

	k8s.WaitUntilServiceAvailable(s.t, s.kubectlOptions, s.httpbinService.Name, 3, 2*time.Second)

	s.testBackendService, err = s.newTestBackend()
	assert.Nil(s.t, err, "initializing test backend")

	k8s.WaitUntilServiceAvailable(s.t, s.kubectlOptions, s.testBackendService.Name, 3, 2*time.Second)

	if s.opts.EnableWebhooks {
		err := generateWebhookCert(s.namespace)
		assert.Nil(s.t, err, "generate certs and create webhook secret")
	}

	err = s.newIngressAPISIXController()
	assert.Nil(s.t, err, "initializing ingress apisix controller")

	err = s.WaitAllIngressControllerPodsAvailable()
	assert.Nil(s.t, err, "waiting for ingress apisix controller ready")
}

func (s *Scaffold) afterEach() {
	defer ginkgo.GinkgoRecover()

	if ginkgo.CurrentSpecReport().Failed() {
		//_, _ = fmt.Fprintln(ginkgo.GinkgoWriter, "Dumping namespace contents")
		//output, _ := k8s.RunKubectlAndGetOutputE(ginkgo.GinkgoT(), s.kubectlOptions, "get", "deploy,sts,svc,pods")
		//if output != "" {
		//	_, _ = fmt.Fprintln(ginkgo.GinkgoWriter, output)
		//}
		//output, _ = k8s.RunKubectlAndGetOutputE(ginkgo.GinkgoT(), s.kubectlOptions, "describe", "pods")
		//if output != "" {
		//	_, _ = fmt.Fprintln(ginkgo.GinkgoWriter, output)
		//}
		//// Get the logs of apisix
		//output = s.GetDeploymentLogs("apisix-deployment-e2e-test")
		//if output != "" {
		//	_, _ = fmt.Fprintln(ginkgo.GinkgoWriter, output)
		//}
		//// Get the logs of ingress
		//output = s.GetDeploymentLogs("ingress-apisix-controller-deployment-e2e-test")
		//if output != "" {
		//	_, _ = fmt.Fprintln(ginkgo.GinkgoWriter, output)
		//}
	}

	err := k8s.DeleteNamespaceE(s.t, s.kubectlOptions, s.namespace)
	assert.Nilf(ginkgo.GinkgoT(), err, "deleting namespace %s", s.namespace)

	for _, f := range s.finializers {
		runWithRecover(f)
	}

	// Wait for a while to prevent the worker node being overwhelming
	// (new cases will be run).
	time.Sleep(3 * time.Second)
}

func runWithRecover(f func()) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if ok {
			// just ignore already closed channel
			if strings.Contains(err.Error(), "close of closed channel") {
				return
			}
		}
		panic(r)
	}()
	f()
}

func (s *Scaffold) GetDeploymentLogs(name string) string {
	cli, err := k8s.GetKubernetesClientE(s.t)
	if err != nil {
		assert.Nilf(ginkgo.GinkgoT(), err, "get client error: %s", err.Error())
	}
	pods, err := cli.CoreV1().Pods(s.namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app=" + name,
	})
	if err != nil {
		return ""
	}
	var buf strings.Builder
	for _, pod := range pods.Items {
		buf.WriteString(fmt.Sprintf("=== pod: %s ===\n", pod.Name))
		logs, err := cli.CoreV1().RESTClient().Get().
			Resource("pods").
			Namespace(s.namespace).
			Name(pod.Name).SubResource("log").
			Param("container", name).
			Do(context.TODO()).
			Raw()
		if err == nil {
			buf.Write(logs)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (s *Scaffold) addFinalizers(f func()) {
	s.finializers = append(s.finializers, f)
}

func (s *Scaffold) renderConfig(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	t := template.Must(template.New(path).Parse(string(data)))
	if err := t.Execute(&buf, s); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// FormatRegistry replace default registry to custom registry if exist
func (s *Scaffold) FormatRegistry(workloadTemplate string) string {
	customRegistry, isExist := os.LookupEnv("REGISTRY")
	if isExist {
		return strings.Replace(workloadTemplate, "localhost:5000", customRegistry, -1)
	} else {
		return workloadTemplate
	}
}

// FormatNamespaceLabel set label to be empty if s.opts.disableNamespaceSelector is true.
func (s *Scaffold) FormatNamespaceLabel(label string) string {
	if s.opts.disableNamespaceSelector {
		return "\"\""
	}
	return label
}

var (
	versionRegex = regexp.MustCompile(`apiVersion: apisix.apache.org/v.*?\n`)
	kindRegex    = regexp.MustCompile(`kind: (.*?)\n`)
)

func (s *Scaffold) replaceApiVersion(yml, ver string) string {
	return versionRegex.ReplaceAllString(yml, "apiVersion: "+ver+"\n")
}

func (s *Scaffold) getKindValue(yml string) string {
	subStr := kindRegex.FindStringSubmatch(yml)
	if len(subStr) < 2 {
		return ""
	}
	return subStr[1]
}

func (s *Scaffold) DisableNamespaceSelector() {
	s.opts.disableNamespaceSelector = true
}

func waitExponentialBackoff(condFunc func() (bool, error)) error {
	backoff := wait.Backoff{
		Duration: 500 * time.Millisecond,
		Factor:   2,
		Steps:    8,
	}
	return wait.ExponentialBackoff(backoff, condFunc)
}

// generateWebhookCert generates signed certs of webhook and create the corresponding secret by running a script.
func generateWebhookCert(ns string) error {
	commandTemplate := `testdata/webhook-create-signed-cert.sh`
	cmd := exec.Command("/bin/sh", commandTemplate, "--namespace", ns)

	output, err := cmd.Output()
	if err != nil {
		ginkgo.GinkgoT().Errorf("%s", output)
		return fmt.Errorf("failed to execute the script: %v", err)
	}

	return nil
}

func (s *Scaffold) CreateVersionedApisixPluginConfig(yml string) error {
	if !strings.Contains(yml, "kind: ApisixPluginConfig") {
		return errors.New("not a ApisixPluginConfig")
	}

	ac := s.replaceApiVersion(yml, s.opts.ApisixPluginConfigVersion)
	return s.CreateResourceFromString(ac)
}

func (s *Scaffold) CreateVersionedApisixResource(yml string) error {
	kindValue := s.getKindValue(yml)
	switch kindValue {
	case "ApisixRoute":
		ar := s.replaceApiVersion(yml, s.opts.APISIXRouteVersion)
		return s.CreateResourceFromString(ar)
	case "ApisixConsumer":
		ac := s.replaceApiVersion(yml, s.opts.APISIXConsumerVersion)
		return s.CreateResourceFromString(ac)
	case "ApisixPluginConfig":
		apc := s.replaceApiVersion(yml, s.opts.ApisixPluginConfigVersion)
		return s.CreateResourceFromString(apc)
	}
	return fmt.Errorf("the resource %s does not support", kindValue)
}

func (s *Scaffold) CreateVersionedApisixResourceWithNamespace(yml, namespace string) error {
	kindValue := s.getKindValue(yml)
	switch kindValue {
	case "ApisixRoute":
		ar := s.replaceApiVersion(yml, s.opts.APISIXRouteVersion)
		return s.CreateResourceFromStringWithNamespace(ar, namespace)
	case "ApisixConsumer":
		ac := s.replaceApiVersion(yml, s.opts.APISIXConsumerVersion)
		return s.CreateResourceFromStringWithNamespace(ac, namespace)
	case "ApisixPluginConfig":
		apc := s.replaceApiVersion(yml, s.opts.ApisixPluginConfigVersion)
		return s.CreateResourceFromStringWithNamespace(apc, namespace)
	}
	return fmt.Errorf("the resource %s does not support", kindValue)
}

func ApisixResourceVersion() *apisixResourceVersionInfo {
	return apisixResourceVersion
}
