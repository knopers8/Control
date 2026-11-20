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
	"fmt"
	"reflect"
	"sort"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	aliecsv1alpha1 "github.com/AliceO2Group/Control/control-operator/api/v1alpha1"
)

const (
	// labelSharedPod marks Pods that host several Tasks and are owned by an Environment.
	labelSharedPod = "aliecs.alice.cern/shared-pod"
	// labelTaskPod is set on Tasks living in a shared Pod and holds the Pod name.
	labelTaskPod  = "aliecs.alice.cern/pod"
	labelPodGroup = "podGroup"
)

// podKey identifies one shared Pod: Tasks of an Environment on the same node with the same group.
type podKey struct {
	node         string // node name as given in the Environment
	resolvedNode string // actual cluster node name
	group        string
}

// pendingTask is a Task which is ready to be created, but belongs to a shared Pod.
type pendingTask struct {
	task          *aliecsv1alpha1.Task
	containerName string
}

type pendingPods map[podKey][]pendingTask

func (p pendingPods) add(key podKey, task *aliecsv1alpha1.Task, containerName string) {
	p[key] = append(p[key], pendingTask{task: task, containerName: containerName})
}

// sortedKeys returns the keys in a deterministic order.
func (p pendingPods) sortedKeys() []podKey {
	keys := make([]podKey, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].node != keys[j].node {
			return keys[i].node < keys[j].node
		}
		return keys[i].group < keys[j].group
	})
	return keys
}

func sharedPodName(env *aliecsv1alpha1.Environment, key podKey) string {
	return strings.ToLower(fmt.Sprintf("aliecs-pod-%s-%s-%s", env.Name, key.node, key.group))
}

// podLevelSpec returns a copy of the spec without the parts which are allowed to differ among Tasks
// sharing a Pod (everything else has to be identical, as the Pod is one scheduling unit).
func podLevelSpec(spec v1.PodSpec) v1.PodSpec {
	s := *spec.DeepCopy()
	s.Containers = nil
	s.InitContainers = nil
	s.Volumes = nil
	s.TerminationGracePeriodSeconds = nil
	return s
}

// differingPodFields returns the names of the top-level PodSpec fields which differ between a and b.
func differingPodFields(a, b v1.PodSpec) []string {
	va, vb := reflect.ValueOf(a), reflect.ValueOf(b)
	var diff []string
	for i := 0; i < va.NumField(); i++ {
		if !equality.Semantic.DeepEqual(va.Field(i).Interface(), vb.Field(i).Interface()) {
			diff = append(diff, va.Type().Field(i).Name)
		}
	}
	return diff
}

// buildSharedPod merges the Pod specs of Tasks into one Pod, with one container per Task. The Tasks are
// updated to point to the Pod and container. The Pod and the Tasks are owned by the Environment.
func buildSharedPod(env *aliecsv1alpha1.Environment, key podKey, tasks []pendingTask, scheme *runtime.Scheme) (*v1.Pod, error) {
	if len(tasks) == 0 {
		return nil, fmt.Errorf("no tasks for pod group %q on node %s", key.group, key.node)
	}

	podName := sharedPodName(env, key)
	if errs := validation.IsDNS1123Subdomain(podName); len(errs) > 0 {
		return nil, fmt.Errorf("invalid pod name %q: %s", podName, strings.Join(errs, "; "))
	}

	first := tasks[0]
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: env.Namespace,
			Labels: map[string]string{
				"environment":  env.Name,
				"node":         key.node,
				labelPodGroup:  key.group,
				labelSharedPod: "true",
			},
		},
		Spec: podLevelSpec(first.task.Spec.Pod),
	}
	baseLevel := podLevelSpec(first.task.Spec.Pod)

	containerNames := map[string]string{}
	volumes := map[string]v1.Volume{}
	initContainers := map[string]v1.Container{}
	var grace *int64

	for _, pt := range tasks {
		spec := pt.task.Spec.Pod
		if len(spec.Containers) != 1 {
			return nil, fmt.Errorf("task %s must define exactly one container to be placed in a shared pod, has %d",
				pt.task.Name, len(spec.Containers))
		}
		if errs := validation.IsDNS1123Label(pt.containerName); len(errs) > 0 {
			return nil, fmt.Errorf("invalid container name %q for task %s: %s", pt.containerName, pt.task.Name, strings.Join(errs, "; "))
		}
		if other, dup := containerNames[pt.containerName]; dup {
			return nil, fmt.Errorf("tasks %s and %s have the same container name %q in pod %s", other, pt.task.Name, pt.containerName, podName)
		}
		containerNames[pt.containerName] = pt.task.Name

		if diff := differingPodFields(baseLevel, podLevelSpec(spec)); len(diff) > 0 {
			return nil, fmt.Errorf("tasks %s and %s have conflicting pod-level settings (%s) and cannot share pod %s",
				first.task.Name, pt.task.Name, strings.Join(diff, ", "), podName)
		}

		container := *spec.Containers[0].DeepCopy()
		container.Name = pt.containerName
		pod.Spec.Containers = append(pod.Spec.Containers, container)

		for _, vol := range spec.Volumes {
			if existing, ok := volumes[vol.Name]; ok {
				if !equality.Semantic.DeepEqual(existing, vol) {
					return nil, fmt.Errorf("volume %q is defined differently by tasks sharing pod %s", vol.Name, podName)
				}
				continue
			}
			volumes[vol.Name] = vol
			pod.Spec.Volumes = append(pod.Spec.Volumes, *vol.DeepCopy())
		}
		for _, ic := range spec.InitContainers {
			if existing, ok := initContainers[ic.Name]; ok {
				if !equality.Semantic.DeepEqual(existing, ic) {
					return nil, fmt.Errorf("init container %q is defined differently by tasks sharing pod %s", ic.Name, podName)
				}
				continue
			}
			initContainers[ic.Name] = ic
			pod.Spec.InitContainers = append(pod.Spec.InitContainers, *ic.DeepCopy())
		}
		if g := spec.TerminationGracePeriodSeconds; g != nil && (grace == nil || *g > *grace) {
			v := *g
			grace = &v
		}
	}

	pod.Spec.TerminationGracePeriodSeconds = grace
	pod.Spec.RestartPolicy = v1.RestartPolicyNever

	if err := controllerutil.SetControllerReference(env, pod, scheme); err != nil {
		return nil, fmt.Errorf("setting owner of pod %s: %w", podName, err)
	}

	for _, pt := range tasks {
		pt.task.Spec.PodName = podName
		pt.task.Spec.ContainerName = pt.containerName
		if pt.task.Labels == nil {
			pt.task.Labels = map[string]string{}
		}
		pt.task.Labels[labelTaskPod] = podName
		pt.task.Labels[labelPodGroup] = key.group
		if err := controllerutil.SetControllerReference(env, pt.task, scheme); err != nil {
			return nil, fmt.Errorf("setting owner of task %s: %w", pt.task.Name, err)
		}
	}

	return pod, nil
}
