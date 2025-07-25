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

package v1

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webappv1 "my.domain/guestbook/api/v1"
)

// nolint:unused
// log is for logging in this package.
var guestbooklog = logf.Log.WithName("guestbook-resource")

// SetupGuestbookWebhookWithManager registers the webhook for Guestbook in the manager.
func SetupGuestbookWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&webappv1.Guestbook{}).
		WithValidator(&GuestbookCustomValidator{}).
		WithDefaulter(&GuestbookCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-webapp-my-domain-v1-guestbook,mutating=true,failurePolicy=fail,sideEffects=None,groups=webapp.my.domain,resources=guestbooks,verbs=create;update,versions=v1,name=mguestbook-v1.kb.io,admissionReviewVersions=v1

// GuestbookCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Guestbook when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type GuestbookCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &GuestbookCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Guestbook.
func (d *GuestbookCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	// guestbook, ok := obj.(*webappv1.Guestbook)

	// if !ok {
	// 	return fmt.Errorf("expected an Guestbook object but got %T", obj)
	// }
	// guestbooklog.Info("Defaulting for Guestbook", "name", guestbook.GetName())

	// TODO(user): fill in your defaulting logic.

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-webapp-my-domain-v1-guestbook,mutating=false,failurePolicy=fail,sideEffects=None,groups=webapp.my.domain,resources=guestbooks,verbs=create;update,versions=v1,name=vguestbook-v1.kb.io,admissionReviewVersions=v1

// GuestbookCustomValidator struct is responsible for validating the Guestbook resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type GuestbookCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &GuestbookCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Guestbook.
func (v *GuestbookCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	// guestbook, ok := obj.(*webappv1.Guestbook)
	// if strings.Contains(strings.ToLower(guestbook.Spec.AppName), "random") {
	// 	return nil, fmt.Errorf("spec.appName must not contain the word %T", guestbook.Spec.AppName)
	// }
	// if !ok {
	// 	return nil, fmt.Errorf("expected a Guestbook object but got %T", obj)
	// }
	// guestbooklog.Info("Validation for Guestbook upon creation", "name", guestbook.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Guestbook.
func (v *GuestbookCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	// guestbook, ok := newObj.(*webappv1.Guestbook)
	// if !ok {
	// 	return nil, fmt.Errorf("expected a Guestbook object for the newObj but got %T", newObj)
	// }
	// guestbooklog.Info("Validation for Guestbook upon update", "name", guestbook.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Guestbook.
func (v *GuestbookCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	// guestbook, ok := obj.(*webappv1.Guestbook)
	// if !ok {
	// 	return nil, fmt.Errorf("expected a Guestbook object but got %T", obj)
	// }
	// guestbooklog.Info("Validation for Guestbook upon deletion", "name", guestbook.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
