package ingress

import (
	"context"
	"github.com/apache/apisix-ingress-controller/pkg/config"
	"github.com/apache/apisix-ingress-controller/pkg/kube"
	"github.com/apache/apisix-ingress-controller/pkg/providers/namespace"
	"github.com/apache/apisix-ingress-controller/pkg/providers/translation"
	providertypes "github.com/apache/apisix-ingress-controller/pkg/providers/types"
	"github.com/apache/apisix-ingress-controller/pkg/providers/utils"
	"k8s.io/client-go/tools/cache"
)

const (
	ProviderName = "Ingress"
)

type ingressCommon struct {
	*providertypes.Common

	namespaceProvider namespace.WatchingNamespaceProvider
	translator        IngressTranslator
}

var _ Provider = (*ingressProvider)(nil)

type Provider interface {
	providertypes.Provider

	ResourceSync()
}

type ingressProvider struct {
	name string

	ingressController *ingressController

	ingressInformer cache.SharedIndexInformer
}

func NewProvider(common *providertypes.Common, namespaceProvider namespace.WatchingNamespaceProvider, translator translation.Translator) (Provider, error) {
	p := &ingressProvider{
		name: ProviderName,
	}

	kubeFactory := common.KubeClient.NewSharedIndexInformerFactory()

	svcLister := kubeFactory.Core().V1().Services().Lister()
	ingressLister := kube.NewIngressLister(
		kubeFactory.Networking().V1().Ingresses().Lister(),
		kubeFactory.Networking().V1beta1().Ingresses().Lister(),
		kubeFactory.Extensions().V1beta1().Ingresses().Lister(),
	)
	switch common.Config.Kubernetes.IngressVersion {
	case config.IngressNetworkingV1:
		p.ingressInformer = kubeFactory.Networking().V1().Ingresses().Informer()
	case config.IngressNetworkingV1beta1:
		p.ingressInformer = kubeFactory.Networking().V1beta1().Ingresses().Informer()
	default:
		p.ingressInformer = kubeFactory.Extensions().V1beta1().Ingresses().Informer()
	}

	c := &ingressCommon{
		Common:            common,
		namespaceProvider: namespaceProvider,
		translator:        newIngressTranslator(svcLister, translator),
	}

	p.ingressController = newIngressController(c, ingressLister, p.ingressInformer)

	return p, nil
}

func (p *ingressProvider) Run(ctx context.Context) {
	e := utils.ParallelExecutor{}

	e.Add(func() {
		p.ingressInformer.Run(ctx.Done())
	})

	e.Add(func() {
		p.ingressController.run(ctx)
	})

	e.Wait()
}

func (p *ingressProvider) ResourceSync() {
	e := utils.ParallelExecutor{}

	e.Add(p.ingressController.ResourceSync)

	e.Wait()
}
