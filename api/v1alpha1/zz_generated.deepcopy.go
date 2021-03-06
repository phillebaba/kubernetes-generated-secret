// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecret) DeepCopyInto(out *GeneratedSecret) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecret.
func (in *GeneratedSecret) DeepCopy() *GeneratedSecret {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecret)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GeneratedSecret) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecretData) DeepCopyInto(out *GeneratedSecretData) {
	*out = *in
	if in.Length != nil {
		in, out := &in.Length, &out.Length
		*out = new(int)
		**out = **in
	}
	if in.Exclude != nil {
		in, out := &in.Exclude, &out.Exclude
		*out = make([]CharacterOption, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecretData.
func (in *GeneratedSecretData) DeepCopy() *GeneratedSecretData {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecretData)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecretList) DeepCopyInto(out *GeneratedSecretList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GeneratedSecret, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecretList.
func (in *GeneratedSecretList) DeepCopy() *GeneratedSecretList {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecretList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GeneratedSecretList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecretSpec) DeepCopyInto(out *GeneratedSecretSpec) {
	*out = *in
	in.SecretMeta.DeepCopyInto(&out.SecretMeta)
	if in.DataList != nil {
		in, out := &in.DataList, &out.DataList
		*out = make([]GeneratedSecretData, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecretSpec.
func (in *GeneratedSecretSpec) DeepCopy() *GeneratedSecretSpec {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecretSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GeneratedSecretStatus) DeepCopyInto(out *GeneratedSecretStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GeneratedSecretStatus.
func (in *GeneratedSecretStatus) DeepCopy() *GeneratedSecretStatus {
	if in == nil {
		return nil
	}
	out := new(GeneratedSecretStatus)
	in.DeepCopyInto(out)
	return out
}
