package controller

import (
	"context"

	k8sBatchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	batchv1 "my.domain/api/v1"
)

// ClusterScanReconciler reconciles a ClusterScan object
type ClusterScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ClusterScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the ClusterScan instance
	clusterScan := &batchv1.ClusterScan{}
	err := r.Get(ctx, req.NamespacedName, clusterScan)
	if err != nil {
		log.Error(err, "unable to fetch ClusterScan")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Define a new Job object
	job := newJobForCustomResource(clusterScan)

	// Check if this Job already exists
	found := &batchv1.ClusterScan{}
	err = r.Get(ctx, types.NamespacedName{Name: job.Name, Namespace: job.Namespace}, found)
	if err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch Job")
			return ctrl.Result{}, err
		}

		log.Info("Creating a new Job", "Job.Namespace", job.Namespace, "Job.Name", job.Name)
		err = r.Create(ctx, job)
		if err != nil {
			log.Error(err, "unable to create Job")
			return ctrl.Result{}, err
		}

		// Job created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Job already exists - don't requeue
	log.Info("Skip reconcile: Job already exists", "Job.Namespace", found.Namespace, "Job.Name", found.Name)
	return ctrl.Result{}, nil
}

// newJobForCustomResource returns a job with the same name/namespace as the cr
func newJobForCustomResource(cr *batchv1.ClusterScan) *batchv1.ClusterScan {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &batchv1.ClusterScan{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-job",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: batchv1.ClusterScanSpec{
			JobTemplate: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "kube-linter",
							Image: "stackrox/kube-linter:0.2.2",
							Args:  []string{"lint", "../../example-files"},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "dir-to-lint",
									MountPath: "../../example-files",
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
					Volumes: []corev1.Volume{
						{
							Name: "dir-to-lint",
							VolumeSource: corev1.VolumeSource{
								HostPath: &corev1.HostPathVolumeSource{
									Path: "../../example-files",
								},
							},
						},
					},
				},
			},
			CronJobTemplate: batchv1beta1.CronJobSpec{
				Schedule: "*/1 * * * *",
				JobTemplate: batchv1beta1.JobTemplateSpec{
					Spec: k8sBatchv1.JobSpec{
						Template: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  "kube-linter",
										Image: "stackrox/kube-linter:0.2.2",
										Args:  []string{"lint", "../../example-files"},
										VolumeMounts: []corev1.VolumeMount{
											{
												Name:      "dir-to-lint",
												MountPath: "../../example-files",
											},
										},
									},
								},
								RestartPolicy: corev1.RestartPolicyOnFailure,
								Volumes: []corev1.Volume{
									{
										Name: "dir-to-lint",
										VolumeSource: corev1.VolumeSource{
											HostPath: &corev1.HostPathVolumeSource{
												Path: "../../example-files",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.ClusterScan{}).
		Complete(r)
}
