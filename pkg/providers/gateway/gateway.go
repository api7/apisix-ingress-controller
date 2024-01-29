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
package gateway

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/apisix-ingress-controller/pkg/kube"
	"github.com/apache/apisix-ingress-controller/pkg/log"
	gatewaytypes "github.com/apache/apisix-ingress-controller/pkg/providers/gateway/types"
	"github.com/apache/apisix-ingress-controller/pkg/providers/utils"
	"github.com/apache/apisix-ingress-controller/pkg/types"
	"go.uber.org/zap"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayv1beta1 "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type gatewayController struct {
	controller *Provider
	workqueue  workqueue.RateLimitingInterface
	workers    int
}

func newGatewayController(c *Provider) *gatewayController {
	ctl := &gatewayController{
		controller: c,
		workqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.NewItemFastSlowRateLimiter(1*time.Second, 60*time.Second, 5), "Gateway"),
		workers:    1,
	}

	ctl.controller.gatewayInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    ctl.onAdd,
		UpdateFunc: ctl.onUpdate,
		DeleteFunc: ctl.OnDelete,
	})
	return ctl
}

func (c *gatewayController) run(ctx context.Context) {
	log.Info("gateway controller started")
	defer log.Info("gateway controller exited")
	defer c.workqueue.ShutDown()

	if !cache.WaitForCacheSync(ctx.Done(), c.controller.gatewayInformer.HasSynced) {
		log.Error("cache sync failed")
		return
	}

	for i := 0; i < c.workers; i++ {
		go c.runWorker(ctx)
	}
	<-ctx.Done()
}

func (c *gatewayController) runWorker(ctx context.Context) {
	for {
		obj, quit := c.workqueue.Get()
		if quit {
			return
		}
		err := c.sync(ctx, obj.(*types.Event))
		c.workqueue.Done(obj)
		c.handleSyncErr(obj, err)
	}
}

func (c *gatewayController) sync(ctx context.Context, ev *types.Event) error {
	gatewayEvent := ev.Object.(kube.GatewayEvent)
	namespace, name, err := cache.SplitMetaNamespaceKey(gatewayEvent.Key)
	if err != nil {
		log.Errorw("found Gateway resource with invalid meta namespace key",
			zap.Error(err),
			zap.String("key", gatewayEvent.Key),
		)
		return err
	}
	var gatev1 *gatewayv1.Gateway
	var gatev1beta *gatewayv1beta1.Gateway
	var generation int64
	switch gatewayEvent.GroupVersion {
	case kube.GatewayV1:
		gatev1, err = c.controller.gatewayListerV1.Gateways(namespace).Get(name)
		if err != nil {
			return err
		}
		generation = gatev1.Generation
	case kube.GatewayV1beta1:
		gatev1beta, err = c.controller.gatewayListerV1beta1.Gateways(namespace).Get(name)
		if err != nil {
			return err
		}

		fmt.Println("gateway returned by gatewaylister for v1beta1 ", gatev1beta)
		generation = gatev1beta.Generation
	}

	if err != nil {
		if !k8serrors.IsNotFound(err) {
			log.Errorw("failed to get Gateway",
				zap.Error(err),
				zap.String("key", gatewayEvent.Key),
			)
			return err
		}
		if ev.Type != types.EventDelete {
			log.Warnw("Gateway was deleted before it can be delivered",
				zap.String("key", gatewayEvent.Key),
			)
			// Don't need to retry.
			return nil
		}
	}

	if ev.Type == types.EventDelete {
		if gatev1 != nil && gatev1beta != nil {
			// We still find the resource while we are processing the DELETE event,
			// that means object with same namespace and name was created, discarding
			// this stale DELETE event.
			log.Warnw("discard the stale Gateway delete event since it exists",
				zap.String("key", gatewayEvent.Key),
			)
			return nil
		}

		switch gatewayEvent.GroupVersion {
		case kube.GatewayV1:
			err = c.controller.RemoveListeners(gatev1.Namespace, gatev1.Namespace)
		case kube.GatewayV1beta1:
			err = c.controller.RemoveListeners(gatev1beta.Namespace, gatev1beta.Namespace)
		}
		if err != nil {
			return err
		}
	} else {
		var gatewayClassName string
		switch gatewayEvent.GroupVersion {
		case kube.GatewayV1:
			gatewayClassName = string(gatev1.Spec.GatewayClassName)
		case kube.GatewayV1beta1:
			gatewayClassName = string(gatev1beta.Spec.GatewayClassName)
		}

		if c.controller.HasGatewayClass(gatewayClassName) {
			// TODO: handle listeners
			var listeners map[string]*gatewaytypes.ListenerConf
			switch gatewayEvent.GroupVersion {
			case kube.GatewayV1:
				gateway := gatev1
				listeners, err = c.controller.translator.TranslateGatewayV1(gateway)
				if err != nil {
					return err
				}

				err = c.controller.AddListeners(gateway.Namespace, gateway.Name, listeners)
				if err != nil {
					return err
				}
			case kube.GatewayV1beta1:
				gateway := gatev1beta
				listeners, err = c.controller.translator.TranslateGatewayV1beta1(gateway)
				if err != nil {
					return err
				}

				err = c.controller.AddListeners(gateway.Namespace, gateway.Name, listeners)
				if err != nil {
					return err
				}
			}

		} else {
			gatewayClass, err := c.controller.gatewayClassLister.Get(gatewayClassName)
			if err != nil {
				return err
			}
			if gatewayClass.Spec.ControllerName == GatewayClassName {
				log.Warn("gatewayClass not synced")
				return fmt.Errorf("wait gatewayClass %s synced", gatewayClassName)
			}
		}
	}

	// TODO The current implementation does not fully support the definition of Gateway.
	// We can update `spec.addresses` with the current data plane information.
	// At present, we choose to directly update `GatewayStatus.Addresses`
	// to indicate that we have picked the Gateway resource.
	switch gatewayEvent.GroupVersion {
	case kube.GatewayV1:
		c.recordStatusv1(gatev1, string(gatewayv1.ListenerReasonReady), metav1.ConditionTrue, generation)
	case kube.GatewayV1beta1:
		c.recordStatusv1beta(gatev1beta, string(gatewayv1.ListenerReasonReady), metav1.ConditionTrue, generation)
	}
	return nil
}

func (c *gatewayController) handleSyncErr(obj interface{}, err error) {
	if err == nil {
		c.workqueue.Forget(obj)
		c.controller.MetricsCollector.IncrSyncOperation("gateway", "success")
		return
	}
	event := obj.(*types.Event)
	if k8serrors.IsNotFound(err) && event.Type != types.EventDelete {
		log.Infow("sync gateway but not found, ignore",
			zap.String("event_type", event.Type.String()),
		)
		c.workqueue.Forget(event)
		return
	}
	log.Warnw("sync gateway failed, will retry",
		zap.Any("object", obj),
		zap.Error(err),
	)
	c.workqueue.AddRateLimited(obj)
	c.controller.MetricsCollector.IncrSyncOperation("gateway", "failure")
}

func (c *gatewayController) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorw("found gateway resource with bad meta namespace key",
			zap.Error(err),
			zap.Any("obj", obj),
		)
		return
	}
	if !c.controller.NamespaceProvider.IsWatchingNamespace(key) {
		return
	}
	log.Debugw("gateway add event arrived",
		zap.Any("object", obj),
	)
	gateway := kube.MustNewGateway(obj)
	c.workqueue.Add(&types.Event{
		Type: types.EventAdd,
		Object: kube.GatewayEvent{
			Key:          key,
			GroupVersion: gateway.GroupVersion(),
		},
	})
}

func (c *gatewayController) onUpdate(oldObj, newObj interface{}) {
}

func (c *gatewayController) OnDelete(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorw("failed to handle deletion Gateway meta key",
			zap.Error(err),
			zap.Any("obj", obj),
		)
		return
	}

	gateway := kube.MustNewGateway(obj)

	c.workqueue.Add(&types.Event{
		Type: types.EventDelete,
		Object: kube.GatewayEvent{
			Key:          key,
			GroupVersion: gateway.GroupVersion(),
		},
		Tombstone: gateway,
	})
}

func (c *gatewayController) recordStatusv1beta(v *gatewayv1beta1.Gateway, reason string, status metav1.ConditionStatus, generation int64) {
	v = v.DeepCopy()

	gatewayCondition := metav1.Condition{
		Type:               string(gatewayv1.ListenerConditionReady),
		Reason:             reason,
		Status:             status,
		Message:            "Gateway's status has been successfully updated",
		ObservedGeneration: generation,
	}

	if v.Status.Conditions == nil {
		conditions := make([]metav1.Condition, 0)
		v.Status.Conditions = conditions
	} else {
		meta.SetStatusCondition(&v.Status.Conditions, gatewayCondition)
	}

	lbips, err := utils.IngressLBStatusIPs(c.controller.Cfg.IngressPublishService, c.controller.Cfg.IngressStatusAddress, c.controller.ListerInformer.SvcLister)
	if err != nil {
		log.Errorw("failed to get APISIX gateway external IPs",
			zap.Error(err),
		)
	}

	v.Status.Addresses = utils.CoreV1ToGatewayV1beta1Addr(lbips)
	if _, errRecord := c.controller.gatewayClient.GatewayV1beta1().Gateways(v.Namespace).UpdateStatus(context.TODO(), v, metav1.UpdateOptions{}); errRecord != nil {
		log.Errorw("failed to record status change for Gateway resource",
			zap.Error(errRecord),
			zap.String("name", v.Name),
			zap.String("namespace", v.Namespace),
		)
	}
}

func (c *gatewayController) recordStatusv1(v *gatewayv1.Gateway, reason string, status metav1.ConditionStatus, generation int64) {
	v = v.DeepCopy()

	gatewayCondition := metav1.Condition{
		Type:               string(gatewayv1.ListenerConditionReady),
		Reason:             reason,
		Status:             status,
		Message:            "Gateway's status has been successfully updated",
		ObservedGeneration: generation,
	}

	if v.Status.Conditions == nil {
		conditions := make([]metav1.Condition, 0)
		v.Status.Conditions = conditions
	} else {
		meta.SetStatusCondition(&v.Status.Conditions, gatewayCondition)
	}

	lbips, err := utils.IngressLBStatusIPs(c.controller.Cfg.IngressPublishService, c.controller.Cfg.IngressStatusAddress, c.controller.ListerInformer.SvcLister)
	if err != nil {
		log.Errorw("failed to get APISIX gateway external IPs",
			zap.Error(err),
		)
	}

	v.Status.Addresses = utils.CoreV1ToGatewayV1beta1Addr(lbips)
	if _, errRecord := c.controller.gatewayClient.GatewayV1().Gateways(v.Namespace).UpdateStatus(context.TODO(), v, metav1.UpdateOptions{}); errRecord != nil {
		log.Errorw("failed to record status change for Gateway resource",
			zap.Error(errRecord),
			zap.String("name", v.Name),
			zap.String("namespace", v.Namespace),
		)
	}
}
