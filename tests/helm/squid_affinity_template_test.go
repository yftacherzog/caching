package helm_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"

	"github.com/konflux-ci/caching/tests/testhelpers"
)

var _ = Describe("Helm Template Squid Affinity Configuration", func() {
	Describe("Default Configuration", func() {
		It("should include pod anti-affinity rules by default", func() {
			output, err := renderSquidTemplate(testhelpers.HelmValues{
				ReplicaCount: 2,
			})
			Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")
			squidStatefulSet := extractSquidDeploymentSection(output)

			// Check that the squid statefulset is present
			Expect(squidStatefulSet).To(ContainSubstring("name: "+testhelpers.SquidStatefulSetName), "Should contain squid statefulset")
			Expect(squidStatefulSet).To(ContainSubstring("namespace: "+testhelpers.SquidNamespace), "Should be in the Squid namespace")

			// Check for anti-affinity configuration
			Expect(squidStatefulSet).To(ContainSubstring("podAntiAffinity"), "Should contain podAntiAffinity")
			Expect(squidStatefulSet).To(ContainSubstring("preferredDuringSchedulingIgnoredDuringExecution"), "Should use preferred anti-affinity")
			Expect(squidStatefulSet).To(ContainSubstring("kubernetes.io/hostname"), "Should use hostname topology key")
			Expect(squidStatefulSet).To(ContainSubstring("weight: 100"), "Should have weight 100")

			// Verify label selector
			Expect(squidStatefulSet).To(ContainSubstring("app.kubernetes.io/name: "+testhelpers.SquidStatefulSetName), "Should target squid pods")
			Expect(squidStatefulSet).To(ContainSubstring("app.kubernetes.io/component: "+testhelpers.SquidComponentLabel), "Should target squid component")
		})

		It("should not include required anti-affinity (only preferred)", func() {
			output, err := renderSquidTemplate(testhelpers.HelmValues{})
			Expect(err).NotTo(HaveOccurred())
			squidStatefulSet := extractSquidDeploymentSection(output)

			// Should have preferred but not required anti-affinity
			Expect(squidStatefulSet).To(ContainSubstring("preferredDuringSchedulingIgnoredDuringExecution"), "Should have preferred anti-affinity")
			Expect(squidStatefulSet).NotTo(ContainSubstring("requiredDuringSchedulingIgnoredDuringExecution"), "Should not have required anti-affinity")
		})
	})

	Describe("Disabled Affinity", func() {
		It("should not include any affinity when user sets empty affinity", func() {
			output, err := renderSquidTemplate(testhelpers.HelmValues{
				Affinity: json.RawMessage("{}"),
			})
			Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")

			// Should contain the statefulset but no affinity section
			// Extract just the squid statefulset section for more precise checking
			squidDeploymentSection := extractSquidDeploymentSection(output)
			Expect(squidDeploymentSection).To(ContainSubstring("name: "+testhelpers.SquidStatefulSetName), "Should contain squid statefulset")
			Expect(squidDeploymentSection).NotTo(ContainSubstring("affinity:"), "Squid statefulset should not contain affinity section")
			Expect(squidDeploymentSection).NotTo(ContainSubstring("podAntiAffinity"), "Squid statefulset should not contain podAntiAffinity")
		})
	})

	Describe("Custom Affinity", func() {
		It("should use custom node affinity instead of template defaults", func() {
			customAffinity := `{
				"nodeAffinity": {
					"requiredDuringSchedulingIgnoredDuringExecution": {
						"nodeSelectorTerms": [{
							"matchExpressions": [{
								"key": "node-type",
								"operator": "In",
								"values": ["proxy-nodes"]
							}]
						}]
					}
				}
			}`

			output, err := renderSquidTemplate(testhelpers.HelmValues{
				Affinity: json.RawMessage(customAffinity),
			})
			Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")

			// Should have user's node affinity
			squidDeploymentSection := extractSquidDeploymentSection(output)
			Expect(squidDeploymentSection).To(ContainSubstring("nodeAffinity"), "Should contain nodeAffinity from user")
			Expect(squidDeploymentSection).To(ContainSubstring("node-type"), "Should contain user's node selector")
			Expect(squidDeploymentSection).To(ContainSubstring("proxy-nodes"), "Should contain user's node values")

			// Should NOT have template's pod anti-affinity
			Expect(squidDeploymentSection).NotTo(ContainSubstring("podAntiAffinity"), "Should not have template's podAntiAffinity when user provides custom affinity")
		})

		It("should support custom pod anti-affinity with different settings", func() {
			customAffinity := `{
				"podAntiAffinity": {
					"requiredDuringSchedulingIgnoredDuringExecution": [{
						"labelSelector": {
							"matchLabels": {
								"app": "custom-squid"
							}
						},
						"topologyKey": "topology.kubernetes.io/zone"
					}]
				}
			}`

			output, err := renderSquidTemplate(testhelpers.HelmValues{
				Affinity: json.RawMessage(customAffinity),
			})
			Expect(err).NotTo(HaveOccurred())

			// Should have user's custom pod anti-affinity
			squidDeploymentSection := extractSquidDeploymentSection(output)
			Expect(squidDeploymentSection).To(ContainSubstring("podAntiAffinity"), "Should have podAntiAffinity")
			Expect(squidDeploymentSection).To(ContainSubstring("requiredDuringSchedulingIgnoredDuringExecution"), "Should have required anti-affinity from user")
			Expect(squidDeploymentSection).To(ContainSubstring("topology.kubernetes.io/zone"), "Should use user's topology key")
			Expect(squidDeploymentSection).To(ContainSubstring("app: custom-squid"), "Should use user's label selector")

			// Should NOT have template's preferred anti-affinity
			Expect(squidDeploymentSection).NotTo(ContainSubstring("preferredDuringSchedulingIgnoredDuringExecution"), "Should not have template's preferred rules when user provides custom")
		})
	})

	It("should include custom volumes and volumeMounts in squid container", func() {
		output, err := renderSquidTemplate(testhelpers.HelmValues{
			Volumes: []corev1.Volume{
				{
					Name: "custom-config",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "my-configmap",
							},
						},
					},
				},
			},
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      "custom-config",
					MountPath: "/etc/custom-config",
					ReadOnly:  true,
				},
			},
		})
		Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")

		squidDeploymentSection := extractSquidDeploymentSection(output)

		// Verify custom volume is present
		Expect(squidDeploymentSection).To(ContainSubstring("name: custom-config"), "Should contain custom volume name")
		Expect(squidDeploymentSection).To(ContainSubstring("name: my-configmap"), "Should contain configmap reference")

		// Verify custom volumeMount is present
		Expect(squidDeploymentSection).To(ContainSubstring("mountPath: /etc/custom-config"), "Should contain custom mount path")
	})
})
