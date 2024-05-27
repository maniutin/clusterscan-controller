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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	batchv1beta1 "k8s.io/api/batch/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sBatchv1 "k8s.io/api/batch/v1"
	batchv1 "my.domain/api/v1"
)

var _ = Describe("ClusterScan Controller", func() {
	Context("When reconciling a resource", func() {
		resourceName := "test-scan-3"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default",
		}
		clusterscan := &batchv1.ClusterScan{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind ClusterScan")

			err := k8sClient.Get(ctx, typeNamespacedName, clusterscan)
			if err != nil || errors.IsNotFound(err) {
				// fmt.Printf("IsNotFound: %v", err)
				resource := &batchv1.ClusterScan{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					TypeMeta: metav1.TypeMeta{
						APIVersion: "batch.my.domain/v1",
						Kind:       "ClusterScan",
					},
					Spec: batchv1.ClusterScanSpec{
						JobTemplate: corev1.PodTemplateSpec{
							Spec: corev1.PodSpec{
								Containers: []corev1.Container{
									{
										Name:  "test-container",
										Image: "busybox",
										Command: []string{
											"echo",
											"hello",
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
													Args:  []string{"lint", "../../files-to-lint"},
													VolumeMounts: []corev1.VolumeMount{
														{
															Name:      "dir-to-lint",
															MountPath: "../../files-to-lint",
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
															Path: "../../files-to-lint",
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Schedule: "*/1 * * * *",
					},
					Status: batchv1.ClusterScanStatus{
						JobStatus: k8sBatchv1.JobStatus{
							Active: 1,
						},
					},
				}

				By("Creating the ClusterScan resource")

				err := k8sClient.Create(ctx, resource)
				Expect(err).NotTo(HaveOccurred())

			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &batchv1.ClusterScan{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance ClusterScan")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})

		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &ClusterScanReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
