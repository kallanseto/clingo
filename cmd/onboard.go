package cmd

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
)

// Project details for onboarding
type Project struct {
	Name         string `json: "name"`
	Team         string `json: "team"`
	Email        string `json: "email"`
	Owner        string `json: "owner"`
	Service      string `json: "service"`
	Application  string `json: "application"`
	Domain       string `json: "domain"`
	CPU          int    `json: "cpu"`
	Memory       int    `json: "memory"`
	Egressip     string `json: "egressip"`
	Netid        string `json: "netid"`
	Snatip       string `json: "snatip"`
	Namespacevip string `json: "namespacevip"`
}

func (p Project) newProjectDir(BASEDIR string) error {
	return os.Mkdir(BASEDIR+"/"+p.Name, 0755)
}

// EgressAllocations type for managing egressIP allocations
type EgressAllocations map[string]string

func netid(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	i := int(h.Sum32())

	rand.Seed(int64(i))

	return strconv.Itoa(rand.Intn(1000000))
}

func (p *Project) allocateEgressIP(filename string) {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	e := make(EgressAllocations)
	json.Unmarshal(bs, &e)

	if len(e) == 0 {
		fmt.Println("Info: No egressIPs configured for this cluster")
		return
	}

	allocated := false
	for k, v := range e {
		if v == "" && !allocated {
			e[k] = p.Name
			p.Egressip, p.Netid = k, netid(k)
			allocated = true
		}
	}

	if !allocated {
		log.Fatal("Error: cannot allocate egressIP")
	}

	bs, err = json.MarshalIndent(e, "", "")
	ioutil.WriteFile(filename, bs, 0644)

	return
}

func (p Project) writeNamespaceManifest(BASEDIR string) bool {
	m := `apiVersion: v1
kind: Namespace
metadata:
  annotations:
    collectord.io/index: openshift_` + p.Name + `
    ing.com.au/team: ` + p.Team + `
    ing.com.au/email: ` + p.Email + `
    ing.com.au/owner: ` + p.Owner + `
    ing.com.au/service: ` + p.Service + `
    ing.com.au/application: ` + p.Application + `
    ing.com.au/domain: ` + p.Domain + `
    ing.com.au/egressip: ` + p.Egressip + `
    ing.com.au/snatip: ` + p.Snatip + `
    ing.com.au/namespacevip: ` + p.Namespacevip + `
  labels:
    application: ` + p.Application + `
    service: ` + p.Service + `
    domain: ` + p.Domain + `
  name: ` + p.Name

	err := ioutil.WriteFile(BASEDIR+p.Name+"/namespace.yaml", []byte(m), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func (p Project) writeNetnamespaceManifest(BASEDIR string) bool {
	if p.Egressip == "" {
		return false
	}
	m := `apiVersion: network.openshift.io/v1
egressIPs:
- ` + p.Egressip + `
kind: NetNamespace
metadata:
  name: ` + p.Name + `
netname: ` + p.Name + `
netid: ` + p.Netid

	err := ioutil.WriteFile(BASEDIR+p.Name+"/netnamespace.yaml", []byte(m), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func (p Project) writeResourceQuota(BASEDIR string) bool {
	if p.CPU == 0 || p.Memory == 0 {
		return false
	}
	m := `apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute
spec:
  hard:
    cpu: ` + strconv.Itoa(p.CPU) + `
    memory: ` + strconv.Itoa(p.Memory) + "Gi"

	err := ioutil.WriteFile(BASEDIR+p.Name+"/resourcequota.yaml", []byte(m), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func (p Project) writeProjectManifests(BASEDIR string) bool {
	m := `namespace: ` + p.Name + `
commonLabels:
  prometheus: appdeployment
resources:
  - ../../../common/
  - namespace.yaml`

	if p.writeNetnamespaceManifest(BASEDIR) {
		m += `
  - netnamespace.yaml`
	}

	if p.writeResourceQuota(BASEDIR) {
		m += `
patches:
  - resourcequota.yaml
`
	}

	p.writeNamespaceManifest(BASEDIR)

	err := ioutil.WriteFile(BASEDIR+p.Name+"/kustomization.yaml", []byte(m), 0644)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func (p Project) updateClusterKustomization(BASEDIR string) {
	f, err := os.OpenFile(BASEDIR+"kustomization.yaml", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.WriteString("  - " + p.Name + "\n"); err != nil {
		log.Fatal(err)
	}
	return
}
