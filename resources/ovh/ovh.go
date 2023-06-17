package ovh

import (
	"encoding/json"
	"log"
	"os"
	"path"

	files "smatflow/platform-installer/files"
	structs "smatflow/platform-installer/structs"
)

type ResourceJSONData struct {
	Resources []*Resource `json:"resource"`
}

type Resource struct {
	OVHDomainZoneRecord []OVHDomainZoneRecord `json:"ovh_domain_zone_record"`
}

type OVHDomainZoneRecord map[string]*structs.DomainZoneRecord

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

func (r *ResourceJSONData) InitResourcesFiles() {
	// Create file if not exist
	for _, file := range r.GetFiles() {
		files.CreateIfNotExistsWithContent(file, "{}")
	}
}

func (r *ResourceJSONData) GetFiles() [1]string {
	return [...]string{r.GetOVHfile()}
}

func (r *ResourceJSONData) ParseOVHresourcesJSON() error {
	err := json.Unmarshal(files.ReadFile(r.GetOVHfile()), &r)
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
		files.WriteInFile(r.GetOVHfile(), "{}")
		return
	}
	files.WriteInFile(r.GetOVHfile(), string(data))
}

// Resource methods
func (r *Resource) GetOVHDomainZoneRecord() OVHDomainZoneRecord {
	if len(r.OVHDomainZoneRecord) > 0 {
		return r.OVHDomainZoneRecord[0]
	}
	return nil
}

func (r *Resource) AddDomainZoneRerord(ref string, domain *structs.DomainZoneRecord) {
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
