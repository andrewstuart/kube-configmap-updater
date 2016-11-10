// Copyright 2016 Andrew Stuart

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/watch"
	"k8s.io/client-go/rest"
)

const (
	updateKey = "astuart.co/configMapBehavior"
)

func main() {
	cfg, err := rest.InClusterConfig()

	if err != nil {
		log.Fatal("Cluster config", err)
	}

	cli, err := kubernetes.NewForConfig(cfg)

	for {
		w, err := cli.ConfigMaps("").Watch(v1.ListOptions{})
		if err != nil {
			log.Println("Watch error", err)
		}

		for evt := range w.ResultChan() {

			log.Println("Watch event triggered: %#v", evt)

			et := watch.EventType(evt.Type)
			if et != watch.Added && et != watch.Modified {
				continue
			}
			switch item := evt.Object.(type) {
			case *v1.ConfigMap:
				n, ns := item.Name, item.Namespace

				log.Printf("Configmap %s/%s updated\n", n, ns)

				pods, err := cli.Pods(ns).List(v1.ListOptions{LabelSelector: updateKey})
				if err != nil {
					log.Println("Pod query error", err)
					continue
				}

				//Loop through all pods that match our selector for astuart.co/configMapUpdates
			podLoop:
				for _, pod := range pods.Items {
					for _, vol := range pod.Spec.Volumes {

						//If a volume is found in the spec that matches the updated volume name, delete the pod
						if vol.ConfigMap != nil && vol.ConfigMap.Name == n {
							switch pod.ObjectMeta.Labels[updateKey] {
							case "Delete":
								cli.Pods(ns).Delete(pod.Name, nil)
							default:
								log.Println("Unknown behavior: ", pod.ObjectMeta.Labels[updateKey])
							}
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
