package v1alpha1

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Status Updater Tests", func() {
	It("works", func() {

		By("Create a dummy client mount")
		cm := &ClientMount{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: ClientMountSpec{
				Node:         "test",
				DesiredState: "mounted",
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

		updater := NewStatusUpdater[*ClientMount, *ClientMountStatus](cm)

		Expect(k8sClient.Create(context.TODO(), cm)).To(Succeed())

		Eventually(func() error {
			return k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(cm), cm)
		}).Should(Succeed())

		By("Update the Status section using the updater")
		cm.Status.Mounts = []ClientMountInfoStatus{
			{
				State: "mounted",
				Ready: false,
			},
		}
		cm.Status.Error = NewResourceError("Status Updater Test Error",
			fmt.Errorf("Error")).WithFatal()

		Expect(updater.Close(context.TODO(), k8sClient)).To(Succeed())

		Eventually(func(g Gomega) string {
			newcm := &ClientMount{}
			g.Expect(k8sClient.Get(context.TODO(), client.ObjectKeyFromObject(cm), newcm)).To(Succeed())
			return newcm.Status.Error.Error()
		}).Should(Equal(cm.Status.Error.Error()))

	})
})
