/*
 * Copyright 2021-2025 Hewlett Packard Enterprise Development LP
 * Other additional copyright holders may be indicated within.
 *
 * The entirety of this work is licensed under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 *
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha4

import (
	"github.com/DataWorkflowServices/dws/utils/updater"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClientMountDeviceLustre defines the lustre device information for mounting
type ClientMountDeviceLustre struct {
	// Lustre fsname
	FileSystemName string `json:"fileSystemName"`

	// List of mgsAddresses of the form [address]@[lnet]
	MgsAddresses string `json:"mgsAddresses"`
}

// ClientMountNVMeDesc uniquely describes an NVMe namespace
type ClientMountNVMeDesc struct {
	// Serial number of the base NVMe device
	DeviceSerial string `json:"deviceSerial"`

	// Id of the Namespace on the NVMe device (e.g., "2")
	NamespaceID string `json:"namespaceID"`

	// Globally unique namespace ID
	NamespaceGUID string `json:"namespaceGUID"`
}

// ClientMountLVMDeviceType specifies the go type for LVMDeviceType
type ClientMountLVMDeviceType string

const (
	// ClientMountLVMDeviceTypeNVMe specifies the NVMe constant device type
	ClientMountLVMDeviceTypeNVMe ClientMountLVMDeviceType = "nvme"
)

// ClientMountDeviceLVM defines an LVM device by the VG/LV pair and optionally
// the drives that are the PVs.
type ClientMountDeviceLVM struct {
	// Type of underlying block deices used for the PVs
	// +kubebuilder:validation:Enum=nvme
	DeviceType ClientMountLVMDeviceType `json:"deviceType"`

	// List of NVMe namespaces that are used by the VG
	NVMeInfo []ClientMountNVMeDesc `json:"nvmeInfo,omitempty"`

	// LVM volume group name
	VolumeGroup string `json:"volumeGroup,omitempty"`

	// LVM logical volume name
	LogicalVolume string `json:"logicalVolume,omitempty"`
}

// ClientMountDeviceReference is an reference to a different Kubernetes object
// where device information can be found
type ClientMountDeviceReference struct {
	// Object reference for the device information
	ObjectReference corev1.ObjectReference `json:"objectReference"`

	// Optional private data for the driver
	Data int `json:"data,omitempty"`
}

// ClientMountDeviceType specifies the go type for device type
type ClientMountDeviceType string

const (
	// ClientMountDeviceTypeLustre is used to define the device as a Lustre file system
	ClientMountDeviceTypeLustre ClientMountDeviceType = "lustre"

	// ClientMountDeviceTypeLVM is used to define the device as a LVM logical volume
	ClientMountDeviceTypeLVM ClientMountDeviceType = "lvm"

	// ClientMountDeviceTypeReference is used when the device information is described in
	// a separate Kubernetes resource. The clientmountd (or another controller doing the mounts)
	// must know how to interpret the resource to extract the device information.
	ClientMountDeviceTypeReference ClientMountDeviceType = "reference"
)

// ClientMountDevice defines the device to mount
type ClientMountDevice struct {
	// +kubebuilder:validation:Enum=lustre;lvm;reference
	Type ClientMountDeviceType `json:"type"`

	// Lustre specific device information
	Lustre *ClientMountDeviceLustre `json:"lustre,omitempty"`

	// LVM logical volume specific device information
	LVM *ClientMountDeviceLVM `json:"lvm,omitempty"`

	DeviceReference *ClientMountDeviceReference `json:"deviceReference,omitempty"`
}

// ClientMountInfo defines a single mount
type ClientMountInfo struct {
	// Client path for mount target
	MountPath string `json:"mountPath"`

	// UserID to set for the mount
	UserID uint32 `json:"userID,omitempty"`

	// GroupID to set for the mount
	GroupID uint32 `json:"groupID,omitempty"`

	// SetPermissions will set UserID and GroupID on the mount if true
	SetPermissions bool `json:"setPermissions"`

	// Options for the file system mount
	Options string `json:"options"`

	// Description of the device to mount
	Device ClientMountDevice `json:"device"`

	// mount type
	// +kubebuilder:validation:Enum=lustre;xfs;gfs2;none
	Type string `json:"type"`

	// TargetType determines whether the mount target is a file or a directory
	// +kubebuilder:validation:Enum=file;directory
	TargetType string `json:"targetType"`

	// Compute is the name of the compute node which shares this mount if present. Empty if not shared.
	Compute string `json:"compute,omitempty"`
}

// ClientMountState specifies the go type for MountState
type ClientMountState string

// ClientMountState string constants
const (
	ClientMountStateMounted   ClientMountState = "mounted"
	ClientMountStateUnmounted ClientMountState = "unmounted"
)

// ClientMountSpec defines the desired state of ClientMount
type ClientMountSpec struct {
	// Name of the client node that is targeted by this mount
	Node string `json:"node"`

	// Desired state of the mount point
	// +kubebuilder:validation:Enum=mounted;unmounted
	DesiredState ClientMountState `json:"desiredState"`

	// List of mounts to create on this client
	// +kubebuilder:validation:MinItems=1
	Mounts []ClientMountInfo `json:"mounts"`
}

// ClientMountInfoStatus is the status for a single mount point
type ClientMountInfoStatus struct {
	// Current state
	// +kubebuilder:validation:Enum=mounted;unmounted
	State ClientMountState `json:"state"`

	// Ready indicates whether status.state has been achieved
	Ready bool `json:"ready"`
}

// ClientMountStatus defines the observed state of ClientMount
type ClientMountStatus struct {
	// List of mount statuses
	Mounts []ClientMountInfoStatus `json:"mounts"`

	// Rollup of each mounts ready status
	AllReady bool `json:"allReady"`

	// Error information
	ResourceError `json:",inline"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="DESIREDSTATE",type="string",JSONPath=".spec.desiredState",description="The desired state"
//+kubebuilder:printcolumn:name="READY",type="boolean",JSONPath=".status.allReady",description="True if desired state is achieved"
//+kubebuilder:printcolumn:name="ERROR",type="string",JSONPath=".status.error.severity"
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// ClientMount is the Schema for the clientmounts API
type ClientMount struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClientMountSpec   `json:"spec,omitempty"`
	Status ClientMountStatus `json:"status,omitempty"`
}

func (c *ClientMount) GetStatus() updater.Status[*ClientMountStatus] {
	return &c.Status
}

//+kubebuilder:object:root=true

// ClientMountList contains a list of ClientMount
type ClientMountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClientMount `json:"items"`
}

// GetObjectList returns a list of Client references.
func (c *ClientMountList) GetObjectList() []client.Object {
	objectList := []client.Object{}

	for i := range c.Items {
		objectList = append(objectList, &c.Items[i])
	}

	return objectList
}

func init() {
	SchemeBuilder.Register(&ClientMount{}, &ClientMountList{})
}
