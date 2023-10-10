package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"text/tabwriter"
)

type ClusterRole struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   ClusterMeta  `json:"metadata"`
	Rules      []ClusterRule `json:"rules"`
}

type ClusterMeta struct {
	Annotations        map[string]string `json:"annotations"`
	CreationTimestamp  string            `json:"creationTimestamp"`
	Labels             map[string]string `json:"labels"`
	Name               string            `json:"name"`
	ResourceVersion    string            `json:"resourceVersion"`
	UID                string            `json:"uid"`
}

type ClusterRule struct {
	APIGroups     []string `json:"apiGroups"`
	ResourceNames []string `json:"resourceNames,omitempty"`
	Resources     []string `json:"resources"`
	Verbs         []string `json:"verbs"`
}

func main() {
	out, err := exec.Command("kubectl", "get", "clusterroles", "-o", "json").Output()
	if err != nil {
		panic(err)
	}

var clusterRolesList struct {
    Items []ClusterRole `json:"items"`
}

if err := json.Unmarshal(out, &clusterRolesList); err != nil {
    panic(err)
}

w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
fmt.Fprintln(w, "User\tKind\tRole(Name)\tapiGroups\tResources\tVerbs")
fmt.Fprintln(w, "----\t----\t----------\t---------\t---------\t-----")

for _, role := range clusterRolesList.Items {
    var firstAPIGroup, firstResource, firstVerb bool = true, true, true

    for _, rule := range role.Rules {
        for _, apiGroup := range rule.APIGroups {
            if firstAPIGroup {
                fmt.Fprintf(w, "%s\t%s\t%s\t%s\t", "UserHere", role.Kind, role.Metadata.Name, apiGroup)
                firstAPIGroup = false
            } else {
                fmt.Fprintf(w, "\t\t\t%s\t", apiGroup)
            }

            for _, resource := range rule.Resources {
                if firstResource {
                    fmt.Fprintf(w, "%s\t", resource)
                    firstResource = false
                } else {
                    fmt.Fprintf(w, "\t\t\t\t%s\t", resource)
                }

                for _, verb := range rule.Verbs {
                    if firstVerb {
                        fmt.Fprintf(w, "%s\t\n", verb)
                        firstVerb = false
                    } else {
                        fmt.Fprintf(w, "\t\t\t\t\t%s\t\n", verb)
                    }
                }
                firstVerb = true
            }
            firstResource = true
        }
        firstAPIGroup = true
    }

    fmt.Fprintln(w, "----\t----\t----------\t---------\t---------\t-----")
}
w.Flush()

}

