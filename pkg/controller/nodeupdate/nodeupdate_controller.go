package nodeupdate

import (
	"context"
	"fmt"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	powervsclient "github.com/openshift/cluster-api-provider-powervs/pkg/client"
)

const (
	requeueDurationWhenVMNotReady = 1 * time.Minute
	requeueDurationWhenVMNotFound = 2 * time.Minute
)

var _ reconcile.Reconciler = &providerIDReconciler{}

type providerIDReconciler struct {
	client               client.Client
	PowerVSClientBuilder powervsclient.PowerVSClientBuilderFuncType
}

// Reconcile make sure a node has a ProviderID set. The providerID is the ID
// of the machine on powervs. The ID is the PvmInstanceID
// as its guarantee to be unique in a cluster.
func (r *providerIDReconciler) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	klog.Infof("%s: Reconciling node", request.NamespacedName)

	// Fetch the Node instance
	node := corev1.Node{}
	err := r.client.Get(context.Background(), request.NamespacedName, &node)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			klog.Infof("%s: Node not found - do nothing", request.NamespacedName)
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, fmt.Errorf("error getting node: %v", err)
	}

	if node.Spec.ProviderID != "" {
		return reconcile.Result{}, nil
	}

	klog.Infof("%s: ProviderID is not updated in the node - update it", node.Name)

	apiKey, err := powervsclient.GetAPIKey(r.client, powervsclient.DefaultCredentialSecret, powervsclient.DefaultCredentialNamespace)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to read the API key from the secret: %v", err)
	}

	c, err := powervsclient.NewClientMinimal(apiKey)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("%s: failed to create NewClientMinimal, with error: %v", node.Name, err)
	}

	instances, err := c.GetCloudInstances()
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("%s: failed to get the cloud instances, with error: %v", node.Name, err)
	}

	var instance *models.PVMInstanceReference
	// TODO: Need to add concurrency
out:
	for _, i := range instances {
		cl, err := powervsclient.NewValidatedClient(r.client, powervsclient.DefaultCredentialSecret, powervsclient.DefaultCredentialNamespace, i.Guid, "")
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("%s: failed to create NewValidatedClient, with error: %v", node.Name, err)
		}

		ins, err := cl.GetInstances()
		if ins == nil {
			return reconcile.Result{}, fmt.Errorf("%s: failed to GetInstances for %s, with error: %v", node.Name, i.Name, err)
		}

		for _, in := range ins.PvmInstances {
			if *in.ServerName == node.Name {
				instance = in
				break out
			}
		}
	}

	if instance != nil {
		node.Spec.ProviderID = powervsclient.FormatProviderID(*instance.PvmInstanceID)
	} else {
		// TODO: enable this block later
		//klog.Infof("%s: Virtual Machine of this node doesn't exists - delete the node", node.Name)
		//if err := r.client.Delete(context.Background(), &node); err != nil {
		//	return reconcile.Result{}, fmt.Errorf("%s: Error deleting Node, with error: %v", node.Name, err)
		//}
		//return reconcile.Result{}, nil
		klog.Infof("%s: Virtual Machine of this node doesn't exists, requeuing after 2 mins", node.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDurationWhenVMNotFound}, nil
	}

	if *instance.Status != powervsclient.InstanceStateNameActive {
		klog.Infof("%s: Virtual Machine of this node isn't ready - requeue for 1 minute", node.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: requeueDurationWhenVMNotReady}, nil
	}

	if err = r.client.Update(context.Background(), &node); err != nil {
		return reconcile.Result{}, fmt.Errorf("%s: failed updating node, with error: %v", node.Name, err)
	}

	return reconcile.Result{}, nil
}

// Add registers a new provider ID reconciler controller with the controller manager
func Add(mgr manager.Manager) error {
	reconciler, err := NewProviderIDReconciler(mgr)

	if err != nil {
		return fmt.Errorf("error building reconciler: %v", err)
	}

	c, err := controller.New("providerID-controller", mgr, controller.Options{Reconciler: reconciler})
	if err != nil {
		return err
	}

	//Watch node changes
	err = c.Watch(&source.Kind{Type: &corev1.Node{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// NewProviderIDReconciler creates a new providerID reconciler
func NewProviderIDReconciler(mgr manager.Manager) (*providerIDReconciler, error) {
	r := providerIDReconciler{
		client: mgr.GetClient(),
	}
	return &r, nil
}
