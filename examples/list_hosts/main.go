package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/infra-monkey/go-freeipa/freeipa"
)

func main() {
	krb5Principal := os.Getenv("IPA_USERNAME")
	krb5Realm := os.Getenv("IPA_REALM")

	krb5KtFd, err := os.Open(os.Getenv("IPA_KEYTAB"))
	if err != nil {
		log.Fatal(err)
	}
	defer krb5KtFd.Close()

	krb5Fd, err := os.Open("/etc/krb5.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer krb5Fd.Close()

	krb5ConnectOption := &freeipa.KerberosConnectOptions{
		Krb5ConfigReader: krb5Fd,
		KeytabReader:     krb5KtFd,
		Username:         krb5Principal,
		Realm:            krb5Realm,
	}

	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // WARNING DO NOT USE THIS OPTION IN PRODUCTION
		},
	}

	c, err := freeipa.ConnectWithKerberos(os.Getenv("IPA_HOST"), tspt, krb5ConnectOption)
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.HostFind("", &freeipa.HostFindArgs{}, &freeipa.HostFindOptionalArgs{
		Sizelimit: freeipa.Int(0),
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, host := range res.Result {
		fmt.Printf("host[%s]\n", host.Fqdn)
	}
}
