package main

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type ResourceJSONData struct {
	Resources []*Resource `json:"resource"`
}

type Resource struct {
	OVHDomainZoneRecord []OVHDomainZoneRecord `json:"ovh_domain_zone_record"`
}

type OVHDomainZoneRecord map[string]*DomainZoneRecord

// Ovh domain record add resource
func (j *ResourceJSONData) GetResource() *Resource {
	if len(j.Resources) > 0 {
		return j.Resources[0]
	}

	res := &Resource{}

	j.Resources = append(j.Resources, res)

	return res
}

func (j *ResourceJSONData) GetOVHfile() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get the current dir %s", err)
	}
	return path.Join(pwd, "infrastrure/terraform/modules/ovh", "resource_auto.tf.json")
}

func (j *ResourceJSONData) GetPromoxfile() string {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Cannot get the current dir %s", err)
	}
	return path.Join(pwd, "infrastrure/terraform/modules/proxmox", "resource_auto.tf.json")
}

func (r *ResourceJSONData) InitResourcesFiles() {
	// Create file if not exist
	for _, file := range r.GetFiles() {
		File.createIfNotExistsWithContent(file, "{}")
	}
}

func (r *ResourceJSONData) GetFiles() [2]string {
	return [...]string{r.GetOVHfile(), r.GetPromoxfile()}
}

func (r *ResourceJSONData) ParseOVHresourcesJSON() error {
	err := json.Unmarshal(File.readFile(r.GetOVHfile()), &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceJSONData) ParseProxmoxResourcesJSON() error {
	err := json.Unmarshal(File.readFile(r.GetPromoxfile()), &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceJSONData) WriteOVHresources() {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}
	var isEmpty = true

	for _, res := range r.Resources {
		if len(res.OVHDomainZoneRecord) == 0 {
			break
		} else {
			isEmpty = false
		}
	}

	if isEmpty {
		File.writeInFile(r.GetOVHfile(), "{}")
		return
	}
	File.writeInFile(r.GetOVHfile(), string(data))
}

func (r *ResourceJSONData) WriteProxmoxResources() {
	data, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	File.writeInFile(r.GetPromoxfile(), string(data))
}

// Resource methods
func (r *Resource) GetOVHDomainZoneRecord() OVHDomainZoneRecord {
	if len(r.OVHDomainZoneRecord) > 0 {
		return r.OVHDomainZoneRecord[0]
	}
	return nil
}

func (r *Resource) AddDomainZoneRerord(ref string, domain *DomainZoneRecord) {
	ovh_dzr := r.GetOVHDomainZoneRecord()

	if ovh_dzr == nil {
		r.OVHDomainZoneRecord = append(r.OVHDomainZoneRecord, OVHDomainZoneRecord{ref: domain})
		return
	}

	ovh_dzr[ref] = domain
}

func (r *Resource) DeleteDomainZoneRerord(ref string) {
	ozr := r.GetOVHDomainZoneRecord()
	delete(ozr, ref)

	if len(ozr) == 0 {
		r.OVHDomainZoneRecord = nil
	}
}

// Initiliaze some function
func init() {
	resource := &ResourceJSONData{}
	// Create file id not exist yet
	resource.InitResourcesFiles()
}
