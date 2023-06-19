package ovh

import (
	"encoding/json"
	"log"
	"path"
	"smatflow/platform-installer/lib/files"
	"smatflow/platform-installer/lib/structs"
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

func (j *ResourceJSONData) GetFile() string {
	pwd := files.GetPwd()

	return path.Join(pwd, "infrastrure/terraform/modules/ovh", "resource_auto.tf.json")
}

func (r ResourceJSONData) InitResourcesFiles() {
	// Create file if not exist
	files.CreateIfNotExistsWithContent(r.GetFile(), "{}")
}

func (r *ResourceJSONData) ParseResourcesJSON() error {
	err := json.Unmarshal(files.ReadFile(r.GetFile()), &r)
	if err != nil {
		return err
	}
	return nil
}

func (r *ResourceJSONData) WriteResources() {
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
		files.WriteInFile(r.GetFile(), "{}")
		return
	}
	files.WriteInFile(r.GetFile(), string(data))
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
