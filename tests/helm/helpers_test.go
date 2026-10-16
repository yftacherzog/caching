package helm_test

import (
	"strings"

	"github.com/konflux-ci/caching/tests/testhelpers"
)

// chartPath is the path to the Helm chart under test
const chartPath = "./caching"

// renderSquidTemplate renders squid in isolation so nginx resources cannot
// satisfy assertions intended for squid resources.
func renderSquidTemplate(values testhelpers.HelmValues) (string, error) {
	values.Nginx = &testhelpers.NginxValues{Enabled: false}
	return testhelpers.RenderHelmTemplate(chartPath, values)
}

// renderSquidTemplateWithKubeVersion renders squid in isolation for a
// specific Kubernetes version.
func renderSquidTemplateWithKubeVersion(values testhelpers.HelmValues, kubeVersion string) (string, error) {
	values.Nginx = &testhelpers.NginxValues{Enabled: false}
	return testhelpers.RenderHelmTemplateWithKubeVersion(chartPath, values, kubeVersion)
}

// renderNginxTemplate renders nginx in isolation so squid resources cannot
// satisfy assertions intended for nginx resources.
func renderNginxTemplate(values testhelpers.HelmValues) (string, error) {
	values.Squid = &testhelpers.SquidValues{Enabled: testhelpers.BoolPtr(false)}
	return testhelpers.RenderHelmTemplate(chartPath, values)
}

// renderNginxTemplateWithKubeVersion renders nginx in isolation for a
// specific Kubernetes version.
func renderNginxTemplateWithKubeVersion(values testhelpers.HelmValues, kubeVersion string) (string, error) {
	values.Squid = &testhelpers.SquidValues{Enabled: testhelpers.BoolPtr(false)}
	return testhelpers.RenderHelmTemplateWithKubeVersion(chartPath, values, kubeVersion)
}

// extractSection is a generic helper function that extracts a specific template section from helm output
// based on the Helm source comment marker (e.g., "# Source: caching/templates/deployment.yaml")
func extractSection(helmOutput, sourceMarker string) string {
	lines := strings.Split(helmOutput, "\n")
	var sectionLines []string
	inSection := false

	for _, line := range lines {
		// Start capturing when we find the source marker
		if strings.Contains(line, sourceMarker) {
			inSection = true
			continue
		}

		// Stop capturing when we hit the next resource
		if inSection && strings.HasPrefix(line, "---") {
			break
		}

		if inSection {
			sectionLines = append(sectionLines, line)
		}
	}

	return strings.Join(sectionLines, "\n")
}

// extractSquidDeploymentSection extracts just the squid statefulset YAML for precise testing
func extractSquidDeploymentSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/deployment.yaml")
}

// extractNginxStatefulSetSection extracts just the nginx statefulset YAML for precise testing
func extractNginxStatefulSetSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-statefulset.yaml")
}

// extractNginxConfigMapSection extracts just the nginx configmap YAML for precise testing
func extractNginxConfigMapSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-configmap.yaml")
}

// extractNginxServiceSection extracts just the nginx service YAML (non-headless) for precise testing
func extractNginxServiceSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-service.yaml")
}

// extractNginxHeadlessServiceSection extracts just the nginx headless service YAML for precise testing
func extractNginxHeadlessServiceSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-headless-service.yaml")
}

// extractSquidServiceSection extracts just the squid service YAML (non-headless) for precise testing
func extractSquidServiceSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/service.yaml")
}

// extractSquidHeadlessServiceSection extracts just the squid headless service YAML for precise testing
func extractSquidHeadlessServiceSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/headless-service.yaml")
}

// extractSquidConfigMapSection extracts just the squid configmap YAML for precise testing
func extractSquidConfigMapSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/configmap.yaml")
}

// extractNginxExporterConfigMapSection extracts just the nginx exporter configmap YAML for precise testing
func extractNginxExporterConfigMapSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-exporter-configmap.yaml")
}

// extractNginxServiceMonitorSection extracts just the nginx servicemonitor YAML for precise testing
func extractNginxServiceMonitorSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-servicemonitor.yaml")
}

// extractNginxNamespaceSection extracts just the nginx namespace YAML for precise testing
func extractNginxNamespaceSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/nginx-namespace.yaml")
}

// extractSquidServiceMonitorSection extracts just the squid servicemonitor YAML for precise testing
func extractSquidServiceMonitorSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/servicemonitor.yaml")
}

// extractTrustManagerBundleSection extracts just the trust-manager bundle YAML for precise testing
func extractTrustManagerBundleSection(helmOutput string) string {
	return extractSection(helmOutput, "# Source: caching/templates/trust-manager-bundle.yaml")
}
