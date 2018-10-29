package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ankyra/escape-inventory/dao"
	"github.com/ankyra/escape-inventory/dao/types"
)

func CallWebHook(namespace, unit, version, releaseId, username, url string) {
	if url == "" {
		return
	}
	prj := types.NewProject(namespace)
	app := types.NewApplication(namespace, unit)
	prjHooks, err := dao.GetNamespaceHooks(prj)
	if err != nil {
		log.Println("ERROR: Failed to get Inventory Project Hooks from database:", err)
		return
	}
	unitHooks, err := dao.GetApplicationHooks(app)
	if err != nil {
		log.Println("ERROR: Failed to get Inventory Application Hooks from database:", err)
		return
	}
	downstreamHooks, err := dao.GetDownstreamHooks(app)
	if err != nil {
		log.Println("ERROR: Failed to get Upstream Hooks from database:", err)
		return
	}
	fmt.Println("Downstream hooks:")
	fmt.Println(downstreamHooks)
	for _, hooks := range downstreamHooks {
		fmt.Println(hooks)
	}
	data := map[string]interface{}{
		"event":            "NEW_UPLOAD",
		"project":          namespace,
		"project_hooks":    prjHooks,
		"unit":             unit,
		"unit_hooks":       unitHooks,
		"downstream_hooks": downstreamHooks,
		"version":          version,
		"release":          namespace + "/" + releaseId,
		"username":         username,
	}
	body, err := json.Marshal(data)
	if err != nil {
		log.Println("ERROR: Failed to marshal webhook request:", err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Println("ERROR: Failed to create webhook request:", err)
		return
	}
	log.Println("INFO: Calling webhook:", url)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: Failed to call webhook:", err)
		return
	}
	if resp.StatusCode != 200 {
		log.Println("ERROR: Failed to call webhook, expecting status code 200, but got", resp.StatusCode)
		return
	}
}
