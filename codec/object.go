package codec

import (
	"bytes"
	"errors"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	// initialise the extensions scheme
	_ "k8s.io/kubernetes/pkg/apis/extensions/v1beta1"
	oapi "github.com/openshift/origin/pkg/template/api/v1"
	oauth "github.com/openshift/origin/pkg/oauth/api/v1"
)

// Object describes some discrete instance of a Kind.
//
// Objects are tracked as raw []byte data. They can be decoded into a variety
// of forms, but the decoded version is never stored on the object itself.
//
// Mutators like AddLabels typically decode the object, mutate it, and then
// re-encode it. The reason for this is that those operations can assume
// a very simply and generic type and trust that the decoder/encoder pair can
// preserve structure.
//
// Some operations, like Ref() and Meta(), can produce structured subsets of
// the original data, but will not preserve the entire object. These can be
// used to read data easily, but should not be used to mutate an object.
//
// Operations like JSON() and YAML() will decode the data into a generic form
// and then re-encode the data.
type Object struct {
	data []byte
	dec  DecodeFunc
}

// Metadata provides just the basic metadata fields of an object.
//
// It contains more data than an ObjectReference, since it includes standard
// metadata values like Name, Labels, and Annotations.
//
// This is a readable structure, but should not be used to write changes.
type Metadata struct {
	unversioned.TypeMeta `json:",inline"`
	api.ObjectMeta       `json:"metadata,omitempty"`
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

// YAML takes the raw data and (re-)encodes it as YAML.
func (m *Object) YAML() ([]byte, error) {
	// We are not assuming that the m.data is YAML, though for the current
	// iteration of Helm that assumption probably holds. Instead, we are
	// relying upon the decoder to decode to a generic format, and then we
	// re-encode it into YAML.
	var d interface{}

	if err := m.Object(&d); err != nil {
		return []byte{}, err
	}

	b := new(bytes.Buffer)
	err := YAML.Encode(b).One(&d)
	return b.Bytes(), err
}

// JSON re-encodes the data on this object into JSON format.
func (m *Object) JSON() ([]byte, error) {
	var d interface{}
	if err := m.Object(&d); err != nil {
		return []byte{}, err
	}
	b := new(bytes.Buffer)
	err := JSON.Encode(b).One(&d)
	return b.Bytes(), err
}

// Meta returns a *Metadata
//
// This contains more information than an ObjectReference. It is valid for the
// core Kubernetes kinds, but is not guranteed to work for all possible kinds.
func (m *Object) Meta() (*Metadata, error) {
	u := &Metadata{}
	err := m.dec(m.data, u)
	return u, err
}

// AddLabels adds the map of labels to an object, regardless of kind.
//
// This looks for a top-level metadata entry, and adds labels there. Because
// this decodes the object into a generic format and then re-encodes it, this
// is a more costly operation than unmarshaling into a specific type and
// mutating the properties. However, it works on all types.
func (m *Object) AddLabels(labels map[string]string) error {
	return m.addMDItem("labels", labels)
}

// AddAnnotations adds the map of annotations to an object, regardless of kind.
//
// This looks for a top-level metadata entry and adds annotations inside of that.
//
// See the notes on AddLabels for performance implications.
func (m *Object) AddAnnotations(ann map[string]string) error {
	return m.addMDItem("annotations", ann)
}

// addMDItem adds the given key/hash combo to a generic object.
//
// TODO: In the future we might want to make this more flexible. If it turns
// out that adding annotations and labels in close sequence is a common thing,
// we should facilitate that.
func (m *Object) addMDItem(key string, value map[string]string) error {
	var d interface{}
	if err := m.dec(m.data, &d); err != nil {
		return err
	}
	if val, ok := d.(map[string]interface{}); ok {
		md, ok := val["metadata"]
		if !ok {
			val["metadata"] = map[string]interface{}{
				key: value,
			}
		} else if mmd, ok := md.(map[string]interface{}); ok {
			if l, ok := mmd[key]; ok {
				ll := l.(map[string]interface{})
				for k, v := range value {
					ll[k] = v
				}
			} else {
				mmd[key] = value
			}
		}
	} else {
		return errors.New("Top level object is not a map")
	}

	var b bytes.Buffer
	if err := YAML.Encode(&b).One(d); err != nil {
		return err
	}

	m.data = b.Bytes()

	return nil
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

// DaemonSet decodes a manifest into a DaemonSet.
func (m *Object) DaemonSet() (*v1beta1.DaemonSet, error) {
	o := new(v1beta1.DaemonSet)
	return o, m.Object(o)
}

// Job decodes a manifest into a Job.
func (m *Object) Job() (*v1beta1.Job, error) {
	o := new(v1beta1.Job)
	return o, m.Object(o)
}

// Ingress decodes a manifest into a Ingress.
func (m *Object) Ingress() (*v1beta1.Ingress, error) {
	o := new(v1beta1.Ingress)
	return o, m.Object(o)
}

// Deployment decodes a manifest into a Deployment.
func (m *Object) Deployment() (*v1beta1.Deployment, error) {
	o := new(v1beta1.Deployment)
	return o, m.Object(o)
}

// HorizontalPodAutoscaler decodes a manifest into a HorizontalPodAutoscaler.
func (m *Object) HorizontalPodAutoscaler() (*v1beta1.HorizontalPodAutoscaler, error) {
	o := new(v1beta1.HorizontalPodAutoscaler)
	return o, m.Object(o)
}

// List decodes a manifest into a List
func (m *Object) List() (*v1.List, error) {
	o := new(v1.List)
	return o, m.Object(o)
}

// Template decodes a manifest into a Template
func (m *Object) Template() (*oapi.Template, error) {
	o := new(oapi.Template)
	return o, m.Object(o)
}

// OAuthClient decodes a manifest into a OAuthClient
func (m *Object) OAuthClient() (*oauth.OAuthClient, error) {
	o := new(oauth.OAuthClient)
	return o, m.Object(o)
}
