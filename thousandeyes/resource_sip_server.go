package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceSIPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.SIPServer{}, schemas),
		Create: resourceSIPServerCreate,
		Read:   resourceSIPServerRead,
		Update: resourceSIPServerUpdate,
		Delete: resourceSIPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceSIPServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetSIPServer(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceSIPServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.SIPServer{}).(*thousandeyes.SIPServer)
	// While most ThousandEyes updates only require updated fields and specifically
	// disallow some fields on update, SIP Server tests actually require a few fields
	// within the targetSipCredentials object to be retained on update.
	// Calls without port, protocol, or sipRegistrar will fail, whereas sipProxy
	// being absent will cause the update to remove the  value.
	// Unlike other cases, we can send all non-updated values within targetSipCredentials
	// without being rejected.
	fullUpdate := buildSIPServerStruct(d)
	update.TargetSIPCredentials = fullUpdate.TargetSIPCredentials
	_, err := client.UpdateSIPServer(id, *update)
	if err != nil {
		return err
	}
	return resourceSIPServerRead(d, m)
}

func resourceSIPServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteSIPServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceSIPServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildSIPServerStruct(d)
	remote, err := client.CreateSIPServer(*local)
	if err != nil {
		return err
	}
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceSIPServerRead(d, m)
}

func buildSIPServerStruct(d *schema.ResourceData) *thousandeyes.SIPServer {
	return ResourceBuildStruct(d, &thousandeyes.SIPServer{}).(*thousandeyes.SIPServer)
}
