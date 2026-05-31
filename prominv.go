package main

import (
	"context"
	"fmt"
	"gitlab/prominv/children"
	"gitlab/prominv/config"
	"log"
	"maps"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"github.com/tidwall/sjson"
)

var (
	cfg *config.Config
)

// newAPI returns a new instance of the Prometheus v1 API
func newAPI() v1.API {
	client, err := api.NewClient(api.Config{
		Address: cfg.URL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	return v1.NewAPI(client)
}

func runPromQL(query string) model.Vector {
	v1api := newAPI()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus at %s: %v\n", cfg.URL, err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	if result.Type() != model.ValVector {
		log.Fatal("Not a vector result")
	}
	return result.(model.Vector)
}

func makeInventory() {
	children := children.NewChildren()
	var err error
	inventory := "{}"
	inventory, err = sjson.Set(inventory, "_meta", "hostvars")
	if err != nil {
		log.Fatal(err)
	}
	inventory, err = sjson.Set(inventory, "all", "children")
	if err != nil {
		log.Fatal(err)
	}
	results := runPromQL(cfg.Query)
	// In the world of Ansible, some variable names are reserved.  This list is far from exhaustive!
	labelsReplace := map[string]string{
		"notify":      "notify_team",
		"environment": "env",
	}
	// Iterate over the returned metrics
	for _, result := range results {
		labels := result.Metric
		_, ok := labels["instance"]
		if !ok {
			//log.Fatalf("Instance:%s has no instance label", labels["__name__"])
			continue
		}
		instance := string(labels["instance"])
		// Take a clone of labels to iterate over as it's going to be modified in-situ.
		for l := range maps.Clone(labels) {
			replaceWith, ok := labelsReplace[string(l)]
			if ok {
				// A label is in the reserved list. Create a replacement and then delete the original.
				labels[model.LabelName(replaceWith)] = labels[l]
				delete(labels, l)
			}
		}
		// All instances get added to the prometheus child group
		children.AddMember("all", instance)
		if err != nil {
			log.Fatalf("Failed to add %s to the \"prometheus\" group: %v", instance, err)
		}
		// Iterate over the GroupBy labels defined in the Config
		for _, groups := range cfg.Labels.GroupBy {
			// If there's an "groupBy" label, populate an inventory group for it
			if groupBy, ok := labels[model.LabelName(groups)]; ok {
				// Hyphens are an invalid character in Ansible group names.
				groupBySanitised := strings.ReplaceAll(strings.ToLower(string(groupBy)), "-", "_")
				groupNameSanitised := strings.ReplaceAll(strings.ToLower(groups), "-", "_")
				// Make all group names lowercase and in the format <group_name>_<group_name_content>
				groupName := fmt.Sprintf("%s_%s", groupNameSanitised, groupBySanitised)
				children.AddMember(groupName, instance)
			}
		}
		// This conditional populates an "up" child group if the value of the "up" metric is 1.
		if int(result.Value) == 1 {
			children.AddMember("up", instance)
		}

		instanceEscaped := strings.Replace(instance, ".", "\\.", -1)
		// Delete labels that we don't want to appear within the hostvars map
		for _, l := range cfg.Labels.Delete {
			delete(labels, model.LabelName(l))
		}
		hostvarsKey := fmt.Sprintf("_meta.hostvars.%s", instanceEscaped)
		inventory, err = sjson.Set(inventory, hostvarsKey, labels)
	}

	// Passing "false" to GetAllChildren causes the "all" group to be excluded.  It shouldn't be a member of itself.
	inventory, err = sjson.Set(inventory, "all.children", children.GetAllChildren(false))
	if err != nil {
		log.Fatalf("Failed to create all.children: %s", err)
	}
	// Iterate over the child groups and add the members to them
	for _, childName := range children.GetAllChildren(true) {
		members, err := children.MemberSlice(childName)
		if err != nil {
			log.Fatalf("Failed to generate member slice for \"%s\": %v", childName, err)
		}
		sjKey := fmt.Sprintf("%s.hosts", childName)
		inventory, err = sjson.Set(inventory, sjKey, members)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println(inventory)
}

func main() {
	var err error
	// The config file can only be defined via a command line flag or an Environment Variable.
	cfgFilename := config.GetConfigFilename()
	cfg, err = config.ParseConfig(cfgFilename)
	if err != nil {
		log.Fatalf("Cannot parse config: %v", err)
	}
	if config.DoList() {
		makeInventory()
	} else {
		fmt.Printf("%+v\n", cfg)
	}
}
