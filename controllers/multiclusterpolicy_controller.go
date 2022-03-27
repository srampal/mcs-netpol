/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"

	"github.com/coreos/go-iptables/iptables"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	//	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"

	policyv1alpha1 "github.com/srampal/mcs-netpol/api/v1alpha1"
	mcsv1a1 "sigs.k8s.io/mcs-api/pkg/apis/v1alpha1"
)

// MultiClusterPolicyReconciler reconciles a MultiClusterPolicy object
type MultiClusterPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	ipt *iptables.IPTables
}

//+kubebuilder:rbac:groups=policy.submariner.io,resources=multiclusterpolicies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=policy.submariner.io,resources=multiclusterpolicies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=policy.submariner.io,resources=multiclusterpolicies/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=pods/status,verbs=get;list;watch
//+kubebuilder:rbac:groups="multicluster.x-k8s.io",resources=serviceimports,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MultiClusterPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MultiClusterPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	loggr := log.FromContext(ctx)

	loggr.Info("\n Entered MCS-NETPOL reconciler \n")

	var mcsPol policyv1alpha1.MultiClusterPolicy

	if err := r.Get(ctx, req.NamespacedName, &mcsPol); err != nil {
		loggr.Error(err, "unable to fetch McsPol")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// mcsPolFinalizerName := "multiclusterpolicies.submariner.io/finalizer"

	if mcsPol.ObjectMeta.DeletionTimestamp.IsZero() {

		loggr.Info("\n Found DeletionTimestamp Zero! \n")

		/**
		      Removing finalizer logic for now

		  		// Register the finalizer and update the object
		  		if !controllerutil.ContainsFinalizer(&mcsPol, mcsPolFinalizerName) {
		  			controllerutil.AddFinalizer(&mcsPol, mcsPolFinalizerName)
		  		}
		  		loggr.Info("\n Registering Finalizer \n")
		  		if err := r.Update(ctx, &mcsPol); err != nil {
		  			loggr.Info("\n Error registering finalizer into mcsPol! \n")
		  			return ctrl.Result{}, err
		  		}
		  **/
	} else {
		if err := r.handleDeletion(ctx, req, &mcsPol); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
		/**
				if controllerutil.ContainsFinalizer(&mcsPol, mcsPolFinalizerName) {

					if err := r.handleDeletion(ctx, req, &mcsPol); err != nil {
						return ctrl.Result{}, err
					}
				}

				controllerutil.RemoveFinalizer(&mcsPol, mcsPolFinalizerName)
				if err := r.Update(ctx, &mcsPol); err != nil {
					loggr.Info("\n Error removing finalizer from mcsPol! \n")
					return ctrl.Result{}, err
				}
		**/
	}

	if err := r.handleCreateOrUpdate(ctx, req, &mcsPol); err != nil {
		return ctrl.Result{}, err
	}

	err1 := errors.New("MCS-NETPOL reconciler() not yet fully implemented")
	loggr.Error(err1, "Exitting ..")
	return ctrl.Result{}, nil

}

func (r *MultiClusterPolicyReconciler) handleCreateOrUpdate(ctx context.Context, req ctrl.Request,
	mcsPol *policyv1alpha1.MultiClusterPolicy) error {

	loggr := log.FromContext(ctx)

	loggr.Info("\n Entered MCS-NETPOL handleCreateOrUpdate()  \n")

	ipt, err := r.mcsPolInitIptablesChain(ctx)

	if err != nil {
		loggr.Error(err, "Could not create IPTables chain")
		return err
	}
	fmt.Printf("\n Created IPTables %+v\n", ipt)

	if ipt != nil {
		fmt.Printf("\n Created IPTables %#v\n", *ipt)
	}

	// Get list of all mcsPols

	var mcsPols policyv1alpha1.MultiClusterPolicyList

	err = r.List(ctx, &mcsPols, client.InNamespace(req.NamespacedName.Namespace))

	if err != nil {

		fmt.Printf("\n Could not retrieve policy list! \n")
		return errors.Wrap(err, "could not retrieve policy list")
	}

	fmt.Printf("Got policy list! \n")

	for i, mcsPol := range mcsPols.Items {
		fmt.Printf("Got policy %d)  %#v\n", i, mcsPol)
		if err1 := r.programPolicyInDataplane(ctx, req, &mcsPol); err1 != nil {
			loggr.Error(err, "Error programming policy !")
			return err1
		}
	}

	err1 := errors.New("MCS-NETPOL handleCreateOrUpdate() not yet fully implemented")
	loggr.Error(err1, "Exitting ..")
	return nil
}

func (r *MultiClusterPolicyReconciler) programPolicyInDataplane(ctx context.Context, req ctrl.Request,
	mcsPol *policyv1alpha1.MultiClusterPolicy) error {

	loggr := log.FromContext(ctx)

	loggr.Info("\n Entered programPolicyInDataplane()  \n")

	// Get the list of pods selected by this policy
	selector, err := metav1.LabelSelectorAsSelector(&mcsPol.Spec.PodSelector)
	if err != nil {
		loggr.Error(err, "Error creating pod label selector in programPolicyinDataplane\n")
		return errors.Wrap(err, "error creating pod label selector")
	}

	pods := &corev1.PodList{}

	err = r.List(ctx, pods, client.InNamespace(req.NamespacedName.Namespace), client.MatchingLabelsSelector{Selector: selector})
	if err != nil {
		loggr.Error(err, "Error getting selected pod list")
		return errors.Wrap(err, "error getting selected pod list")
	}

	loggr.Info("\n Got the pod list \n")
	fmt.Printf("\n Got the pod list %d items \n", len(pods.Items))

	// Get the list of policy peers, rules and service imports in this policy
	// and add the corresponding iptable rule

	serviceImport := &mcsv1a1.ServiceImport{}
	var siName types.NamespacedName
	// Hard coded for now, will be programmable eventually
	siName.Namespace = "submariner-operator"

	loggr.Info("\n Policy Rules & peers ... \n")

	for i, rule := range mcsPol.Spec.Egress {
		fmt.Printf("\n Got rule %d)\n", i)
		for j, peer := range rule.To {
			fmt.Printf("\n Got peer %d)\n", j)
			for k, si := range peer.ServiceImportRefs {
				fmt.Printf("\n Got serviceImport %d)    %#v\n", k, si)

				/* For each ServiceImport, program rules for each source pod in the list */
				// create namespaced name for the service import

				siName.Name = si
				if err := r.Get(ctx, siName, serviceImport); err != nil {
					fmt.Printf("\n Could not retrieve si (%s, %s) \n",
						siName.Namespace, siName.Name)
					loggr.Error(err, "unable to fetch ServiceImport Object")
					return err
				} else {
					fmt.Printf("\n Got ServiceImport object %#v \n", serviceImport)
					//program iptables rules for each pod and this service import
					// Skipping validation checks for now
					dstIp := serviceImport.Spec.IPs[0]
					for l, pod := range pods.Items {

						if pod.Status.Phase != "Running" {
							fmt.Printf("\n Skipping pod  no. %d not in running state\n", l)
							continue
						}

						srcIp := pod.Status.PodIP

						if err := r.ipt.Insert("filter", "MCS-OUTPUT", 1,
							"-s", srcIp, "-d", dstIp, "-j", "DROP"); err != nil {
							loggr.Error(err, "Error inserting rule in MCS-OUTPUT chain\n")
							return errors.Wrap(err, "unable to insert Jump in OUTPUT chain")
						} else {
							fmt.Printf("\n Success programming iptables rule s %s d %s\n", srcIp, dstIp)
						}

					}

				}
			}
		}
	}

	// For each pod and each policy peer, program an iptables rules

	err1 := errors.New("programPolicyInDataplane() not yet fully implemented")
	loggr.Error(err1, "Exitting ..")
	return nil

}

func (r *MultiClusterPolicyReconciler) handleDeletion(ctx context.Context, req ctrl.Request,
	mcsPol *policyv1alpha1.MultiClusterPolicy) error {

	loggr := log.FromContext(ctx)

	loggr.Info("\n Entered MCS-NETPOL handleDeletion()  \n")

	err1 := errors.New("MCS-NETPOL handleDeletion() not yet fully implemented")
	loggr.Error(err1, "Exitting ..")
	return nil
}

func (r *MultiClusterPolicyReconciler) mcsPolInitIptablesChain(ctx context.Context) (*iptables.IPTables, error) {

	loggr := log.FromContext(ctx)

	loggr.Info("\n Entered mcsPolInitIpTablesChain \n")

	if r.ipt != nil {
		return r.ipt, nil
	}

	ipt, err := iptables.New(iptables.IPFamily(iptables.ProtocolIPv4), iptables.Timeout(5))

	if err != nil {
		return nil, errors.Wrap(err, "error creating IP tables")
	}

	/* Create new MCS chain */

	err = initMcsOutputChain(ctx, ipt)

	if err != nil {
		return nil, errors.Wrap(err, "error creating MCS chain")
	}

	r.ipt = ipt

	return ipt, nil
}

func initMcsOutputChain(ctx context.Context, ipt *iptables.IPTables) error {

	loggr := log.FromContext(ctx)

	loggr.Info("\n Entered initMcsOutputChain \n")

	if err := createChainIfNotExists(ctx, ipt, "filter", "MCS-OUTPUT"); err != nil {
		loggr.Error(err, "Error in createChainIfNotExists() \n")
		return errors.Wrap(err, "unable to create MCS-OUTPUT chain in iptables")
	}

	loggr.Info("\n Created new chain MCS-OUTPUT\n")

	forwardToMcsOutputRuleSpec := []string{"-j", "MCS-OUTPUT"}

	if err := ipt.Insert("filter", "OUTPUT", 1, forwardToMcsOutputRuleSpec...); err != nil {
		loggr.Error(err, "Error inserting Jump to MCS-OUTPUT chain\n")
		return errors.Wrap(err, "unable to insert Jump in OUTPUT chain")
	}

	loggr.Info("\n Added Jump in OUTPUT chain to MCS-OUTPUT\n")

	return nil
}

func createChainIfNotExists(ctx context.Context, ipt *iptables.IPTables, table, chain string) error {

	loggr := log.FromContext(ctx)

	existingChains, err := ipt.ListChains(table)

	if err != nil {
		loggr.Error(err, "Error ipt.ListChains \n")
		return errors.Wrap(err, "error listing IP table chains")
	}

	for _, val := range existingChains {
		if val == chain {
			// Chain already exists
			fmt.Printf("\n Chain already exists \n")
			return nil
		}

	}

	return errors.Wrap(ipt.NewChain(table, chain), "error creating IP table chain")

	/*
	           err = ipt.NewChain(table, chain)

	           if err != nil {
	   	    return errors.New("error creating IP table chain")
	           } else {
	               return nil
	           }
	*/
}

// SetupWithManager sets up the controller with the Manager.
func (r *MultiClusterPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&policyv1alpha1.MultiClusterPolicy{}).
		Complete(r)
}
