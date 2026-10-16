package helm_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/konflux-ci/caching/tests/testhelpers"
)

var _ = Describe("Helm Template Nginx Affinity Configuration", func() {
	Describe("Default Configuration", func() {
		It("should include pod anti-affinity rules by default", func() {
			output, err := renderNginxTemplate(testhelpers.HelmValues{
				Nginx: &testhelpers.NginxValues{
					Enabled:      true,
					ReplicaCount: 2,
				},
			})
			Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")

			nginxSection := extractNginxStatefulSetSection(output)
			Expect(nginxSection).To(ContainSubstring("podAntiAffinity"), "Should contain podAntiAffinity")
			Expect(nginxSection).To(ContainSubstring("preferredDuringSchedulingIgnoredDuringExecution"), "Should use preferred anti-affinity")
			Expect(nginxSection).To(ContainSubstring("kubernetes.io/hostname"), "Should use hostname topology key")
			Expect(nginxSection).To(ContainSubstring("weight: 100"), "Should have weight 100")
			Expect(nginxSection).To(ContainSubstring("app.kubernetes.io/name: "+testhelpers.NginxStatefulSetName), "Should target nginx pods")
			Expect(nginxSection).To(ContainSubstring("app.kubernetes.io/component: "+testhelpers.NginxComponentLabel), "Should target nginx component")
		})

		It("should not include required anti-affinity", func() {
			output, err := renderNginxTemplate(testhelpers.HelmValues{
				Nginx: &testhelpers.NginxValues{Enabled: true},
			})
			Expect(err).NotTo(HaveOccurred())

			nginxSection := extractNginxStatefulSetSection(output)
			Expect(nginxSection).To(ContainSubstring("preferredDuringSchedulingIgnoredDuringExecution"), "Should have preferred anti-affinity")
			Expect(nginxSection).NotTo(ContainSubstring("requiredDuringSchedulingIgnoredDuringExecution"), "Should not have required anti-affinity")
		})
	})

	Describe("Disabled Affinity", func() {
		It("should not include any affinity when user sets empty affinity", func() {
			output, err := renderNginxTemplate(testhelpers.HelmValues{
				Nginx: &testhelpers.NginxValues{
					Enabled:  true,
					Affinity: json.RawMessage("{}"),
				},
			})
			Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")

			nginxSection := extractNginxStatefulSetSection(output)
			Expect(nginxSection).NotTo(ContainSubstring("affinity:"), "Nginx statefulset should not contain affinity section")
			Expect(nginxSection).NotTo(ContainSubstring("podAntiAffinity"), "Nginx statefulset should not contain podAntiAffinity")
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
								"values": ["cache-nodes"]
							}]
						}]
					}
				}
			}`

			output, err := renderNginxTemplate(testhelpers.HelmValues{
				Nginx: &testhelpers.NginxValues{
					Enabled:  true,
					Affinity: json.RawMessage(customAffinity),
				},
			})
			Expect(err).NotTo(HaveOccurred(), "Helm template rendering should succeed")

			nginxSection := extractNginxStatefulSetSection(output)
			Expect(nginxSection).To(ContainSubstring("nodeAffinity"), "Should contain nodeAffinity from user")
			Expect(nginxSection).To(ContainSubstring("node-type"), "Should contain user's node selector")
			Expect(nginxSection).To(ContainSubstring("cache-nodes"), "Should contain user's node values")
			Expect(nginxSection).NotTo(ContainSubstring("podAntiAffinity"), "Should not have template's podAntiAffinity when user provides custom affinity")
		})

		It("should support custom pod anti-affinity with different settings", func() {
			customAffinity := `{
				"podAntiAffinity": {
					"requiredDuringSchedulingIgnoredDuringExecution": [{
						"labelSelector": {
							"matchLabels": {
								"app": "custom-nginx"
							}
						},
						"topologyKey": "topology.kubernetes.io/zone"
					}]
				}
			}`

			output, err := renderNginxTemplate(testhelpers.HelmValues{
				Nginx: &testhelpers.NginxValues{
					Enabled:  true,
					Affinity: json.RawMessage(customAffinity),
				},
			})
			Expect(err).NotTo(HaveOccurred())

			nginxSection := extractNginxStatefulSetSection(output)
			Expect(nginxSection).To(ContainSubstring("podAntiAffinity"), "Should have podAntiAffinity")
			Expect(nginxSection).To(ContainSubstring("requiredDuringSchedulingIgnoredDuringExecution"), "Should have required anti-affinity from user")
			Expect(nginxSection).To(ContainSubstring("topology.kubernetes.io/zone"), "Should use user's topology key")
			Expect(nginxSection).To(ContainSubstring("app: custom-nginx"), "Should use user's label selector")
			Expect(nginxSection).NotTo(ContainSubstring("preferredDuringSchedulingIgnoredDuringExecution"), "Should not have template's preferred rules when user provides custom")
		})
	})
})
