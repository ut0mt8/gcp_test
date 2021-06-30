package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/compute/v1"
	"os"
)

// wesh commit

func main() {

	project := flag.String("project", "", "GCP project ID (Required)")
	flag.Parse()

	if *project == "" {
		flag.PrintDefaults()
		os.Exit(1)

	}

	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	aggrList, err := computeService.Instances.AggregatedList(*project).Do()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for zone, values := range aggrList.Items {
		if values.Instances != nil {
			for _, instance := range values.Instances {
				ip := "none"
				lifecycle := "undef"
				for _, nif := range instance.NetworkInterfaces {
					ip = nif.NetworkIP
				}
				for _, meta := range instance.Metadata.Items {
					if meta.Key == "lifecycle" {
						lifecycle = *meta.Value
					}
				}
				fmt.Printf("%s\t%s\t%s\t%s\t%s\n", zone, instance.Name, ip, instance.Status, lifecycle)
			}
		}
	}
}
