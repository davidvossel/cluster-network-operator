package hypershift

import (
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"testing"
)

func TestParseHostedControlPlane(t *testing.T) {
	testCases := []struct {
		name                     string
		inputUnstructuredContent string
		expectedOutput           *HostedControlPlane
	}{
		{
			name: "Picks up expected IBMCloud HCP fields with apiserver networking set",
			expectedOutput: &HostedControlPlane{
				ClusterID:                    "31df7fa9-b1a7-4a66-98ef-c6920bf213d8",
				ControllerAvailabilityPolicy: HighlyAvailable,
				NodeSelector:                 nil,
				AdvertiseAddress:             "172.20.0.1",
				AdvertisePort:                2040,
			},
			inputUnstructuredContent: `
apiVersion: hypershift.openshift.io/v1beta1
kind: HostedControlPlane
spec:
  autoscaling: {}
  clusterID: 31df7fa9-b1a7-4a66-98ef-c6920bf213d8
  controllerAvailabilityPolicy: HighlyAvailable
  dns:
    baseDomain: tf71faa489656c98b18e2-a383e1dc466c308d41a756a1a66c2b6a-c000.us-south.satellite.test.appdomain.cloud
  infrastructureAvailabilityPolicy: HighlyAvailable
  issuerURL: https://kubernetes.default.svc
  networking:
    apiServer:
      advertiseAddress: 172.20.0.1
      port: 2040
`,
		},
		{
			name: "Picks up defaults appropriately",
			expectedOutput: &HostedControlPlane{
				ClusterID:                    "31df7fa9-b1a7-4a66-98ef-c6920bf213d8",
				ControllerAvailabilityPolicy: HighlyAvailable,
				NodeSelector:                 nil,
				AdvertiseAddress:             HostedClusterDefaultAdvertiseAddressIPV4,
				AdvertisePort:                int(HostedClusterDefaultAdvertisePort),
			},
			inputUnstructuredContent: `
apiVersion: hypershift.openshift.io/v1beta1
kind: HostedControlPlane
spec:
  autoscaling: {}
  clusterID: 31df7fa9-b1a7-4a66-98ef-c6920bf213d8
  controllerAvailabilityPolicy: HighlyAvailable
  dns:
    baseDomain: tf71faa489656c98b18e2-a383e1dc466c308d41a756a1a66c2b6a-c000.us-south.satellite.test.appdomain.cloud
  infrastructureAvailabilityPolicy: HighlyAvailable
  issuerURL: https://kubernetes.default.svc
`,
		},
		{
			name: "Picks up default ipv6 address",
			expectedOutput: &HostedControlPlane{
				ClusterID:                    "31df7fa9-b1a7-4a66-98ef-c6920bf213d8",
				ControllerAvailabilityPolicy: HighlyAvailable,
				NodeSelector:                 nil,
				AdvertiseAddress:             HostedClusterDefaultAdvertiseAddressIPV6,
				AdvertisePort:                2040,
			},
			inputUnstructuredContent: `
apiVersion: hypershift.openshift.io/v1beta1
kind: HostedControlPlane
spec:
  autoscaling: {}
  clusterID: 31df7fa9-b1a7-4a66-98ef-c6920bf213d8
  controllerAvailabilityPolicy: HighlyAvailable
  dns:
    baseDomain: tf71faa489656c98b18e2-a383e1dc466c308d41a756a1a66c2b6a-c000.us-south.satellite.test.appdomain.cloud
  infrastructureAvailabilityPolicy: HighlyAvailable
  issuerURL: https://kubernetes.default.svc
  networking:
    serviceNetwork:
    - cidr: "2001::/16"
    apiServer:
      port: 2040
`,
		},
		{
			name: "Picks up toleration",
			expectedOutput: &HostedControlPlane{
				ClusterID:                    "31df7fa9-b1a7-4a66-98ef-c6920bf213d8",
				ControllerAvailabilityPolicy: "HighlyAvailable",
				NodeSelector: map[string]string{
					"kubernetes.io/os": "linux",
				},
				Tolerations: []corev1.Toleration{
					{
						Key:               "node.kubernetes.io/not-ready",
						Operator:          "Exists",
						Value:             "",
						Effect:            "NoExecute",
						TolerationSeconds: nil,
					},
					{
						Key:               "node.kubernetes.io/unreachable",
						Operator:          "Exists",
						Value:             "",
						Effect:            "NoExecute",
						TolerationSeconds: nil,
					},
					{
						Key:               "node.kubernetes.io/memory-pressure",
						Operator:          "Exists",
						Value:             "",
						Effect:            "NoSchedule",
						TolerationSeconds: nil,
					},
					{
						Key:               "key1",
						Operator:          "Equal",
						Value:             "value1",
						Effect:            "NoSchedule",
						TolerationSeconds: nil,
					},
					{
						Key:               "key1",
						Operator:          "Exists",
						Value:             "",
						Effect:            "NoSchedule",
						TolerationSeconds: nil,
					},
				},
				AdvertiseAddress: "172.20.0.1",
				AdvertisePort:    6443,
				PriorityClass:    "",
			},
			inputUnstructuredContent: `
apiVersion: hypershift.openshift.io/v1beta1
kind: HostedControlPlane
spec:
  autoscaling: {}
  clusterID: 31df7fa9-b1a7-4a66-98ef-c6920bf213d8
  controllerAvailabilityPolicy: HighlyAvailable
  dns:
    baseDomain: tf71faa489656c98b18e2-a383e1dc466c308d41a756a1a66c2b6a-c000.us-south.satellite.test.appdomain.cloud
  infrastructureAvailabilityPolicy: HighlyAvailable
  nodeSelector:
    kubernetes.io/os: linux
  networking:
    clusterNetwork:
    - cidr: 10.132.0.0/14
    networkType: OVNKubernetes
    serviceNetwork:
    - cidr: 172.31.0.0/16
  olmCatalogPlacement: management
  platform:
    kubevirt:
      baseDomainPassthrough: true
      generateID: tgw7vsjfjm
    type: KubeVirt
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 300
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 300
  - effect: NoSchedule
    key: node.kubernetes.io/memory-pressure
    operator: Exists
  - key: "key1"
    operator: "Equal"
    value: "value1"
    effect: "NoSchedule"
  - key: "key1"
    operator: "Exists"
    effect: "NoSchedule"
status:
  conditions:
  - lastTransitionTime: "2024-07-17T05:55:02Z"
    message: Configuration passes validation
    observedGeneration: 1
    reason: AsExpected
    status: "True"
    type: ValidHostedControlPlaneConfiguration
`,
		},
	}
	g := NewGomegaWithT(t)
	for _, tc := range testCases {
		rawHostedControlPlane, err := yaml.ToJSON([]byte(tc.inputUnstructuredContent))
		g.Expect(err).NotTo(HaveOccurred())
		object, err := runtime.Decode(unstructured.UnstructuredJSONScheme, rawHostedControlPlane)
		g.Expect(err).NotTo(HaveOccurred())
		hcpUnstructured, ok := object.(*unstructured.Unstructured)
		g.Expect(ok).To(BeTrue())
		actualOutput, err := ParseHostedControlPlane(hcpUnstructured)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(actualOutput).To(Equal(tc.expectedOutput))
	}
}
