package dag

import (
	"context"
	"errors"

	"github.com/derailed/popeye/internal/client"
	"github.com/derailed/popeye/internal/dao"
	pv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// ListPodSecurityPolicies list all included PodSecurityPolicies.
func ListPodSecurityPolicies(ctx context.Context) (map[string]*pv1beta1.PodSecurityPolicy, error) {
	dps, err := listAllPodSecurityPolicys(ctx)
	if err != nil {
		return map[string]*pv1beta1.PodSecurityPolicy{}, err
	}

	f := mustExtractFactory(ctx)
	res := make(map[string]*pv1beta1.PodSecurityPolicy, len(dps))
	for fqn, dp := range dps {
		if includeNS(f.Client(), dp.Namespace) {
			res[fqn] = dp
		}
	}

	return res, nil
}

// ListAllPodSecurityPolicys fetch all PodSecurityPolicys on the cluster.
func listAllPodSecurityPolicys(ctx context.Context) (map[string]*pv1beta1.PodSecurityPolicy, error) {
	ll, err := fetchPodSecurityPolicys(ctx)
	if err != nil {
		return nil, err
	}

	dps := make(map[string]*pv1beta1.PodSecurityPolicy, len(ll.Items))
	for i := range ll.Items {
		dps[metaFQN(ll.Items[i].ObjectMeta)] = &ll.Items[i]
	}

	return dps, nil
}

// FetchPodSecurityPolicys retrieves all PodSecurityPolicys on the cluster.
func fetchPodSecurityPolicys(ctx context.Context) (*pv1beta1.PodSecurityPolicyList, error) {
	f, cfg := mustExtractFactory(ctx), mustExtractConfig(ctx)
	if cfg.Flags.StandAlone {
		return f.Client().DialOrDie().PolicyV1beta1().PodSecurityPolicies().List(ctx, metav1.ListOptions{})
	}

	var res dao.Resource
	res.Init(f, client.NewGVR("policy/v1beta1/podsecuritypolicies"))
	oo, err := res.List(ctx, client.AllNamespaces)
	if err != nil {
		return nil, err
	}
	var ll pv1beta1.PodSecurityPolicyList
	for _, o := range oo {
		var psp pv1beta1.PodSecurityPolicy
		err = runtime.DefaultUnstructuredConverter.FromUnstructured(o.(*unstructured.Unstructured).Object, &psp)
		if err != nil {
			return nil, errors.New("expecting configmap resource")
		}
		ll.Items = append(ll.Items, psp)
	}

	return &ll, nil

}
