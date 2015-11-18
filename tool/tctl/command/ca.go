/*
Copyright 2015 Gravitational, Inc.

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
package command

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gravitational/teleport/lib/auth/native"
	"github.com/gravitational/teleport/lib/services"

	"github.com/gravitational/teleport/Godeps/_workspace/src/github.com/buger/goterm"
	"github.com/gravitational/teleport/Godeps/_workspace/src/github.com/gravitational/trace"
)

func (cmd *Command) GenerateKeyPair(privateKeyPath, publicKeyPath, passphrase string) error {
	priv, pub, err := native.New().GenerateKeyPair(passphrase)
	if err != nil {
		return trace.Wrap(err)
	}
	if err := ioutil.WriteFile(privateKeyPath, priv, 0600); err != nil {
		return trace.Wrap(err)
	}
	if err := ioutil.WriteFile(publicKeyPath, pub, 0666); err != nil {
		return trace.Wrap(err)
	}
	cmd.printOK("Public and private keys have been written")
	return nil
}

func (cmd *Command) ResetHostCertificateAuthority(confirm bool) {
	if !confirm && !cmd.confirm("Reseting private and public keys for Host CA. This will invalidate all signed host certs. Continue?") {
		cmd.printError(fmt.Errorf("aborted by user"))
		return
	}
	if err := cmd.client.ResetHostCertificateAuthority(); err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("Certificate authority keys have been regenerated")
}

func (cmd *Command) GetHostPublicCertificate() {
	key, err := cmd.client.GetHostPublicCertificate()
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("Host CA Key")
	fmt.Fprintf(cmd.out, string(key.PubValue))
}

func (cmd *Command) ResetUserCertificateAuthority(confirm bool) {
	if !confirm && !cmd.confirm("Reseting private and public keys for User CA. This will invalidate all signed user certs. Continue?") {
		cmd.printError(fmt.Errorf("aborted by user"))
		return
	}
	if err := cmd.client.ResetUserCertificateAuthority(); err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("Certificate authority keys have been regenerated")
}

func (cmd *Command) GetUserPublicCertificate() {
	key, err := cmd.client.GetUserPublicCertificate()
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("User CA Key")
	fmt.Fprintf(cmd.out, string(key.PubValue))
}

func (cmd *Command) UpsertRemoteCertificate(id, fqdn, certType, path string, ttl time.Duration) {
	val, err := cmd.readInput(path)
	if err != nil {
		cmd.printError(err)
		return
	}
	cert := services.PublicCertificate{
		FQDN:     fqdn,
		Type:     certType,
		ID:       id,
		PubValue: val,
	}
	if err := cmd.client.UpsertRemoteCertificate(cert, ttl); err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("Remote cert have been upserted")
}

func (cmd *Command) GetRemoteCertificates(fqdn, certType string) {
	certs, err := cmd.client.GetRemoteCertificates(certType, fqdn)
	if err != nil {
		cmd.printError(err)
		return
	}
	fmt.Fprintf(cmd.out, remoteCertsView(certs))
}

func (cmd *Command) DeleteRemoteCertificate(id, fqdn, certType string) {
	err := cmd.client.DeleteRemoteCertificate(certType, fqdn, id)
	if err != nil {
		cmd.printError(err)
		return
	}
	cmd.printOK("certificate deleted")
}

func remoteCertsView(certs []services.PublicCertificate) string {
	t := goterm.NewTable(0, 10, 5, ' ', 0)
	fmt.Fprint(t, "Type\tFQDN\tID\tValue\n")
	if len(certs) == 0 {
		return t.String()
	}
	for _, c := range certs {
		fmt.Fprintf(t, "%v\t%v\t%v\t%v\n", c.Type, c.FQDN, c.ID, string(c.PubValue))
	}
	return t.String()
}
