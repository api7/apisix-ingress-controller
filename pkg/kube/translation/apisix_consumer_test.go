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
package translation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	configv2beta3 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2beta3"
	apisixv1 "github.com/apache/apisix-ingress-controller/pkg/types/apisix/v1"
)

func TestTranslateApisixConsumer(t *testing.T) {
	ac := &configv2beta3.ApisixConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jack",
			Namespace: "qa",
		},
		Spec: configv2beta3.ApisixConsumerSpec{
			AuthParameter: configv2beta3.ApisixConsumerAuthParameter{
				BasicAuth: &configv2beta3.ApisixConsumerBasicAuth{
					Value: &configv2beta3.ApisixConsumerBasicAuthValue{
						Username: "jack",
						Password: "jacknice",
					},
				},
			},
		},
	}
	consumer, err := (&translator{}).TranslateApisixConsumer(ac)
	assert.Nil(t, err)
	assert.Len(t, consumer.Plugins, 1)
	cfg := consumer.Plugins["basic-auth"].(*apisixv1.BasicAuthConsumerConfig)
	assert.Equal(t, "jack", cfg.Username)
	assert.Equal(t, "jacknice", cfg.Password)

	ac = &configv2beta3.ApisixConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jack",
			Namespace: "qa",
		},
		Spec: configv2beta3.ApisixConsumerSpec{
			AuthParameter: configv2beta3.ApisixConsumerAuthParameter{
				KeyAuth: &configv2beta3.ApisixConsumerKeyAuth{
					Value: &configv2beta3.ApisixConsumerKeyAuthValue{
						Key: "qwerty",
					},
				},
			},
		},
	}
	consumer, err = (&translator{}).TranslateApisixConsumer(ac)
	assert.Nil(t, err)
	assert.Len(t, consumer.Plugins, 1)
	cfg2 := consumer.Plugins["key-auth"].(*apisixv1.KeyAuthConsumerConfig)
	assert.Equal(t, "qwerty", cfg2.Key)

	ac = &configv2beta3.ApisixConsumer{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "jack",
			Namespace: "qa",
		},
		Spec: configv2beta3.ApisixConsumerSpec{
			AuthParameter: configv2beta3.ApisixConsumerAuthParameter{
				JwtAuth: &configv2beta3.ApisixConsumerJwtAuth{
					Value: &configv2beta3.ApisixConsumerJwtAuthValue{
						Key:          "foo",
						Secret:       "123",
						PublicKey:    "public",
						PrivateKey:   "private",
						Algorithm:    "HS256",
						Exp:          int64(1000),
						Base64Secret: true,
					},
				},
			},
		},
	}
	consumer, err = (&translator{}).TranslateApisixConsumer(ac)
	assert.Nil(t, err)
	assert.Len(t, consumer.Plugins, 1)
	cfg3 := consumer.Plugins["jwt-auth"].(*apisixv1.JwtAuthConsumerConfig)
	assert.Equal(t, "foo", cfg3.Key)
	assert.Equal(t, "123", cfg3.Secret)
	assert.Equal(t, "public", cfg3.PublicKey)
	assert.Equal(t, "private", cfg3.PrivateKey)
	assert.Equal(t, "HS256", cfg3.Algorithm)
	assert.Equal(t, int64(1000), cfg3.Exp)
	assert.Equal(t, true, cfg3.Base64Secret)
	// No test test cases for secret references as we already test them
	// in plugin_test.go.
}
