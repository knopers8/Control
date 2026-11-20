/*
 * === This file is part of ALICE O² ===
 *
 * Copyright 2026 CERN and copyright holders of ALICE O².
 * Author: Michal Tichak <michal.tichak@cern.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * In applying this license CERN does not waive the privileges and
 * immunities granted to it by virtue of its status as an
 * Intergovernmental Organization or submit itself to any jurisdiction.
 */

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aliecsv1alpha1 "github.com/AliceO2Group/Control/control-operator/api/v1alpha1"
)

var _ = Describe("Environment Controller with shared Pods", Ordered, func() {
	const (
		namespace = "default"
		envName   = "shared-env"
		nodeName  = "shared-node"
		// the cluster knows the FQDN, the Environment refers to the short name
		clusterNode = "shared-node.cern.ch"
		podName     = "aliecs-pod-shared-env-shared-node-g"
	)

	var (
		ctx        context.Context
		envKey     types.NamespacedName
		reconciler *EnvironmentReconciler
	)

	taskDef := func(name, group string) aliecsv1alpha1.TaskDefinition {
		return aliecsv1alpha1.TaskDefinition{
			Name:     name,
			PodGroup: group,
			Spec: aliecsv1alpha1.TaskTemplateSpec{
				Pod: v1.PodSpec{
					HostNetwork: true,
					Containers:  []v1.Container{{Name: "main", Image: "img-" + name}},
				},
			},
		}
	}

	reconcileOnce := func() {
		GinkgoHelper()
		_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: envKey})
		Expect(err).NotTo(HaveOccurred())
	}

	getTask := func(name string) *aliecsv1alpha1.Task {
		GinkgoHelper()
		task := &aliecsv1alpha1.Task{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, task)).To(Succeed())
		return task
	}

	getPod := func() *v1.Pod {
		GinkgoHelper()
		pod := &v1.Pod{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: podName, Namespace: namespace}, pod)).To(Succeed())
		return pod
	}

	BeforeAll(func() {
		ctx = context.Background()
		envKey = types.NamespacedName{Name: envName, Namespace: namespace}
		reconciler = &EnvironmentReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}

		Expect(k8sClient.Create(ctx, &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: clusterNode}})).To(Succeed())
		Expect(k8sClient.Create(ctx, &aliecsv1alpha1.Environment{
			ObjectMeta:    metav1.ObjectMeta{Name: envName, Namespace: namespace},
			TaskTemplates: aliecsv1alpha1.TemplateSpecification{Tasks: map[string][]aliecsv1alpha1.TaskReference{}},
			Spec: aliecsv1alpha1.EnvironmentSpec{
				State: "standby",
				Tasks: map[string][]aliecsv1alpha1.TaskDefinition{
					nodeName: {taskDef("a", "g"), taskDef("b", "g"), taskDef("c", "")},
				},
			},
		})).To(Succeed())

		DeferCleanup(func() {
			// only clears what a failed spec left behind
			pod := &v1.Pod{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Name: podName, Namespace: namespace}, pod); err == nil {
				Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, pod, client.GracePeriodSeconds(0)))).To(Succeed())
			}
			tasks := &aliecsv1alpha1.TaskList{}
			Expect(k8sClient.List(ctx, tasks, client.InNamespace(namespace), client.MatchingLabels{"environment": envName})).To(Succeed())
			for i := range tasks.Items {
				controllerutil.RemoveFinalizer(&tasks.Items[i], taskFinalizer)
				Expect(k8sClient.Update(ctx, &tasks.Items[i])).To(Succeed())
				Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, &tasks.Items[i]))).To(Succeed())
			}
			env := &aliecsv1alpha1.Environment{}
			if err := k8sClient.Get(ctx, envKey, env); err == nil {
				controllerutil.RemoveFinalizer(env, environmentFinalizer)
				Expect(k8sClient.Update(ctx, env)).To(Succeed())
				Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, env))).To(Succeed())
			}
			Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: clusterNode}}))).To(Succeed())
		})
	})

	It("creates one Pod per group and a separate Pod for ungrouped Tasks", func() {
		reconcileOnce() // finalizer
		reconcileOnce() // deployment

		pod := getPod()
		Expect(pod.Spec.Containers).To(HaveLen(2))
		Expect(pod.Spec.Containers[0].Name).To(Equal("a"))
		Expect(pod.Spec.Containers[1].Name).To(Equal("b"))
		Expect(pod.Spec.NodeName).To(Equal(clusterNode))
		Expect(pod.Spec.RestartPolicy).To(Equal(v1.RestartPolicyNever))
		Expect(pod.OwnerReferences).To(HaveLen(1))
		Expect(pod.OwnerReferences[0].Kind).To(Equal("Environment"))

		for _, name := range []string{"a", "b"} {
			task := getTask(nodeName + "-" + name)
			Expect(task.Spec.PodName).To(Equal(podName))
			Expect(task.Spec.ContainerName).To(Equal(name))
			Expect(task.OwnerReferences).To(HaveLen(1))
			Expect(task.OwnerReferences[0].Kind).To(Equal("Environment"))
		}

		ungrouped := getTask(nodeName + "-c")
		Expect(ungrouped.Spec.PodName).To(BeEmpty(), "ungrouped Tasks keep getting their own Pod from the Task controller")
	})

	It("deletes the shared Pod only after all Tasks are gone", func() {
		By("keeping one Task alive with a finalizer, as the Task controller would while it tears OCC down")
		a := getTask(nodeName + "-a")
		controllerutil.AddFinalizer(a, taskFinalizer)
		Expect(k8sClient.Update(ctx, a)).To(Succeed())

		env := &aliecsv1alpha1.Environment{}
		Expect(k8sClient.Get(ctx, envKey, env)).To(Succeed())
		Expect(k8sClient.Delete(ctx, env)).To(Succeed())

		reconcileOnce()

		Expect(getTask(nodeName + "-a").DeletionTimestamp.IsZero()).To(BeFalse())
		Expect(getPod().DeletionTimestamp.IsZero()).To(BeTrue(), "the Pod must outlive the Tasks tearing down")

		By("the Task controller finishing its teardown")
		a = getTask(nodeName + "-a")
		controllerutil.RemoveFinalizer(a, taskFinalizer)
		Expect(k8sClient.Update(ctx, a)).To(Succeed())

		reconcileOnce()

		Expect(getPod().DeletionTimestamp.IsZero()).To(BeFalse(), "all Tasks are gone, the Pod can be deleted")
		Expect(k8sClient.Get(ctx, envKey, &aliecsv1alpha1.Environment{})).To(Succeed(),
			"the Environment waits until the Pod is gone")

		By("the kubelet confirming the Pod termination")
		Expect(k8sClient.Delete(ctx, getPod(), client.GracePeriodSeconds(0))).To(Succeed())

		reconcileOnce()

		Expect(k8sClient.Get(ctx, envKey, &aliecsv1alpha1.Environment{})).To(Satisfy(errors.IsNotFound))
	})
})
