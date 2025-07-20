/*
Copyright 2025.

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

package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "my.domain/guestbook/api/v1"
)

const guestbookFinalizer = "webapp.my.domain/finalizer"

// GuestbookReconciler reconciles a Guestbook object
type GuestbookReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=webapp.my.domain,resources=guestbooks,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.my.domain,resources=guestbooks/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=webapp.my.domain,resources=guestbooks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Guestbook object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *GuestbookReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := logf.FromContext(ctx)
	url := "http://localhost:8080/hello-world-with-body"

	// Get Guestbook instance
	var guestbook webappv1.Guestbook
	if err := r.Get(ctx, req.NamespacedName, &guestbook); err != nil {
		if apierrors.IsNotFound(err) {
			// Deleted resource (finalizer already removed, nothing to do)
			logger.Info("Guestbook deleted (object no longer exists)")
			return ctrl.Result{}, nil // Resource deleted
		}
		return ctrl.Result{}, err
	}
	data := map[string]string{
		"appName": guestbook.Spec.AppName,
		"domain":  guestbook.Spec.Domain,
	}

	// -------------------- DELETE HANDLING --------------------
	if !guestbook.ObjectMeta.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(&guestbook, guestbookFinalizer) {
			jsonData, err := json.Marshal(data)
			if err != nil {
				logger.Error(err, "JSON marshal error", "Error:", err)
				return ctrl.Result{}, err
			}

			req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
			if err != nil {
				logger.Error(err, "Error creating request", "Error:", err)
				return ctrl.Result{}, err
			}
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				logger.Error(err, "DELETE Request Error", "Error:", err)
				return ctrl.Result{}, err
			}

			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				logger.Error(err, "Read error", "Error:", err)
				return ctrl.Result{}, err
			}
			logger.Info("Guestbook is being deleted", "Custom Resource Name:", guestbook.Name, "Namespace:", guestbook.Namespace, "Response Body:", string(body))

			// Remove finalizer so deletion can complete
			controllerutil.RemoveFinalizer(&guestbook, guestbookFinalizer)
			if err := r.Update(ctx, &guestbook); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// -------------------- CREATE HANDLING --------------------
	if !controllerutil.ContainsFinalizer(&guestbook, guestbookFinalizer) {
		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Error(err, "JSON marshal error", "Error:", err)
			return ctrl.Result{}, err
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error(err, "POST Request Error", "Error:", err)
			return ctrl.Result{}, err
		}
		defer resp.Request.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Read error:", err)
			return ctrl.Result{}, err
		}
		logger.Info("Guestbook is being created", "Custom Resource Name:", guestbook.Name, "Namespace:", guestbook.Namespace, "Response Body:", string(body))

		// Add the finalizer so we can handle deletion events later
		controllerutil.AddFinalizer(&guestbook, guestbookFinalizer)
		// Save the finalizer update to the API server
		if err := r.Update(ctx, &guestbook); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	// -------------------- UPDATE HANDLING --------------------
	if guestbook.Generation != guestbook.Status.ObservedGeneration {
		jsonData, err := json.Marshal(data)
		if err != nil {
			logger.Error(err, "JSON marshal error", "Error:", err)
			return ctrl.Result{}, err
		}

		req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error(err, "Error creating request", "Error:", err)
			return ctrl.Result{}, err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			logger.Error(err, "PUT Request Error", "Error:", err)
			return ctrl.Result{}, err
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err, "Read error", "Error:", err)
			return ctrl.Result{}, err
		}
		logger.Info("Guestbook is being updated", "Custom Resource Name:", guestbook.Name, "Namespace:", guestbook.Namespace, "Response Body:", string(body))

		guestbook.Status.ObservedGeneration = guestbook.Generation
		if err := r.Status().Update(ctx, &guestbook); err != nil {
			logger.Error(err, "Failed to update status.observedGeneration")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GuestbookReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.Guestbook{}).
		Named("guestbook").
		Complete(r)
}
