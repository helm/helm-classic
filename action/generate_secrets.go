package action

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/helm/helm/log"

	"golang.org/x/crypto/ssh"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/api/v1"
)

// secretSettings contains the flags to determine how the Secret generation works
type secretSettings struct {
	PrintImportFolders bool
	WriteGeneratedKeys bool
	GenerateSecretsData bool
}

// keypair is an SSH private and public pair
type keypair struct {
	pub  []byte
	priv []byte
}


// createSecretsFromAnnotations creates all the secrets referenced in the ReplicationController's annotations
// using these annotations and rules:
// 	https://github.com/fabric8io/fabric8/blob/master/docs/secretAnnotations.md
func createSecretsFromAnnotations(rc *v1.ReplicationController, namespace string, dry bool, secretFlags *secretSettings) error {
	for secretType, secretDataIdentifiers := range rc.Spec.Template.Annotations {
		if err := createSecretsFromAnnotation(secretDataIdentifiers, secretType, namespace, dry, secretFlags); err != nil {
			return err
		}
	}
	return nil
}

func createSecretsFromAnnotation(secretDataIdentifiers string, secretType string, namespace string, dry bool, secretFlags *secretSettings) error {
	var dataType = strings.Split(secretType, "/")
	if len(dataType) == 2 && dataType[0] == "fabric8.io" {
		switch dataType[1] {
		case "secret-ssh-key":
			// check to see if multiple public and private keys are needed
			items := strings.Split(secretDataIdentifiers, ",")
			for i := range items {
				var name = items[i]
				if err := marshalAndCreateSecret(namespace, name, secretType, nil, dry, secretFlags); err != nil {
					return err
				}
			}
		case "secret-ssh-public-key":
			// if this is just a public key then the secret name is at the start of the string
			f := func(c rune) bool {
				return c == ',' || c == '[' || c == ']'
			}
			secrets := strings.FieldsFunc(secretDataIdentifiers, f)
			numOfSecrets := len(secrets)

			var keysNames []string
			if numOfSecrets > 0 {
				// if multiple secrets
				for i := 1; i < numOfSecrets; i++ {
					keysNames = append(keysNames, secrets[i])
				}
			} else {
				// only single secret required
				keysNames[0] = "ssh-key.pub"
			}

			if err := marshalAndCreateSecret(namespace, secrets[0], secretType, keysNames, dry, secretFlags); err != nil {
				return err
			}
		case "secret-gpg-key:":
			gpgKeyName := []string{"gpg.conf", "secring.gpg", "pubring.gpg", "trustdb.gpg"}
			if err := marshalAndCreateSecret(namespace, secretDataIdentifiers, secretType, gpgKeyName, dry, secretFlags); err != nil {
				return err
			}
		}
	}
	return nil
}

func marshalAndCreateSecret(namespace string, name string, secretType string, keysNames []string, dry bool, secretFlags *secretSettings) error {
	secret := api.Secret{
		TypeMeta: unversioned.TypeMeta{
			Kind: "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name: name,
		},
		Type: api.SecretType(secretType),
		Data: getSecretData(secretType, name, keysNames, secretFlags),
	}
	return marshalAndKubeCtlCreate(secret, namespace, dry)
}

func getSecretData(secretType string, name string, keysNames []string, secretFlags *secretSettings) map[string][]byte {
	var dataType = strings.Split(secretType, "/")
	var data = make(map[string][]byte)

	switch dataType[1] {
	case "secret-ssh-key":
		if secretFlags.PrintImportFolders {
			logSecretImport(name + "/ssh-key")
			logSecretImport(name + "/ssh-key.pub")
		}

		sshKey, err1 := ioutil.ReadFile(name + "/ssh-key")
		sshKeyPub, err2 := ioutil.ReadFile(name + "/ssh-key.pub")

		// if we cant find the public and private key to import, and generation flag is set then lets generate the keys
		if (err1 != nil && err2 != nil) && secretFlags.GenerateSecretsData {
			log.Info("No secrets found on local filesystem, generating SSH public and private key pair\n")
			keypair := generateSSHKeyPair()
			if secretFlags.WriteGeneratedKeys {
				writeFile(name+"/ssh-key", keypair.priv)
				writeFile(name+"/ssh-key.pub", keypair.pub)
			}
			data["ssh-key"] = keypair.priv
			data["ssh-key.pub"] = keypair.pub

		} else if (err1 != nil || err2 != nil) && secretFlags.GenerateSecretsData {
			log.Info("Found some keys to import but with errors so unable to generate SSH public and private key pair. %s\n", name)
			check(err1)
			check(err2)
		} else {
			// if we're not generating the keys and there's an error importing them then still create the secret but with empty data
			check(err1)
			check(err2)

			data["ssh-key"] = sshKey
			data["ssh-key.pub"] = sshKeyPub
		}
	case "secret-ssh-public-key":
		for i := 0; i < len(keysNames); i++ {
			if secretFlags.PrintImportFolders {
				logSecretImport(name + "/" + keysNames[i])
			}

			sshPub, err := ioutil.ReadFile(name + "/" + keysNames[i])
			// if we cant find the public key to import and generation flag is set then lets generate the key
			if (err != nil) && secretFlags.GenerateSecretsData {
				log.Info("No secrets found on local filesystem, generating SSH public key\n")
				keypair := generateSSHKeyPair()
				if secretFlags.WriteGeneratedKeys {
					writeFile(name+"/ssh-key.pub", keypair.pub)
				}
				data[keysNames[i]] = keypair.pub

			} else {
				// if we're not generating the keys and there's an error importing them then still create the secret but with empty data
				check(err)
				data[keysNames[i]] = sshPub
			}
		}
	case "secret-gpg-key":
		for i := 0; i < len(keysNames); i++ {
			if secretFlags.PrintImportFolders {
				logSecretImport(name + "/" + keysNames[i])
			}
			gpg, err := ioutil.ReadFile(name + "/" + keysNames[i])
			check(err)

			data[keysNames[i]] = gpg
		}
	default:
		log.Warn("No matching data type %s\n", dataType)
	}
	return data
}

func logSecretImport(file string) {
	log.Info("Importing secret: %s\n", file)
}

func check(e error) {
	if e != nil {
		log.Warn("Warning: %s\n", e)
	}
}

func generateSSHKeyPair() keypair {
	priv, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		log.Die("Error generating key", err)
	}

	// Get der format. priv_der []byte
	privDer := x509.MarshalPKCS1PrivateKey(priv)

	// pem.Block
	// blk pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}

	// Resultant private key in PEM format.
	// priv_pem string
	privDem := string(pem.EncodeToMemory(&privBlock))

	// Public Key generation
	sshPublicKey, err := ssh.NewPublicKey(&priv.PublicKey)
	pubBytes := ssh.MarshalAuthorizedKey(sshPublicKey)

	return keypair{
		pub:  []byte(pubBytes),
		priv: []byte(privDem),
	}
}

func writeFile(path string, contents []byte) {
	dir := strings.Split(path, string(filepath.Separator))
	os.MkdirAll("."+string(filepath.Separator)+dir[0], 0700)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err := f.Write(contents); err != nil {
		panic(err)
	}
}