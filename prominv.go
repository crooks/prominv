package main

import (
	"context"
	"fmt"
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

const (
	PrometheusURL = "http://plosysmon02.westernpower.co.uk:9090"
	PromQL        = "up{job=\"node_exporter\"}"
)

var (
	flags *config.Flags
)

type InventoryItem struct {
	App string `json:"app"`
	Env string `json:"env"`
}

func kvExtract(result string) map[string]string {
	openCurly := strings.SplitN(result, "{", 2)[1]
	closeCurly := strings.SplitN(openCurly, "}", 2)[0]
	labelPairs := strings.Split(closeCurly, ",")
	labels := make(map[string]string)
	for _, l := range labelPairs {
		kv := strings.SplitN(strings.Trim(l, " "), "=", 2)
		labels[kv[0]] = strings.Trim(kv[1], "\"")
	}
	return labels
}

// newAPI returns a new instance of the Prometheus v1 API
func newAPI() v1.API {
	client, err := api.NewClient(api.Config{
		Address: PrometheusURL,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	return v1.NewAPI(client)
}

func runPromQL(query string) []model.LabelSet {
	v1api := newAPI()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	lbls, warnings, err := v1api.Series(ctx, []string{query}, time.Now().Add(-time.Hour), time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	return lbls
}

func makeInventory() {
	lbls := runPromQL(PromQL)
	var err error
	inventory := "{}"
	inventory, err = sjson.Set(inventory, "_meta", "hostvars")
	if err != nil {
		log.Fatal(err)
	}
	inventory, err = sjson.Set(inventory, "all", "children:[]")
	if err != nil {
		log.Fatal(err)
	}
	inventory, err = sjson.Set(inventory, "all.children.-1", "prometheus")
	if err != nil {
		log.Fatal(err)
	}
	for _, lbl := range lbls {
		labels := kvExtract(lbl.String())
		instance, ok := labels["instance"]
		if !ok {
			//log.Fatalf("Instance:%s has no instance label", labels["__name__"])
			continue
		}
		instanceEscaped := strings.Replace(instance, ".", "\\.", -1)
		inventory, err = sjson.Set(inventory, "all.prometheus.-1", instanceEscaped)
		if err != nil {
			log.Fatal(err)
		}
		// Delete labels that we don't want to appear within the hostvars map
		delete(labels, "instance")
		delete(labels, "__name__")
		hostvarsKey := fmt.Sprintf("_meta.hostvars.%s", instanceEscaped)
		inventory, err = sjson.Set(inventory, hostvarsKey, labels)
	}
	fmt.Println(inventory)
}

func main() {
	flags = config.ParseFlags()
	if flags.List {
		makeInventory()
	}
}
