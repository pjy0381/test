package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/tabwriter"
)

type ClusterRole struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		Annotations      map[string]string `json:"annotations,omitempty"`
		CreationTimestamp string            `json:"creationTimestamp"`
		Labels            map[string]string `json:"labels,omitempty"`
		Name              string            `json:"name"`
		ResourceVersion   string            `json:"resourceVersion"`
		UID               string            `json:"uid"`
	} `json:"metadata"`
	Rules []struct {
		APIGroups []string `json:"apiGroups"`
		Resources []string `json:"resources"`
		Verbs     []string `json:"verbs"`
	} `json:"rules"`
	AggregationRule struct {
		ClusterRoleSelectors []struct {
			MatchLabels map[string]string `json:"matchLabels"`
		} `json:"clusterRoleSelectors"`
	} `json:"aggregationRule,omitempty"`
}

type ClusterRolesList struct {
	Items []ClusterRole `json:"items"`
}

func main() {
	out, err := exec.Command("kubectl", "get", "clusterrole", "-o", "json").Output()
	if err != nil {
		panic(err)
	}

	var clusterRoles ClusterRolesList
	err = json.Unmarshal(out, &clusterRoles)
	if err != nil {
		panic(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, '\t', 0)
	fmt.Fprintln(w, "Name\tAPIGroups\tResources\tVerbs")

	for _, role := range clusterRoles.Items {
		for _, rule := range role.Rules {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				role.Metadata.Name,
				strings.Join(rule.APIGroups, ", "),
				strings.Join(rule.Resources, ", "),
				strings.Join(rule.Verbs, ", "),
			)
		}
	}
	w.Flush()
}
