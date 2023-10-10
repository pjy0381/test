package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

type ClusterRole struct {
	Name      string
	User      string
	Kind      string
	APIGroups []string
	Resources []string
	Verbs     []string
}

func main() {
	roles := []ClusterRole{
		{
			Name:      "view",
			User:      "kube",
			Kind:      "ClusterRole",
			APIGroups: []string{"CORE", "apps"},
			Resources: []string{"bindings", "configmaps", "endpoints", "events"},
			Verbs:     []string{"get", "list", "watch"},
		},
		{
			Name:      "edit",
			User:      "sphere",
			Kind:      "ClusterRole",
			APIGroups: []string{"CORE"},
			Resources: []string{"bindings", "replicationcontrollers/scale", "replicationcontrollers/status"},
			Verbs:     []string{"get", "list", "watch"},
		},
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Index\tUser\tKind\tRole(Name)\tapiGroups\tResources\tVerbs")
	fmt.Fprintln(w, "-----\t----\t----\t----------\t---------\t---------\t-----")

	for idx, role := range roles {
		for i, apiGroup := range role.APIGroups {
			fmt.Fprintf(w, "%04d\t%s\t%s\t%s\t%s\t%s\t[%s]\n",
				idx+1,
				role.User,
				role.Kind,
				role.Name,
				apiGroup,
				role.Resources[i],
				role.Verbs[i],
			)
		}
		fmt.Fprintln(w, "-----\t----\t----\t----------\t---------\t---------\t-----")
	}
	w.Flush()
}
