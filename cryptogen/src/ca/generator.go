/*
Copyright IBM Corp. 2017 All Rights Reserved.

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
package ca

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"path/filepath"
	// "github.com/hyperledger/fabric/bccsp"
	// "github.com/hyperledger/fabric/bccsp/signer"
	// "github.com/hyperledger/fabric/bccsp/factory"

	"csp"
	"fmt"
	"reflect"
	"encoding/asn1"
	"errors"
)

type CA struct {
	Name string
	//SignKey  *ecdsa.PrivateKey
	Signer   crypto.Signer
	SignCert *x509.Certificate
}

// NewCA creates an instance of CA and saves the signing key pair in
// baseDir/name
func NewCA(baseDir, org, name string) (*CA, error) {

	var response error
	var ca *CA

	err := os.MkdirAll(baseDir, 0755)
	if err == nil {
		priv, signer, err := csp.GeneratePrivateKey(baseDir)
		response = err
		if err == nil {
			// get public signing certificate
			ecPubKey, err := csp.GetECPublicKey(priv)
			response = err
			if err == nil {
				template := x509Template()
				//this is a CA
				template.IsCA = true
				template.KeyUsage |= x509.KeyUsageDigitalSignature |
					x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
					x509.KeyUsageCRLSign
				template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}

				//set the organization for the subject
				subject := subjectTemplate()
				subject.Organization = []string{org}
				subject.CommonName = name

				template.Subject = subject
				template.SubjectKeyId = priv.SKI()

				x509Cert, err := genCertificateECDSA(baseDir, name, &template, &template,
					ecPubKey, signer)
				response = err
				if err == nil {
					ca = &CA{
						Name:     name,
						Signer:   signer,
						SignCert: x509Cert,
					}
				}
			}
		}
	}
	return ca, response
}

// 测试生成根证书
func NewCAOut(baseDir, org, name string, isCa bool) {
	err := os.MkdirAll(baseDir, 0755)
	if err == nil {
		priv, signer, err := csp.GeneratePrivateKey(baseDir)
		if err == nil {
			// get public signing certificate
			ecPubKey, err := csp.GetECPublicKey(priv)
			if err == nil {
				template := x509Template()
				//this is a CA
				template.IsCA = isCa
				template.KeyUsage |= x509.KeyUsageDigitalSignature |
					x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
					x509.KeyUsageCRLSign
				template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}

				//set the organization for the subject
				subject := subjectTemplate()
				subject.Organization = []string{org}
				subject.CommonName = name

				template.Subject = subject
				template.SubjectKeyId = priv.SKI()

				genCertificateECDSA(baseDir, name, &template, &template,
					ecPubKey, signer)
			}
		}
	}
}

// 从根证书文件读取的bytes来获取CA对象
func NewCAFromBytes(baseDir, org, name string, certs []byte) (*CA, error) {
	var response error
	var ca *CA
	certBlock, _ := pem.Decode(certs)
	//fmt.Printf("certBlock=%v", certBlock.Bytes)
	parCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		fmt.Printf("parse certs file error:%v\r\n", err)
		return nil, err
	}
	test := isCACert(parCert) 
	fmt.Printf("test = %v, parisca=%v, keyusage=%v\r\n", test, parCert.IsCA, parCert.KeyUsage)
	// parCert.IsCA = true
	// parCert.KeyUsage |= x509.KeyUsageDigitalSignature |
	// 	x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
	// 	x509.KeyUsageCRLSign
	// parCert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}
	// tca, err := NewCA("/tmp", org, name)
	// parCert := tca.SignCert

	err = os.MkdirAll(baseDir, 0755)
	if err == nil {
		priv, signer, err := csp.GeneratePrivateKey(baseDir)
		response = err
		if err == nil {
			// get public signing certificate
			ecPubKey, err := csp.GetECPublicKey(priv)
			response = err
			if err == nil {
				template := x509Template()
				//this is a CA
				template.IsCA = true
				template.KeyUsage |= x509.KeyUsageDigitalSignature |
					x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign |
					x509.KeyUsageCRLSign
				template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageAny}

				//set the organization for the subject
				subject := subjectTemplate()
				subject.Organization = []string{org}
				subject.CommonName = name

				template.Subject = subject
				template.SubjectKeyId = priv.SKI()

				//fmt.Printf("template subject key id=%v, keyusage=%v\r\n", template.SubjectKeyId, template.KeyUsage)

				// x509Cert, err := genCertificateECDSA(baseDir, name, &template, parCert,
				// 	ecPubKey, signer)
				x509Cert, err := genCertificateECDSA(baseDir, name, &template, parCert,
					ecPubKey, signer)

				// fmt.Printf("my parisca=%v, keyusage=%v\r\n", x509Cert.IsCA, x509Cert.KeyUsage)
				// fmt.Printf("my cert ca test start ------------\r\n")
				// test = isCACert(x509Cert)
				// fmt.Printf("my cert ca test end --------------\r\n")


				response = err
				if err == nil {
					ca = &CA{
						Name:     name,
						Signer:   signer,
						SignCert: x509Cert,
					}
				}
			}
		}
	}

	// opts := &factory.FactoryOpts{
	// 	ProviderName: "SW",
	// 	SwOpts: &factory.SwOpts{
	// 		HashFamily: "SHA2",
	// 		SecLevel:   256,

	// 		FileKeystore: &factory.FileKeystoreOpts{
	// 			KeyStorePath: baseDir,
	// 		},
	// 	},
	// }
	// csp, err := factory.GetBCCSPFromOpts(opts)
	// if err != nil {
	// 	fmt.Printf("get csp error:%v\r\n", err)
	// 	return nil, err
	// }
	// priv, err := csp.KeyGen(&bccsp.ECDSAP256KeyGenOpts{Temporary: false})
	// if err != nil {
	// 	fmt.Printf("ken gen error:%v\r\n", err)
	// 	return nil, err
	// }

	// // create a crypto.Signer
	// sig, err := signer.New(csp, priv)
	// if err != nil {
	// 	fmt.Printf("create signer err:%v\r\n", err)
	// 	return nil, err
	// }
	// ca = &CA{
	// 	Name: name, 
	// 	Signer: sig,
	// 	SignCert: x509Cert,
	// }
	return ca, response
}

// SignCertificate creates a signed certificate based on a built-in template
// and saves it in baseDir/name
func (ca *CA) SignCertificate(baseDir, name string, sans []string, pub *ecdsa.PublicKey,
	ku x509.KeyUsage, eku []x509.ExtKeyUsage) (*x509.Certificate, error) {

	template := x509Template()
	template.KeyUsage = ku
	template.ExtKeyUsage = eku

	//set the organization for the subject
	subject := subjectTemplate()
	subject.CommonName = name

	template.Subject = subject
	template.DNSNames = sans

	cert, err := genCertificateECDSA(baseDir, name, &template, ca.SignCert,
		pub, ca.Signer)

	if err != nil {
		return nil, err
	}

	return cert, nil
}

// default template for X509 subject
func subjectTemplate() pkix.Name {
	return pkix.Name{
		Country:  []string{"US"},
		Locality: []string{"San Francisco"},
		Province: []string{"California"},
	}
}

// default template for X509 certificates
func x509Template() x509.Certificate {

	//generate a serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, serialNumberLimit)

	now := time.Now()
	//basic template to use
	x509 := x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             now,
		NotAfter:              now.Add(3650 * 24 * time.Hour), //~ten years
		BasicConstraintsValid: true,
	}
	return x509

}

// generate a signed X509 certficate using ECDSA
func genCertificateECDSA(baseDir, name string, template, parent *x509.Certificate, pub *ecdsa.PublicKey,
	priv interface{}) (*x509.Certificate, error) {

	//create the x509 public cert
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, pub, priv)
	if err != nil {
		return nil, err
	}

	//write cert out to file
	fileName := filepath.Join(baseDir, name+"-cert.pem")
	certFile, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	//pem encode the cert
	err = pem.Encode(certFile, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certFile.Close()
	if err != nil {
		return nil, err
	}

	x509Cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, err
	}
	return x509Cert, nil
}





/// 以下为测试代码
// isCACert does a few checks on the certificate,
// assuming it's a CA; it returns true if all looks good
// and false otherwise
func isCACert(cert *x509.Certificate) bool {
	_, err := getSubjectKeyIdentifierFromCert(cert)
	if err != nil {
		return false
	}

	if !cert.IsCA {
		return false
	}

	return true
}

// getSubjectKeyIdentifierFromCert returns the Subject Key Identifier for the supplied certificate
// Subject Key Identifier is an identifier of the public key of this certificate
func getSubjectKeyIdentifierFromCert(cert *x509.Certificate) ([]byte, error) {
	var SKI []byte
	for _, ext := range cert.Extensions {
		fmt.Printf("ext.Id=%v\r\n", ext.Id)
		// Subject Key Identifier is identified by the following ASN.1 tag
		// subjectKeyIdentifier (2 5 29 14) (see https://tools.ietf.org/html/rfc3280.html)
		if reflect.DeepEqual(ext.Id, asn1.ObjectIdentifier{2, 5, 29, 14}) {
			fmt.Printf("equal ok\r\n")
			_, err := asn1.Unmarshal(ext.Value, &SKI)
			if err != nil {
				return nil, fmt.Errorf("Failed to unmarshal Subject Key Identifier, err %s", err)
			}

			return SKI, nil
		}
	}
	fmt.Printf("oh no\r\n")
	return nil, errors.New("subjectKeyIdentifier not found in certificate")
}