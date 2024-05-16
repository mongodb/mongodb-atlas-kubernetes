//go:build !ignore_autogenerated

/*
Copyright (C) MongoDB, Inc. 2020-present.

Licensed under the Apache License, Version 2.0 (the "License"); you may
not use this file except in compliance with the License. You may obtain
a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasAuditing) DeepCopyInto(out *AtlasAuditing) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasAuditing.
func (in *AtlasAuditing) DeepCopy() *AtlasAuditing {
	if in == nil {
		return nil
	}
	out := new(AtlasAuditing)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AtlasAuditing) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasAuditingConfig) DeepCopyInto(out *AtlasAuditingConfig) {
	*out = *in
	if in.AuditFilter != nil {
		in, out := &in.AuditFilter, &out.AuditFilter
		*out = new(v1.JSON)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasAuditingConfig.
func (in *AtlasAuditingConfig) DeepCopy() *AtlasAuditingConfig {
	if in == nil {
		return nil
	}
	out := new(AtlasAuditingConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasAuditingSpec) DeepCopyInto(out *AtlasAuditingSpec) {
	*out = *in
	in.AtlasAuditingConfig.DeepCopyInto(&out.AtlasAuditingConfig)
	if in.ProjectIDs != nil {
		in, out := &in.ProjectIDs, &out.ProjectIDs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasAuditingSpec.
func (in *AtlasAuditingSpec) DeepCopy() *AtlasAuditingSpec {
	if in == nil {
		return nil
	}
	out := new(AtlasAuditingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasProject) DeepCopyInto(out *AtlasProject) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasProject.
func (in *AtlasProject) DeepCopy() *AtlasProject {
	if in == nil {
		return nil
	}
	out := new(AtlasProject)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AtlasProject) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasProjectList) DeepCopyInto(out *AtlasProjectList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AtlasProject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasProjectList.
func (in *AtlasProjectList) DeepCopy() *AtlasProjectList {
	if in == nil {
		return nil
	}
	out := new(AtlasProjectList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AtlasProjectList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AtlasProjectSpec) DeepCopyInto(out *AtlasProjectSpec) {
	*out = *in
	in.AtlasProjectSpec.DeepCopyInto(&out.AtlasProjectSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AtlasProjectSpec.
func (in *AtlasProjectSpec) DeepCopy() *AtlasProjectSpec {
	if in == nil {
		return nil
	}
	out := new(AtlasProjectSpec)
	in.DeepCopyInto(out)
	return out
}
