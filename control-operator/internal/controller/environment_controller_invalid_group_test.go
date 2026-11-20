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
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	aliecsv1alpha1 "github.com/AliceO2Group/Control/control-operator/api/v1alpha1"
)

var _ = Describe("Environment Controller with an invalid Pod group", func() {
	const (
		namespace   = "default"
		envName     = "invalid-group-env"
		clusterNode = "invalid-node.cern.ch"
	)

	It("reports the conflict in the Environment status and creates no shared Pod", func() {
		ctx := context.Background()
		envKey := types.NamespacedName{Name: envName, Namespace: namespace}
		reconciler := &EnvironmentReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}

		Expect(k8sClient.Create(ctx, &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: clusterNode}})).To(Succeed())
		DeferCleanup(func() {
			env := &aliecsv1alpha1.Environment{}
			if err := k8sClient.Get(ctx, envKey, env); err == nil {
				controllerutil.RemoveFinalizer(env, environmentFinalizer)
				Expect(k8sClient.Update(ctx, env)).To(Succeed())
				Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, env))).To(Succeed())
			}
			Expect(client.IgnoreNotFound(k8sClient.Delete(ctx, &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: clusterNode}}))).To(Succeed())
		})

		def := func(name string, secCtx *v1.PodSecurityContext) aliecsv1alpha1.TaskDefinition {
			return aliecsv1alpha1.TaskDefinition{
				Name:     name,
				PodGroup: "g",
				Spec: aliecsv1alpha1.TaskTemplateSpec{Pod: v1.PodSpec{
					SecurityContext: secCtx,
					Containers:      []v1.Container{{Name: "main", Image: "img"}},
				}},
			}
		}
		Expect(k8sClient.Create(ctx, &aliecsv1alpha1.Environment{
			ObjectMeta:    metav1.ObjectMeta{Name: envName, Namespace: namespace},
			TaskTemplates: aliecsv1alpha1.TemplateSpecification{Tasks: map[string][]aliecsv1alpha1.TaskReference{}},
			Spec: aliecsv1alpha1.EnvironmentSpec{
				State: "standby",
				Tasks: map[string][]aliecsv1alpha1.TaskDefinition{
					"invalid-node": {def("a", nil), def("b", &v1.PodSecurityContext{FSGroup: new(int64)})},
				},
			},
		})).To(Succeed())

		for range 2 { // finalizer, deployment
			_, err := reconciler.Reconcile(ctx, reconcile.Request{NamespacedName: envKey})
			Expect(err).NotTo(HaveOccurred())
		}

		env := &aliecsv1alpha1.Environment{}
		Expect(k8sClient.Get(ctx, envKey, env)).To(Succeed())
		cond := meta.FindStatusCondition(env.Status.Conditions, "PodsCreated")
		Expect(cond).NotTo(BeNil(), "the status update must not be rejected by the CRD validation")
		Expect(cond.Status).To(Equal(metav1.ConditionFalse))
		Expect(cond.Message).To(ContainSubstring("SecurityContext"))

		pods := &v1.PodList{}
		Expect(k8sClient.List(ctx, pods, client.InNamespace(namespace), client.MatchingLabels{"environment": envName})).To(Succeed())
		Expect(pods.Items).To(BeEmpty())
	})
})
