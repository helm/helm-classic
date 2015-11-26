package codec

import (
	"k8s.io/kubernetes/pkg/api/v1"
	// initialise the extensions scheme
	_ "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	oauth "github.com/openshift/origin/pkg/oauth/api/v1"
)

// Object describes some discrete thing decoded from YAML or JSON.
type Object struct {
	data []byte
	dec  DecodeFunc
}

// DecodeFunc is a func that can decode bytes into an interface.
type DecodeFunc func([]byte, interface{}) error

// Ref returns an ObjectReference with basic information about the object.
//
// This can be used to perform simple operations, as well as to instrospect
// a record enough to know how to unmarshal it.
func (m *Object) Ref() (*v1.ObjectReference, error) {
	or := &v1.ObjectReference{}
	return or, m.dec(m.data, or)
}

// Object decodes the manifest into the given object.
//
// You can use ObjectReference.Kind to figure out what kind of object to
// decode into.
//
// There are several shortcut methods that will allow you to decode directly
// to one of the common types, like Pod(), RC(), and Service().
func (m *Object) Object(v interface{}) error {
	return m.dec(m.data, v)
}

// Pod decodes a manifest into a Pod.
func (m *Object) Pod() (*v1.Pod, error) {
	o := new(v1.Pod)
	return o, m.Object(o)
}

// RC decodes a manifest into a ReplicationController.
func (m *Object) RC() (*v1.ReplicationController, error) {
	o := new(v1.ReplicationController)
	return o, m.Object(o)
}

// Service decodes a manifest into a Service
func (m *Object) Service() (*v1.Service, error) {
	o := new(v1.Service)
	return o, m.Object(o)
}

// PersistentVolume decodes a manifest into a PersistentVolume
func (m *Object) PersistentVolume() (*v1.PersistentVolume, error) {
	o := new(v1.PersistentVolume)
	return o, m.Object(o)
}

// Secret decodes a manifest into a Secret
func (m *Object) Secret() (*v1.Secret, error) {
	o := new(v1.Secret)
	return o, m.Object(o)
}

// Namespace decodes a manifest into a Namespace
func (m *Object) Namespace() (*v1.Namespace, error) {
	o := new(v1.Namespace)
	return o, m.Object(o)
}

// ServiceAccount decodes a manifest into a ServiceAccount.
func (m *Object) ServiceAccount() (*v1.ServiceAccount, error) {
	o := new(v1.ServiceAccount)
	return o, m.Object(o)
}

// OAuthClient decodes a manifest into a OAuthClient.
func (m *Object) OAuthClient() (*oauth.OAuthClient, error) {
	o := new(oauth.OAuthClient)
	return o, m.Object(o)
}
