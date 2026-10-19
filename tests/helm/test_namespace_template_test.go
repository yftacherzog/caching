package helm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/konflux-ci/caching/tests/testhelpers"
)

var _ = Describe("Helm test namespace isolation", func() {
	It("places test infrastructure and proxy discovery in independent namespaces", func() {
		output, err := testhelpers.RenderHelmTemplate(chartPath, testhelpers.SquidHelmValues{
			Squid: &testhelpers.SquidValues{Namespace: "squid-proxy"},
			Nginx: &testhelpers.NginxValues{Enabled: true, Namespace: "nginx-proxy"},
			Test:  &testhelpers.TestValues{Namespace: "test-infra"},
		})
		Expect(err).NotTo(HaveOccurred())

		testPod := extractSection(output, "# Source: caching/templates/test-pod.yaml")
		Expect(testPod).To(ContainSubstring("namespace: test-infra"))
		Expect(testPod).To(ContainSubstring("name: SQUID_NAMESPACE\n      value: \"squid-proxy\""))
		Expect(testPod).To(ContainSubstring("name: NGINX_NAMESPACE\n      value: \"nginx-proxy\""))

		rbac := extractSection(output, "# Source: caching/templates/test-rbac.yaml")
		Expect(rbac).To(ContainSubstring("namespace: test-infra"))
		testServer := extractSection(output, "# Source: caching/templates/https-test-server.yaml")
		Expect(testServer).To(ContainSubstring("namespace: test-infra"))
	})

	It("renders Nginx tests with Squid disabled", func() {
		output, err := testhelpers.RenderHelmTemplate(chartPath, testhelpers.SquidHelmValues{
			Squid: &testhelpers.SquidValues{Enabled: testhelpers.BoolPtr(false)},
			Nginx: &testhelpers.NginxValues{Enabled: true, Namespace: "nginx-only"},
			Test:  &testhelpers.TestValues{Namespace: "test-infra"},
		})
		Expect(err).NotTo(HaveOccurred())

		Expect(extractSquidDeploymentSection(output)).To(BeEmpty())
		Expect(extractNginxStatefulSetSection(output)).To(ContainSubstring("namespace: nginx-only"))
		testPod := extractSection(output, "# Source: caching/templates/test-pod.yaml")
		Expect(testPod).To(ContainSubstring("name: SQUID_ENABLED\n      value: \"false\""))
		Expect(testPod).To(ContainSubstring("name: NGINX_ENABLED\n      value: \"true\""))
	})

	It("renders the test server Bundle as a cluster-scoped resource", func() {
		output, err := renderSquidTemplate(testhelpers.HelmValues{})
		Expect(err).NotTo(HaveOccurred())

		bundle := extractSection(output, "# Source: caching/templates/https-test-server-bundle.yaml")
		Expect(bundle).To(ContainSubstring("kind: Bundle"))
		Expect(bundle).NotTo(ContainSubstring("namespace:"))
	})

	It("rejects an invalid test namespace", func() {
		_, err := testhelpers.RenderHelmTemplate(chartPath, testhelpers.SquidHelmValues{
			Test: &testhelpers.TestValues{Namespace: "Invalid_Namespace"},
		})
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("/test/namespace"))
	})

	It("rejects a test namespace matching the enabled Squid namespace", func() {
		_, err := testhelpers.RenderHelmTemplate(chartPath, testhelpers.SquidHelmValues{
			Squid: &testhelpers.SquidValues{Namespace: "shared-proxy"},
			Test:  &testhelpers.TestValues{Namespace: "shared-proxy"},
		})
		Expect(err).To(MatchError(ContainSubstring("test.namespace must be distinct from squid.namespace")))
	})

	It("rejects a test namespace matching the enabled Nginx namespace", func() {
		_, err := testhelpers.RenderHelmTemplate(chartPath, testhelpers.SquidHelmValues{
			Nginx: &testhelpers.NginxValues{Enabled: true, Namespace: "shared-proxy"},
			Test:  &testhelpers.TestValues{Namespace: "shared-proxy"},
		})
		Expect(err).To(MatchError(ContainSubstring("test.namespace must be distinct from nginx.namespace")))
	})
})
