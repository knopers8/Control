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
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aliecsv1alpha1 "github.com/AliceO2Group/Control/control-operator/api/v1alpha1"
)

var _ = Describe("Task Controller with a shared Pod", Ordered, func() {
	const (
		taskResource  = "shared-task-test"
		sharedPod     = "aliecs-pod-shared-test"
		namespace     = "default"
		containerName = "mine"
		containerID   = "containerd://abc123"
		podIP         = "127.0.0.1"
	)

	var (
		ctx        context.Context
		taskName   types.NamespacedName
		podKey     types.NamespacedName
		occServer  *occServerStub
		reconciler *TaskReconciler
	)

	reconcileOnce := func() {
		GinkgoHelper()
		_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: taskName})
		Expect(err).NotTo(HaveOccurred())
	}

	getTask := func() *aliecsv1alpha1.Task {
		GinkgoHelper()
		task := &aliecsv1alpha1.Task{}
		Expect(k8sClient.Get(ctx, taskName, task)).To(Succeed())
		return task
	}

	BeforeAll(func() {
		ctx = context.Background()
		taskName = types.NamespacedName{Name: taskResource, Namespace: namespace}
		podKey = types.NamespacedName{Name: sharedPod, Namespace: namespace}

		stub, port, stop, err := startOccServerStub()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(stop)
		occServer = stub
		occServer.setState("standby")

		reconciler = &TaskReconciler{Client: k8sClient, Scheme: k8sClient.Scheme(), Recorder: &record.FakeRecorder{}}

		Expect(k8sClient.Create(ctx, &aliecsv1alpha1.Task{
			ObjectMeta: metav1.ObjectMeta{Name: taskResource, Namespace: namespace},
			Spec: aliecsv1alpha1.TaskSpec{
				Control:       aliecsv1alpha1.TaskSpecControl{Mode: "direct", Port: port},
				State:         "standby",
				PodName:       sharedPod,
				ContainerName: containerName,
				Pod:           v1.PodSpec{Containers: []v1.Container{{Name: containerName, Image: "img"}}},
			},
		})).To(Succeed())

		DeferCleanup(func() {
			if c, loaded := clientsForContainers.LoadAndDelete(taskResource); loaded {
				Expect(c.(*OccClient).Close()).To(Succeed())
			}
			pod := &v1.Pod{}
			if err := k8sClient.Get(ctx, podKey, pod); err == nil {
				Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, pod, client.GracePeriodSeconds(0)))).To(Succeed())
			}
			leftover := &aliecsv1alpha1.Task{}
			if err := k8sClient.Get(ctx, taskName, leftover); err == nil {
				controllerutil.RemoveFinalizer(leftover, taskFinalizer)
				Expect(k8sClient.Update(ctx, leftover)).To(Succeed())
				Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, leftover))).To(Succeed())
			}
		})
	})

	It("does not create a Pod of its own while the shared Pod is missing", func() {
		reconcileOnce() // finalizer
		reconcileOnce()

		Expect(k8sClient.Get(ctx, podKey, &v1.Pod{})).To(Satisfy(errors.IsNotFound))
		Expect(k8sClient.Get(ctx, types.NamespacedName{Name: podNameFromTask(taskResource), Namespace: namespace}, &v1.Pod{})).
			To(Satisfy(errors.IsNotFound))
	})

	It("waits for its own container to run", func() {
		Expect(k8sClient.Create(ctx, &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{Name: sharedPod, Namespace: namespace, Labels: map[string]string{labelSharedPod: "true"}},
			Spec: v1.PodSpec{Containers: []v1.Container{
				{Name: containerName, Image: "img"},
				{Name: "neighbour", Image: "img"},
			}},
		})).To(Succeed())

		pod := &v1.Pod{}
		Expect(k8sClient.Get(ctx, podKey, pod)).To(Succeed())
		pod.Status.Phase = v1.PodRunning
		pod.Status.PodIP = podIP
		// our container is not started, while the neighbour is
		pod.Status.ContainerStatuses = []v1.ContainerStatus{
			{Name: "neighbour", State: v1.ContainerState{Running: &v1.ContainerStateRunning{}}},
		}
		Expect(k8sClient.Status().Update(ctx, pod)).To(Succeed())

		reconcileOnce()

		_, connected := clientsForContainers.Load(taskResource)
		Expect(connected).To(BeFalse())
	})

	It("records the container ID and connects once its container is running", func() {
		pod := &v1.Pod{}
		Expect(k8sClient.Get(ctx, podKey, pod)).To(Succeed())
		pod.Status.ContainerStatuses = append(pod.Status.ContainerStatuses, v1.ContainerStatus{
			Name: containerName, ContainerID: containerID, State: v1.ContainerState{Running: &v1.ContainerStateRunning{}},
		})
		Expect(k8sClient.Status().Update(ctx, pod)).To(Succeed())

		reconcileOnce()
		Expect(getTask().Status.ContainerID).To(Equal(containerID))

		reconcileOnce()
		occClient, connected := clientsForContainers.Load(taskResource)
		Expect(connected).To(BeTrue())
		Expect(meta.IsStatusConditionTrue(getTask().Status.Conditions, aliecsv1alpha1.ConditionGRPCConnected)).To(BeTrue())

		connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		Expect(occClient.(*OccClient).WaitUntilConnected(connectCtx)).To(Succeed())

		reconcileOnce() // reads the state
		Expect(getTask().Status.State).To(Equal("standby"))
	})

	It("transitions the container to a requested state", func() {
		task := getTask()
		task.Spec.State = "configured"
		Expect(k8sClient.Update(ctx, task)).To(Succeed())

		reconcileOnce()

		Expect(getTask().Status.State).To(Equal("configured"))
		Expect(occServer.getState()).To(Equal("configured"))
	})

	It("steps the process down via OCC and leaves the shared Pod alone when deleted", func() {
		Expect(k8sClient.Delete(ctx, getTask())).To(Succeed())

		reconcileOnce()

		Expect(occServer.getState()).To(Equal("standby"))
		_, connected := clientsForContainers.Load(taskResource)
		Expect(connected).To(BeFalse())
		Expect(k8sClient.Get(ctx, taskName, &aliecsv1alpha1.Task{})).To(Satisfy(errors.IsNotFound))

		pod := &v1.Pod{}
		Expect(k8sClient.Get(ctx, podKey, pod)).To(Succeed())
		Expect(pod.DeletionTimestamp.IsZero()).To(BeTrue(), "the Environment controller owns the shared Pod")
	})
})
