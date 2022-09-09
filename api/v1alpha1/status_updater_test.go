package v1alpha1

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Status Updater Test", func() {
	var cm = &ClientMount{}
	BeforeEach(func() {
		cm = &ClientMount{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: ClientMountSpec{
				Node:         "test",
				DesiredState: "unmounted",
				Mounts: []ClientMountInfo{
					{
						MountPath: "/",
						Options:   "",
						Device: ClientMountDevice{
							Type: "reference",
						},
						Type:       "none",
						TargetType: "file",
						Compute:    "",
					},
				},
			},
		}

		Expect(k8sClient.Create(context.TODO(), cm)).To(Succeed())

		Eventually(func() error {
			return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(cm), cm)
		}).Should(Succeed())
	})

	AfterEach(func() {
		Expect(k8sClient.Delete(context.TODO(), cm)).To(Succeed())
	})

	It("Updates Status Only", func() {

		By("Create new status updater")
		updater := NewStatusUpdater[*ClientMountStatus](cm)

		By("Toggle some spec and status fields")
		cm.Spec.DesiredState = "mounted"
		cm.Status.Mounts = []ClientMountInfoStatus{
			{
				State: "mounted",
				Ready: true,
			},
		}

		By("Updating the status, not the spec, using the updater")
		Expect(updater.CloseWithStatusUpdate(context.TODO(), k8sClient)).To(Succeed())

		By("Expect the status to change but the spec to remain the same")
		Eventually(func(g Gomega) bool {
			newcm := &ClientMount{}
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(cm), newcm)).To(Succeed())
			g.Expect(newcm.Status.Mounts).To(HaveLen(1))
			return newcm.Status.Mounts[0].Ready && newcm.Spec.DesiredState == "unmounted"
		}).Should(BeTrue())
	})
})
