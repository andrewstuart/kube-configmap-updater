package main

import (
	"log"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/pkg/api/v1"
	"k8s.io/client-go/1.4/pkg/labels"
	"k8s.io/client-go/1.4/pkg/watch"
	"k8s.io/client-go/1.4/rest"
)

const (
	updateKey = "astuart.co/configMapBehavior"
)

func main() {
	sel, err := labels.Parse(updateKey)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("Cluster config", err)
	}

	cli, err := kubernetes.NewForConfig(cfg)

	for {
		w, err := cli.ConfigMaps("").Watch(api.ListOptions{})
		if err != nil {
			log.Println("Watch error", err)
		}

		for evt := range w.ResultChan() {

			et := watch.EventType(evt.Type)
			if et != watch.Added && et != watch.Modified {
				continue
			}
			switch item := evt.Object.(type) {
			case *v1.ConfigMap:
				n, ns := item.Name, item.Namespace

				pods, err := cli.Pods(ns).List(api.ListOptions{LabelSelector: sel})
				if err != nil {
					log.Println("Pod query err", err)
					continue
				}

				//Loop through all pods that match our selector for astuart.co/configMapUpdates
			podLoop:
				for _, pod := range pods.Items {
					for _, vol := range pod.Spec.Volumes {

						//If a volume is found in the spec that matches the updated volume name, delete the pod
						if vol.ConfigMap != nil && vol.ConfigMap.Name == n {
							cli.Pods(ns).Delete(pod.Name, nil)
							continue podLoop
						}
					}

				}
			default:
				log.Println("Unexpected object type: ", evt.Object.GetObjectKind())
			}
		}
	}
}
