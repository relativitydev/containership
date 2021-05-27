// +build !ignore_autogenerated

/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1beta2

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerManagementObject) DeepCopyInto(out *ContainerManagementObject) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerManagementObject.
func (in *ContainerManagementObject) DeepCopy() *ContainerManagementObject {
	if in == nil {
		return nil
	}
	out := new(ContainerManagementObject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ContainerManagementObject) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerManagementObjectList) DeepCopyInto(out *ContainerManagementObjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ContainerManagementObject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerManagementObjectList.
func (in *ContainerManagementObjectList) DeepCopy() *ContainerManagementObjectList {
	if in == nil {
		return nil
	}
	out := new(ContainerManagementObjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ContainerManagementObjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerManagementObjectSpec) DeepCopyInto(out *ContainerManagementObjectSpec) {
	*out = *in
	if in.Images != nil {
		in, out := &in.Images, &out.Images
		*out = make([]Image, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerManagementObjectSpec.
func (in *ContainerManagementObjectSpec) DeepCopy() *ContainerManagementObjectSpec {
	if in == nil {
		return nil
	}
	out := new(ContainerManagementObjectSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerManagementObjectStatus) DeepCopyInto(out *ContainerManagementObjectStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerManagementObjectStatus.
func (in *ContainerManagementObjectStatus) DeepCopy() *ContainerManagementObjectStatus {
	if in == nil {
		return nil
	}
	out := new(ContainerManagementObjectStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Image) DeepCopyInto(out *Image) {
	*out = *in
	if in.SupportedTags != nil {
		in, out := &in.SupportedTags, &out.SupportedTags
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Destinations != nil {
		in, out := &in.Destinations, &out.Destinations
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Image.
func (in *Image) DeepCopy() *Image {
	if in == nil {
		return nil
	}
	out := new(Image)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistriesConfig) DeepCopyInto(out *RegistriesConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistriesConfig.
func (in *RegistriesConfig) DeepCopy() *RegistriesConfig {
	if in == nil {
		return nil
	}
	out := new(RegistriesConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RegistriesConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistriesConfigList) DeepCopyInto(out *RegistriesConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]RegistriesConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistriesConfigList.
func (in *RegistriesConfigList) DeepCopy() *RegistriesConfigList {
	if in == nil {
		return nil
	}
	out := new(RegistriesConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RegistriesConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistriesConfigSpec) DeepCopyInto(out *RegistriesConfigSpec) {
	*out = *in
	if in.Registries != nil {
		in, out := &in.Registries, &out.Registries
		*out = make([]Registry, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistriesConfigSpec.
func (in *RegistriesConfigSpec) DeepCopy() *RegistriesConfigSpec {
	if in == nil {
		return nil
	}
	out := new(RegistriesConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RegistriesConfigStatus) DeepCopyInto(out *RegistriesConfigStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RegistriesConfigStatus.
func (in *RegistriesConfigStatus) DeepCopy() *RegistriesConfigStatus {
	if in == nil {
		return nil
	}
	out := new(RegistriesConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Registry) DeepCopyInto(out *Registry) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Registry.
func (in *Registry) DeepCopy() *Registry {
	if in == nil {
		return nil
	}
	out := new(Registry)
	in.DeepCopyInto(out)
	return out
}