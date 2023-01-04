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
package apisix

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"github.com/apache/apisix-ingress-controller/pkg/config"
	"github.com/apache/apisix-ingress-controller/pkg/kube"
	configv2 "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2"
	"github.com/apache/apisix-ingress-controller/pkg/log"
	"github.com/apache/apisix-ingress-controller/pkg/providers/utils"
	"github.com/apache/apisix-ingress-controller/pkg/types"
)

type apisixGlobalRuleController struct {
	*apisixCommon

	workqueue workqueue.RateLimitingInterface
	workers   int
}

func newApisixGlobalRuleController(common *apisixCommon) *apisixGlobalRuleController {
	c := &apisixGlobalRuleController{
		apisixCommon: common,
		workqueue:    workqueue.NewNamedRateLimitingQueue(workqueue.NewItemFastSlowRateLimiter(1*time.Second, 60*time.Second, 5), "ApisixGlobalRule"),
		workers:      1,
	}

	c.ApisixGlobalRuleInformer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		},
	)
	return c
}

func (c *apisixGlobalRuleController) run(ctx context.Context) {
	log.Info("ApisixGlobalRule controller started")
	defer log.Info("ApisixGlobalRule controller exited")
	defer c.workqueue.ShutDown()

	for i := 0; i < c.workers; i++ {
		go c.runWorker(ctx)
	}
	<-ctx.Done()
}

func (c *apisixGlobalRuleController) runWorker(ctx context.Context) {
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

func (c *apisixGlobalRuleController) sync(ctx context.Context, ev *types.Event) error {
	obj := ev.Object.(kube.ApisixGlobalRuleEvent)
	namespace, name, err := cache.SplitMetaNamespaceKey(obj.Key)
	if err != nil {
		log.Errorf("invalid resource key: %s", obj.Key)
		return err
	}
	var (
		agr kube.ApisixGlobalRule
	)
	agr, err = c.ApisixGlobalRuleLister.ApisixGlobalRule(namespace, name)
	if err != nil {
		if !k8serrors.IsNotFound(err) {
			log.Errorw("failed to get ApisixGlobalRule",
				zap.String("version", obj.GroupVersion),
				zap.String("key", obj.Key),
				zap.Error(err),
			)
			return err
		}

		if ev.Type != types.EventDelete {
			log.Warnw("ApisixGlobalRule was deleted before it can be delivered",
				zap.String("key", obj.Key),
				zap.String("version", obj.GroupVersion),
			)
			return nil
		}
	}
	if ev.Type == types.EventDelete {
		if agr != nil {
			// We still find the resource while we are processing the DELETE event,
			// that means object with same namespace and name was created, discarding
			// this stale DELETE event.
			log.Warnw("discard the stale ApisixGlobalRule delete event since the resource still exists",
				zap.String("key", obj.Key),
			)
			return nil
		}
		agr = ev.Tombstone.(kube.ApisixGlobalRule)
	}

	tctx, err := c.translator.TranslateGlobalRule(agr)

	log.Debugw("translated ApisixGlobalRule",
		zap.Any("globalrules", tctx.GlobalRules),
	)

	m := &utils.Manifest{
		GlobalRules: tctx.GlobalRules,
	}

	var (
		added   *utils.Manifest
		updated *utils.Manifest
		deleted *utils.Manifest
	)

	if ev.Type == types.EventDelete {
		deleted = m
	} else if ev.Type == types.EventAdd {
		added = m
	} else {
		oldCtx, err := c.translator.TranslateGlobalRule(obj.OldObject)
		if err != nil {
			log.Errorw("failed to translate old ApisixGlobalRule",
				zap.String("version", obj.GroupVersion),
				zap.String("event", "update"),
				zap.Error(err),
				zap.Any("ApisixGlobalRule", agr),
			)
			return err
		}

		om := &utils.Manifest{
			GlobalRules: oldCtx.GlobalRules,
		}
		added, updated, deleted = m.Diff(om)
	}

	return c.SyncManifests(ctx, added, updated, deleted)
}

func (c *apisixGlobalRuleController) handleSyncErr(obj interface{}, errOrigin error) {
	ev := obj.(*types.Event)
	event := ev.Object.(kube.ApisixGlobalRuleEvent)
	if k8serrors.IsNotFound(errOrigin) && ev.Type != types.EventDelete {
		log.Infow("sync ApisixGlobalRule but not found, ignore",
			zap.String("event_type", ev.Type.String()),
			zap.String("ApisixGlobalRule", ev.Object.(kube.ApisixGlobalRuleEvent).Key),
		)
		c.workqueue.Forget(event)
		return
	}
	namespace, name, errLocal := cache.SplitMetaNamespaceKey(event.Key)
	if errLocal != nil {
		log.Errorf("invalid resource key: %s", event.Key)
		c.MetricsCollector.IncrSyncOperation("PluginConfig", "failure")
		return
	}
	var apc kube.ApisixGlobalRule
	switch event.GroupVersion {
	case config.ApisixV2:
		apc, errLocal = c.ApisixGlobalRuleLister.V2(namespace, name)
	default:
		errLocal = fmt.Errorf("unsupported ApisixGlobalRule group version %s", event.GroupVersion)
	}
	if errOrigin == nil {
		if ev.Type != types.EventDelete {
			if errLocal == nil {
				switch apc.GroupVersion() {
				case config.ApisixV2:
					c.RecordEvent(apc.V2(), v1.EventTypeNormal, utils.ResourceSynced, nil)
					c.recordStatus(apc.V2(), utils.ResourceSynced, nil, metav1.ConditionTrue, apc.V2().GetGeneration())
				}
			} else {
				log.Errorw("failed list ApisixGlobalRule",
					zap.Error(errLocal),
					zap.String("name", name),
					zap.String("namespace", namespace),
				)
			}
		}
		c.workqueue.Forget(obj)
		c.MetricsCollector.IncrSyncOperation("PluginConfig", "success")
		return
	}
	log.Warnw("sync ApisixGlobalRule failed, will retry",
		zap.Any("object", obj),
		zap.Error(errOrigin),
	)
	if errLocal == nil {
		switch apc.GroupVersion() {
		case config.ApisixV2:
			c.RecordEvent(apc.V2(), v1.EventTypeWarning, utils.ResourceSyncAborted, errOrigin)
			c.recordStatus(apc.V2(), utils.ResourceSyncAborted, errOrigin, metav1.ConditionFalse, apc.V2().GetGeneration())
		}
	} else {
		log.Errorw("failed list ApisixGlobalRule",
			zap.Error(errLocal),
			zap.String("name", name),
			zap.String("namespace", namespace),
		)
	}
	c.workqueue.AddRateLimited(obj)
	c.MetricsCollector.IncrSyncOperation("PluginConfig", "failure")
}

func (c *apisixGlobalRuleController) onAdd(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorf("found ApisixGlobalRule resource with bad meta namespace key: %s", err)
		return
	}
	if !c.namespaceProvider.IsWatchingNamespace(key) {
		return
	}
	log.Debugw("ApisixGlobalRule add event arrived",
		zap.Any("object", obj))

	apc := kube.MustNewApisixGlobalRule(obj)
	c.workqueue.Add(&types.Event{
		Type: types.EventAdd,
		Object: kube.ApisixGlobalRuleEvent{
			Key:          key,
			GroupVersion: apc.GroupVersion(),
		},
	})

	c.MetricsCollector.IncrEvents("PluginConfig", "add")
}

func (c *apisixGlobalRuleController) onUpdate(oldObj, newObj interface{}) {
	prev := kube.MustNewApisixGlobalRule(oldObj)
	curr := kube.MustNewApisixGlobalRule(newObj)
	if prev.ResourceVersion() >= curr.ResourceVersion() {
		return
	}
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err != nil {
		log.Errorf("found ApisixGlobalRule resource with bad meta namespace key: %s", err)
		return
	}
	if !c.namespaceProvider.IsWatchingNamespace(key) {
		return
	}
	log.Debugw("ApisixGlobalRule update event arrived",
		zap.Any("new object", curr),
		zap.Any("old object", prev),
	)
	c.workqueue.Add(&types.Event{
		Type: types.EventUpdate,
		Object: kube.ApisixGlobalRuleEvent{
			Key:          key,
			GroupVersion: curr.GroupVersion(),
			OldObject:    prev,
		},
	})

	c.MetricsCollector.IncrEvents("GlobalRule", "update")
}

func (c *apisixGlobalRuleController) onDelete(obj interface{}) {
	apc, err := kube.NewApisixGlobalRule(obj)
	if err != nil {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			return
		}
		apc = kube.MustNewApisixGlobalRule(tombstone)
	}
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Errorf("found ApisixGlobalRule resource with bad meta namesapce key: %s", err)
		return
	}
	if !c.namespaceProvider.IsWatchingNamespace(key) {
		return
	}
	log.Debugw("ApisixGlobalRule delete event arrived",
		zap.Any("final state", apc),
	)
	c.workqueue.Add(&types.Event{
		Type: types.EventDelete,
		Object: kube.ApisixGlobalRuleEvent{
			Key:          key,
			GroupVersion: apc.GroupVersion(),
		},
		Tombstone: apc,
	})

	c.MetricsCollector.IncrEvents("GlobalRule", "delete")
}

func (c *apisixGlobalRuleController) ResourceSync() {
	objs := c.ApisixGlobalRuleInformer.GetIndexer().List()
	for _, obj := range objs {
		key, err := cache.MetaNamespaceKeyFunc(obj)
		if err != nil {
			log.Errorw("ApisixGlobalRule sync failed, found ApisixGlobalRule resource with bad meta namespace key", zap.String("error", err.Error()))
			continue
		}
		if !c.namespaceProvider.IsWatchingNamespace(key) {
			continue
		}
		apc := kube.MustNewApisixGlobalRule(obj)
		c.workqueue.Add(&types.Event{
			Type: types.EventAdd,
			Object: kube.ApisixGlobalRuleEvent{
				Key:          key,
				GroupVersion: apc.GroupVersion(),
			},
		})
	}
}

// recordStatus record resources status
func (c *apisixGlobalRuleController) recordStatus(at interface{}, reason string, err error, status metav1.ConditionStatus, generation int64) {
	// build condition
	message := utils.CommonSuccessMessage
	if err != nil {
		message = err.Error()
	}
	condition := metav1.Condition{
		Type:               utils.ConditionType,
		Reason:             reason,
		Status:             status,
		Message:            message,
		ObservedGeneration: generation,
	}
	apisixClient := c.KubeClient.APISIXClient

	if kubeObj, ok := at.(runtime.Object); ok {
		at = kubeObj.DeepCopyObject()
	}

	switch v := at.(type) {
	case *configv2.ApisixGlobalRule:
		// set to status
		if v.Status.Conditions == nil {
			conditions := make([]metav1.Condition, 0)
			v.Status.Conditions = conditions
		}
		//
		if utils.VerifyGeneration(&v.Status.Conditions, condition) && !meta.IsStatusConditionPresentAndEqual(v.Status.Conditions, condition.Type, condition.Status) {
			meta.SetStatusCondition(&v.Status.Conditions, condition)
			if _, errRecord := apisixClient.ApisixV2().ApisixGlobalRules(v.Namespace).
				UpdateStatus(context.TODO(), v, metav1.UpdateOptions{}); errRecord != nil {
				log.Errorw("failed to record status change for ApisixGlobalRule",
					zap.Error(errRecord),
					zap.String("name", v.Name),
					zap.String("namespace", v.Namespace),
				)
			}
		}
	default:
		// This should not be executed
		log.Errorf("unsupported resource record: %s", v)
	}
}
