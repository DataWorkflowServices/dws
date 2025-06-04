/*
 * Copyright 2023-2025 Hewlett Packard Enterprise Development LP
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

package v1alpha1

import (
	"os"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	dwsv1alpha5 "github.com/DataWorkflowServices/dws/api/v1alpha5"
	utilconversion "github.com/DataWorkflowServices/dws/github/cluster-api/util/conversion"
)

var convertlog = logf.Log.V(2).WithName("convert-v1alpha1")

func (src *ClientMount) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert ClientMount To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.ClientMount)

	if err := Convert_v1alpha1_ClientMount_To_v1alpha5_ClientMount(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.ClientMount{}
	hasAnno, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	// v1alpha2 removed Error.Recoverable and uses Error.Severity and Error.Type, instead.
	if hasAnno && restored.Status.Error != nil {
		dst.Status.Error.Type = restored.Status.Error.Type
		dst.Status.Error.Severity = restored.Status.Error.Severity
	}
	if src.Status.Error != nil && !src.Status.Error.Recoverable {
		dst.Status.Error.Severity = dwsv1alpha5.SeverityFatal
	}

	// v1alpha2 added a rollup of all the mounts' ready flags
	if hasAnno {
		dst.Status.AllReady = restored.Status.AllReady
	} else {
		dst.Status.AllReady = true
		for _, mount := range src.Status.Mounts {
			if !mount.Ready {
				dst.Status.AllReady = false
				break
			}
		}
	}

	return nil
}

func (dst *ClientMount) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.ClientMount)
	convertlog.Info("Convert ClientMount From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_ClientMount_To_v1alpha1_ClientMount(src, dst, nil); err != nil {
		return err
	}

	// v1alpha2 removed Error.Recoverable, and it must be translated from
	// other fields.
	if src.Status.Error != nil {
		if src.Status.Error.Severity == dwsv1alpha5.SeverityFatal {
			dst.Status.Error.Recoverable = false
		} else {
			dst.Status.Error.Recoverable = true
		}
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)

}

func (src *Computes) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert Computes To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.Computes)

	if err := Convert_v1alpha1_Computes_To_v1alpha5_Computes(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.Computes{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	return nil
}

func (dst *Computes) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.Computes)
	convertlog.Info("Convert Computes From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_Computes_To_v1alpha1_Computes(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

func (src *DWDirectiveRule) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert DWDirectiveRule To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.DWDirectiveRule)

	if err := Convert_v1alpha1_DWDirectiveRule_To_v1alpha5_DWDirectiveRule(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.DWDirectiveRule{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	return nil
}

func (dst *DWDirectiveRule) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.DWDirectiveRule)
	convertlog.Info("Convert DWDirectiveRule From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_DWDirectiveRule_To_v1alpha1_DWDirectiveRule(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

func (src *DirectiveBreakdown) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert DirectiveBreakdown To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.DirectiveBreakdown)

	if err := Convert_v1alpha1_DirectiveBreakdown_To_v1alpha5_DirectiveBreakdown(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.DirectiveBreakdown{}
	hasAnno, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	// v1alpha2 removed Error.Recoverable and uses Error.Severity and Error.Type, instead.
	if hasAnno && restored.Status.Error != nil {
		dst.Status.Error.Type = restored.Status.Error.Type
		dst.Status.Error.Severity = restored.Status.Error.Severity
	}
	if src.Status.Error != nil && !src.Status.Error.Recoverable {
		dst.Status.Error.Severity = dwsv1alpha5.SeverityFatal
	}

	if hasAnno {
		dst.Status.Requires = restored.Status.Requires
	}

	return nil
}

func (dst *DirectiveBreakdown) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.DirectiveBreakdown)
	convertlog.Info("Convert DirectiveBreakdown From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_DirectiveBreakdown_To_v1alpha1_DirectiveBreakdown(src, dst, nil); err != nil {
		return err
	}

	// v1alpha2 removed Error.Recoverable, and it must be translated from
	// other fields.
	if src.Status.Error != nil {
		if src.Status.Error.Severity == dwsv1alpha5.SeverityFatal {
			dst.Status.Error.Recoverable = false
		} else {
			dst.Status.Error.Recoverable = true
		}
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

func (src *PersistentStorageInstance) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert PersistentStorageInstance To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.PersistentStorageInstance)

	if err := Convert_v1alpha1_PersistentStorageInstance_To_v1alpha5_PersistentStorageInstance(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.PersistentStorageInstance{}
	hasAnno, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	// v1alpha2 removed Error.Recoverable and uses Error.Severity and Error.Type, instead.
	if hasAnno && restored.Status.Error != nil {
		dst.Status.Error.Type = restored.Status.Error.Type
		dst.Status.Error.Severity = restored.Status.Error.Severity
	}
	if src.Status.Error != nil && !src.Status.Error.Recoverable {
		dst.Status.Error.Severity = dwsv1alpha5.SeverityFatal
	}

	return nil
}

func (dst *PersistentStorageInstance) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.PersistentStorageInstance)
	convertlog.Info("Convert PersistentStorageInstance From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_PersistentStorageInstance_To_v1alpha1_PersistentStorageInstance(src, dst, nil); err != nil {
		return err
	}

	// v1alpha2 removed Error.Recoverable, and it must be translated from
	// other fields.
	if src.Status.Error != nil {
		if src.Status.Error.Severity == dwsv1alpha5.SeverityFatal {
			dst.Status.Error.Recoverable = false
		} else {
			dst.Status.Error.Recoverable = true
		}
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

func (src *Servers) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert Servers To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.Servers)

	if err := Convert_v1alpha1_Servers_To_v1alpha5_Servers(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.Servers{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	// v1alpha2 introduced Status.ResourceError.
	if restored.Status.Error != nil {
		// Allocate a resource here, because v1alpha1 didn't have this.
		dst.Status.Error = dwsv1alpha5.NewResourceError("")
		dst.Status.Error.DebugMessage = restored.Status.Error.DebugMessage
		dst.Status.Error.UserMessage = restored.Status.Error.UserMessage
		dst.Status.Error.Type = restored.Status.Error.Type
		dst.Status.Error.Severity = restored.Status.Error.Severity
	}

	return nil
}

func (dst *Servers) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.Servers)
	convertlog.Info("Convert Servers From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_Servers_To_v1alpha1_Servers(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

func (src *Storage) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert Storage To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.Storage)

	if err := Convert_v1alpha1_Storage_To_v1alpha5_Storage(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.Storage{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	dst.Spec.Mode = restored.Spec.Mode

	return nil
}

func (dst *Storage) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.Storage)
	convertlog.Info("Convert Storage From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_Storage_To_v1alpha1_Storage(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

func (src *SystemConfiguration) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert SystemConfiguration To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.SystemConfiguration)

	if err := Convert_v1alpha1_SystemConfiguration_To_v1alpha5_SystemConfiguration(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.SystemConfiguration{}
	hasAnno, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	if hasAnno {
		dst.Spec.PortsCooldownInSeconds = restored.Spec.PortsCooldownInSeconds
		dst.Status.Error = restored.Status.Error

		// dst.Spec.ComputeNodes: The destination does not have this.
		// Instead, it finds the computes that are already in the
		// dst.Spec.StorageNodes list.

		dst.Spec.ExternalComputeNodes = restored.Spec.ExternalComputeNodes
	} else {
		// The v1alpha1 resource's spec.ComputeNodes list is a
		// combination of all compute nodes from the spec.StorageNodes
		// list as well as any external computes.
		// The v1alpha1.FindExternalComputes() method walks through
		// the spec.Computes list to determine which ones are external.
		externComputes := src.FindExternalComputes()
		dstExternComputes := make([]dwsv1alpha5.SystemConfigurationExternalComputeNode, len(externComputes))
		idx := 0
		for _, name := range externComputes {
			dstExternComputes[idx].Name = name
			idx++
		}
		dst.Spec.ExternalComputeNodes = dstExternComputes
	}
	return nil
}

func (dst *SystemConfiguration) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.SystemConfiguration)
	convertlog.Info("Convert SystemConfiguration From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_SystemConfiguration_To_v1alpha1_SystemConfiguration(src, dst, nil); err != nil {
		return err
	}

	// The v1alpha1 resource's spec.ComputeNodes list is a combination
	// of all compute nodes from the spec.StorageNodes list as well as any
	// external computes.
	// The v1alpha2 src.Computes() method returns only the compute nodes
	// from the spec.StorageNodes list. To retrieve the external computes we must
	// also use the ComputesExternal() method.
	computes := make([]SystemConfigurationComputeNode, 0)
	for _, name := range src.Computes() {
		computes = append(computes, SystemConfigurationComputeNode{Name: *name})
	}
	for _, name := range src.ComputesExternal() {
		computes = append(computes, SystemConfigurationComputeNode{Name: *name})
	}
	dst.Spec.ComputeNodes = computes

	// +crdbumper:carryforward:begin="SystemConfiguration.ConvertFrom"
	// In a non-test environment, ENVIRONMENT will be set. Don't save Hub data in the
	// annotations in this case. The SystemConfiguration resource can be very large, and
	// the annotation will be too large to store. In a test environment, we want the Hub
	// data saved in the annotations to test hub-spoke-hub and spoke-hub-spoke conversions.
	if _, found := os.LookupEnv("ENVIRONMENT"); found {
		return nil
	}
	// +crdbumper:carryforward:end

	// Preserve Hub data on down-conversion except for metadata
	return utilconversion.MarshalData(src, dst)
}

func (src *Workflow) ConvertTo(dstRaw conversion.Hub) error {
	convertlog.Info("Convert Workflow To Hub", "name", src.GetName(), "namespace", src.GetNamespace())
	dst := dstRaw.(*dwsv1alpha5.Workflow)

	if err := Convert_v1alpha1_Workflow_To_v1alpha5_Workflow(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &dwsv1alpha5.Workflow{}
	hasAnno, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	// EDIT THIS FUNCTION! If the annotation is holding anything that is
	// hub-specific then copy it into 'dst' from 'restored'.
	// Otherwise, you may comment out UnmarshalData() until it's needed.

	if hasAnno {
		dst.Spec.JobID = restored.Spec.JobID
		dst.Status.Requires = restored.Status.Requires
		dst.Status.WorkflowToken = restored.Status.WorkflowToken
	} else {
		dst.Spec.JobID = intstr.FromInt(src.Spec.JobID)
	}

	return nil
}

func (dst *Workflow) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*dwsv1alpha5.Workflow)
	convertlog.Info("Convert Workflow From Hub", "name", src.GetName(), "namespace", src.GetNamespace())

	if err := Convert_v1alpha5_Workflow_To_v1alpha1_Workflow(src, dst, nil); err != nil {
		return err
	}

	dst.Spec.JobID = src.Spec.JobID.IntValue()

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

// The List-based ConvertTo/ConvertFrom routines are never used by the
// conversion webhook, but the conversion-verifier tool wants to see them.
// The conversion-gen tool generated the Convert_X_to_Y routines, should they
// ever be needed.

func resource(resource string) schema.GroupResource {
	return schema.GroupResource{Group: "dws", Resource: resource}
}

func (src *ClientMountList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("ClientMountList"), "ConvertTo")
}

func (dst *ClientMountList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("ClientMountList"), "ConvertFrom")
}

func (src *ComputesList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("ComputesList"), "ConvertTo")
}

func (dst *ComputesList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("ComputesList"), "ConvertFrom")
}

func (src *DWDirectiveRuleList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("DWDirectiveRuleList"), "ConvertTo")
}

func (dst *DWDirectiveRuleList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("DWDirectiveRuleList"), "ConvertFrom")
}

func (src *DirectiveBreakdownList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("DirectiveBreakdownList"), "ConvertTo")
}

func (dst *DirectiveBreakdownList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("DirectiveBreakdownList"), "ConvertFrom")
}

func (src *PersistentStorageInstanceList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("PersistentStorageInstanceList"), "ConvertTo")
}

func (dst *PersistentStorageInstanceList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("PersistentStorageInstanceList"), "ConvertFrom")
}

func (src *ServersList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("ServersList"), "ConvertTo")
}

func (dst *ServersList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("ServersList"), "ConvertFrom")
}

func (src *StorageList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("StorageList"), "ConvertTo")
}

func (dst *StorageList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("StorageList"), "ConvertFrom")
}

func (src *SystemConfigurationList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("SystemConfigurationList"), "ConvertTo")
}

func (dst *SystemConfigurationList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("SystemConfigurationList"), "ConvertFrom")
}

func (src *WorkflowList) ConvertTo(dstRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("WorkflowList"), "ConvertTo")
}

func (dst *WorkflowList) ConvertFrom(srcRaw conversion.Hub) error {
	return apierrors.NewMethodNotSupported(resource("WorkflowList"), "ConvertFrom")
}

// The conversion-gen tool dropped these from zz_generated.conversion.go to
// force us to acknowledge that we are addressing the conversion requirements.

func Convert_v1alpha5_ClientMountStatus_To_v1alpha1_ClientMountStatus(in *dwsv1alpha5.ClientMountStatus, out *ClientMountStatus, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_ClientMountStatus_To_v1alpha1_ClientMountStatus(in, out, s)
}

func Convert_v1alpha5_DirectiveBreakdownStatus_To_v1alpha1_DirectiveBreakdownStatus(in *dwsv1alpha5.DirectiveBreakdownStatus, out *DirectiveBreakdownStatus, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_DirectiveBreakdownStatus_To_v1alpha1_DirectiveBreakdownStatus(in, out, s)
}

func Convert_v1alpha1_ResourceErrorInfo_To_v1alpha5_ResourceErrorInfo(in *ResourceErrorInfo, out *dwsv1alpha5.ResourceErrorInfo, s apiconversion.Scope) error {
	return autoConvert_v1alpha1_ResourceErrorInfo_To_v1alpha5_ResourceErrorInfo(in, out, s)
}

func Convert_v1alpha5_ResourceErrorInfo_To_v1alpha1_ResourceErrorInfo(in *dwsv1alpha5.ResourceErrorInfo, out *ResourceErrorInfo, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_ResourceErrorInfo_To_v1alpha1_ResourceErrorInfo(in, out, s)
}

func Convert_v1alpha5_ServersStatus_To_v1alpha1_ServersStatus(in *dwsv1alpha5.ServersStatus, out *ServersStatus, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_ServersStatus_To_v1alpha1_ServersStatus(in, out, s)
}

func Convert_v1alpha1_WorkflowSpec_To_v1alpha5_WorkflowSpec(in *WorkflowSpec, out *dwsv1alpha5.WorkflowSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha1_WorkflowSpec_To_v1alpha5_WorkflowSpec(in, out, s)
}

func Convert_v1alpha5_WorkflowSpec_To_v1alpha1_WorkflowSpec(in *dwsv1alpha5.WorkflowSpec, out *WorkflowSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_WorkflowSpec_To_v1alpha1_WorkflowSpec(in, out, s)
}

func Convert_v1alpha1_SystemConfigurationSpec_To_v1alpha5_SystemConfigurationSpec(in *SystemConfigurationSpec, out *dwsv1alpha5.SystemConfigurationSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha1_SystemConfigurationSpec_To_v1alpha5_SystemConfigurationSpec(in, out, s)
}

func Convert_v1alpha5_SystemConfigurationSpec_To_v1alpha1_SystemConfigurationSpec(in *dwsv1alpha5.SystemConfigurationSpec, out *SystemConfigurationSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_SystemConfigurationSpec_To_v1alpha1_SystemConfigurationSpec(in, out, s)
}

func Convert_v1alpha5_StorageSpec_To_v1alpha1_StorageSpec(in *dwsv1alpha5.StorageSpec, out *StorageSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_StorageSpec_To_v1alpha1_StorageSpec(in, out, s)
}

func Convert_v1alpha5_WorkflowStatus_To_v1alpha1_WorkflowStatus(in *dwsv1alpha5.WorkflowStatus, out *WorkflowStatus, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_WorkflowStatus_To_v1alpha1_WorkflowStatus(in, out, s)
}

func Convert_v1alpha5_SystemConfigurationStatus_To_v1alpha1_SystemConfigurationStatus(in *dwsv1alpha5.SystemConfigurationStatus, out *SystemConfigurationStatus, s apiconversion.Scope) error {
	return autoConvert_v1alpha5_SystemConfigurationStatus_To_v1alpha1_SystemConfigurationStatus(in, out, s)
}
