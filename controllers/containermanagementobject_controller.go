package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	errors "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	containershipv1beta1 "github.com/relativitydev/containership/api/v1beta1"
	r1azure "github.com/relativitydev/containership/pkg/azure"
	r1slack "github.com/relativitydev/containership/pkg/slack"
)

// Time in seconds for regular looping of CMO configurations
var interval = 600

// ContainerManagementObjectReconciler reconciles a ContainerManagementObject object
type ContainerManagementObjectReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=containership.app,resources=containermanagementobjects,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=containership.app,resources=containermanagementobjects/status,verbs=get;update;patch

// Reconcile is business logic to take on every CMO event
func (r *ContainerManagementObjectReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()

	// Fetch the ContainerManagementObject instance. Determine which version there is.
	instance := &containershipv1beta1.ContainerManagementObject{}

	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, errors.Wrap(err, "Error reading the object - requeue the request")
	}

	for _, instanceSpec := range instance.Spec.Images {
		promotion := r1azure.ImagePromotion{
			SourceImage:      instanceSpec.SourceImage,
			TargetRepository: instanceSpec.TargetRepository,
			SupportedTags:    instanceSpec.SupportedTags,
		}
		for _, destination := range instanceSpec.Destinations.AzureContainerRegistries {
			promotion.Destinations = append(promotion.Destinations, r1azure.PromotionDestination{
				Name:           destination.Name,
				Ring:           destination.Ring,
				SubscriptionID: destination.SubscriptionID,
				ResourceGroup:  destination.ResourceGroup,
			})
		}

		events := r1azure.ImageProcessor(ctx, promotion)

		if len(events) > 0 {
			// publish events created prior to this point in the reconcile function
			publishEvents(r, instance, events)
		}
	}

	return ctrl.Result{RequeueAfter: time.Second * time.Duration(interval)}, nil
}

// SetupWithManager configures how the controller watches CMO
func (r *ContainerManagementObjectReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&containershipv1beta1.ContainerManagementObject{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Complete(r)
}

// There is no need to reconcile deleted CMOs
func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// publishEvents will publish all previously added events to k8s
func publishEvents(r *ContainerManagementObjectReconciler, instance *containershipv1beta1.ContainerManagementObject, events []corev1.Event) {
	slackClient := r1slack.NewClient(instance.Spec.SlackWebhookEndpoint)

	// publish events created prior to this point in the reconcile function
	for _, event := range events {
		r.Recorder.Event(instance, event.Type, event.Reason, event.Message)

		if instance.Spec.SlackWebhookEndpoint != "" {
			// send events to slack too
			err := slackClient.SendMessage(fmt.Sprintf("*%s*:\n%s", event.Type, event.Message))
			if err != nil {
				r.Log.Error(err, "Slack webhook failed")
			}
		}
	}
}
