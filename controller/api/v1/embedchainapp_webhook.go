/*
Copyright 2024.

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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var embedchainapplog = logf.Log.WithName("embedchainapp-resource")

// SetupWebhookWithManager will setup the manager to manage the webhooks
func (r *EmbedchainApp) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-aiapps-learncloudnative-com-v1-embedchainapp,mutating=true,failurePolicy=fail,sideEffects=None,groups=aiapps.learncloudnative.com,resources=embedchainapps,verbs=create;update,versions=v1,name=membedchainapp.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &EmbedchainApp{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *EmbedchainApp) Default() {
	embedchainapplog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-aiapps-learncloudnative-com-v1-embedchainapp,mutating=false,failurePolicy=fail,sideEffects=None,groups=aiapps.learncloudnative.com,resources=embedchainapps,verbs=create;update,versions=v1,name=vembedchainapp.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &EmbedchainApp{}

func (r *EmbedchainApp) validate() error {
	embedchainapplog.Info("validate", "name", r.Name)
	// TODO: implement the validation logic
	// This is where we could check the validity of configRef and secretRef (i.e. make sure they exist)

	// For reference, check this: https://book.kubebuilder.io/cronjob-tutorial/webhook-implementation
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *EmbedchainApp) ValidateCreate() (admission.Warnings, error) {
	embedchainapplog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil, r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *EmbedchainApp) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	embedchainapplog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil, r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *EmbedchainApp) ValidateDelete() (admission.Warnings, error) {
	embedchainapplog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil, nil
}
