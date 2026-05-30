package main

import (
	"context"
	"fmt"
	"gitlab/prominv/children"
	"gitlab/prominv/config"
	"log"
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
	// Iterate over the returned metrics
	for _, result := range results {
		labels := result.Metric
		_, ok := labels["instance"]
		if !ok {
			//log.Fatalf("Instance:%s has no instance label", labels["__name__"])
			continue
		}
		instance := string(labels["instance"])
		// All instances get added to the prometheus child group
		children.AddMember("prometheus", instance)
		if err != nil {
			log.Fatalf("Failed to add %s to the \"prometheus\" group: %v", instance, err)
		}
		// If there's an "groupBy" label, populate an inventory group for it
		if groupBy, ok := labels[model.LabelName(cfg.Labels.GroupBy)]; ok {
			// Make all group names lowercase
			groupName := strings.ToLower(string(groupBy))
			children.AddMember(groupName, instance)
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

	allChildren := children.GetAllChildren()
	inventory, err = sjson.Set(inventory, "all.children", allChildren)
	// Iterate over the child groups and add the members to them
	for _, childName := range allChildren {
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
