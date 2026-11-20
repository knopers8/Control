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
	"strings"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	aliecsv1alpha1 "github.com/AliceO2Group/Control/control-operator/api/v1alpha1"
)

func newTestEnvironment() *aliecsv1alpha1.Environment {
	return &aliecsv1alpha1.Environment{
		ObjectMeta: metav1.ObjectMeta{Name: "env1", Namespace: "alice-tasks", UID: types.UID("env-uid")},
	}
}

func testPendingTask(name, container string, mutate func(*v1.PodSpec)) pendingTask {
	spec := v1.PodSpec{
		NodeName:    "node1",
		HostNetwork: true,
		Containers:  []v1.Container{{Name: "orig", Image: "img-" + container}},
		Volumes:     []v1.Volume{{Name: "tmp", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/tmp"}}}},
	}
	if mutate != nil {
		mutate(&spec)
	}
	return pendingTask{
		task: &aliecsv1alpha1.Task{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "alice-tasks", Labels: map[string]string{"environment": "env1"}},
			Spec:       aliecsv1alpha1.TaskSpec{Pod: spec},
		},
		containerName: container,
	}
}

func TestBuildSharedPod(t *testing.T) {
	key := podKey{node: "node1", resolvedNode: "node1", group: "dataflow"}
	a := testPendingTask("node1-readout", "readout", func(s *v1.PodSpec) { s.TerminationGracePeriodSeconds = ptr.To(int64(10)) })
	b := testPendingTask("node1-stfb", "stfb", func(s *v1.PodSpec) {
		s.TerminationGracePeriodSeconds = ptr.To(int64(30))
		s.Volumes = append(s.Volumes, v1.Volume{Name: "extra"})
	})

	pod, err := buildSharedPod(newTestEnvironment(), key, []pendingTask{a, b}, newTestScheme(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pod.Name != "aliecs-pod-env1-node1-dataflow" {
		t.Errorf("unexpected pod name %q", pod.Name)
	}
	if len(pod.Spec.Containers) != 2 || pod.Spec.Containers[0].Name != "readout" || pod.Spec.Containers[1].Name != "stfb" {
		t.Errorf("unexpected containers: %+v", pod.Spec.Containers)
	}
	if pod.Spec.Containers[1].Image != "img-stfb" {
		t.Errorf("container content was not preserved: %+v", pod.Spec.Containers[1])
	}
	if len(pod.Spec.Volumes) != 2 {
		t.Errorf("expected volumes to be merged by name, got %+v", pod.Spec.Volumes)
	}
	if g := pod.Spec.TerminationGracePeriodSeconds; g == nil || *g != 30 {
		t.Errorf("expected max grace period 30, got %v", g)
	}
	if pod.Spec.RestartPolicy != v1.RestartPolicyNever {
		t.Errorf("expected restart policy Never, got %s", pod.Spec.RestartPolicy)
	}
	if pod.Labels[labelSharedPod] != "true" || pod.Labels["environment"] != "env1" {
		t.Errorf("unexpected labels: %v", pod.Labels)
	}
	if ref := metav1.GetControllerOf(pod); ref == nil || ref.Kind != "Environment" {
		t.Errorf("pod should be controlled by the Environment, got %v", ref)
	}

	for _, pt := range []pendingTask{a, b} {
		if pt.task.Spec.PodName != pod.Name || pt.task.Spec.ContainerName != pt.containerName {
			t.Errorf("task %s not pointing to its pod/container: %+v", pt.task.Name, pt.task.Spec)
		}
		if pt.task.Labels[labelTaskPod] != pod.Name {
			t.Errorf("task %s misses pod label: %v", pt.task.Name, pt.task.Labels)
		}
		if ref := metav1.GetControllerOf(pt.task); ref == nil || ref.Kind != "Environment" {
			t.Errorf("task %s should be controlled by the Environment, got %v", pt.task.Name, ref)
		}
	}
}

func TestBuildSharedPodErrors(t *testing.T) {
	key := podKey{node: "node1", resolvedNode: "node1", group: "g"}
	tests := []struct {
		name    string
		tasks   []pendingTask
		wantErr string
	}{
		{
			name: "conflicting pod-level settings",
			tasks: []pendingTask{
				testPendingTask("t1", "c1", nil),
				testPendingTask("t2", "c2", func(s *v1.PodSpec) { s.HostNetwork = false }),
			},
			wantErr: "conflicting pod-level settings (HostNetwork)",
		},
		{
			name: "same volume defined differently",
			tasks: []pendingTask{
				testPendingTask("t1", "c1", nil),
				testPendingTask("t2", "c2", func(s *v1.PodSpec) {
					s.Volumes = []v1.Volume{{Name: "tmp", VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: "/other"}}}}
				}),
			},
			wantErr: "volume",
		},
		{
			name: "duplicate container names",
			tasks: []pendingTask{
				testPendingTask("t1", "same", nil),
				testPendingTask("t2", "same", nil),
			},
			wantErr: "same container name",
		},
		{
			name: "task with multiple containers",
			tasks: []pendingTask{
				testPendingTask("t1", "c1", func(s *v1.PodSpec) { s.Containers = append(s.Containers, v1.Container{Name: "second"}) }),
			},
			wantErr: "exactly one container",
		},
		{
			name:    "invalid container name",
			tasks:   []pendingTask{testPendingTask("t1", "Not_Valid", nil)},
			wantErr: "invalid container name",
		},
		{
			name:    "no tasks",
			wantErr: "no tasks",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildSharedPod(newTestEnvironment(), key, tt.tasks, newTestScheme(t))
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestIsContainerFailed(t *testing.T) {
	failedStatus := v1.ContainerStatus{Name: "other", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 1}}}
	okDone := v1.ContainerStatus{Name: "mine", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 0}}}
	running := v1.ContainerStatus{Name: "mine", Ready: true, State: v1.ContainerState{Running: &v1.ContainerStateRunning{}}}

	tests := []struct {
		name       string
		pod        *v1.Pod
		wantFailed bool
	}{
		{
			name:       "neighbour failed, mine running",
			pod:        &v1.Pod{Status: v1.PodStatus{Phase: v1.PodRunning, ContainerStatuses: []v1.ContainerStatus{failedStatus, running}}},
			wantFailed: false,
		},
		{
			name:       "mine failed",
			pod:        &v1.Pod{Status: v1.PodStatus{Phase: v1.PodRunning, ContainerStatuses: []v1.ContainerStatus{{Name: "mine", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 2}}}}}},
			wantFailed: true,
		},
		{
			name:       "pod Failed because of neighbour, mine exited cleanly",
			pod:        &v1.Pod{Status: v1.PodStatus{Phase: v1.PodFailed, ContainerStatuses: []v1.ContainerStatus{failedStatus, okDone}}},
			wantFailed: false,
		},
		{
			name:       "pod Failed without status of mine (e.g. evicted)",
			pod:        &v1.Pod{Status: v1.PodStatus{Phase: v1.PodFailed, Reason: "Evicted"}},
			wantFailed: true,
		},
		{
			name: "init container failure affects everyone",
			pod: &v1.Pod{Status: v1.PodStatus{Phase: v1.PodPending,
				InitContainerStatuses: []v1.ContainerStatus{{Name: "init", State: v1.ContainerState{Terminated: &v1.ContainerStateTerminated{ExitCode: 1}}}}}},
			wantFailed: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, reason := isContainerFailed(tt.pod, "mine"); got != tt.wantFailed {
				t.Errorf("isContainerFailed() = %v (%s), want %v", got, reason, tt.wantFailed)
			}
		})
	}
}

func TestSharedPodTaskRequests(t *testing.T) {
	scheme := newTestScheme(t)
	_ = v1.AddToScheme(scheme)
	mk := func(name, node, pod string) *aliecsv1alpha1.Task {
		return &aliecsv1alpha1.Task{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "ns", Labels: map[string]string{labelTaskPod: pod}},
			Spec:       aliecsv1alpha1.TaskSpec{NodeName: node},
		}
	}
	c := fakeClient(t, scheme, mk("a", "node1", "p1"), mk("b", "node1", "p1"), mk("c", "node1", "p2"))
	r := &TaskReconciler{Client: c, NodeName: "node1"}

	pod := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "p1", Namespace: "ns", Labels: map[string]string{labelSharedPod: "true"}}, Spec: v1.PodSpec{NodeName: "node1"}}
	if reqs := r.sharedPodTaskRequests(t.Context(), pod); len(reqs) != 2 {
		t.Errorf("expected 2 requests, got %v", reqs)
	}

	pod.Spec.NodeName = "node2"
	if reqs := r.sharedPodTaskRequests(t.Context(), pod); len(reqs) != 0 {
		t.Errorf("expected no requests for pod on other node, got %v", reqs)
	}

	delete(pod.Labels, labelSharedPod)
	pod.Spec.NodeName = "node1"
	if reqs := r.sharedPodTaskRequests(t.Context(), pod); len(reqs) != 0 {
		t.Errorf("expected no requests for non-shared pod, got %v", reqs)
	}
}

func fakeClient(t *testing.T, scheme *runtime.Scheme, objs ...client.Object) client.Client {
	t.Helper()
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}
