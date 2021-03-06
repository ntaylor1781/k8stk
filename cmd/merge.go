/*
This Source Code Form is subject to the terms of the Mozilla Public
License, v. 2.0. If a copy of the MPL was not distributed with this
file, You can obtain one at https://mozilla.org/MPL/2.0/.

*/
package cmd

import (
    "fmt"
    "strconv"
    "github.com/mjnt/k8stk/util"
    
    
    "github.com/spf13/cobra"
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
    Use:   "merge [flags] KUBECONFIG KUBECONFIG [KUBECONFIG...]",
    Short: "Merge multiple kubeconfig files",
    Long: `Merges any number of kubeconfig files passed, utilizing the first one as the base.`,
    Args: cobra.MinimumNArgs(2),
    Example: "k8stk merge -o newConfig ~/.kube/config ~/.kube/new_cluster_config",
    Run: func(cmd *cobra.Command, args []string) {
        mergeKubeconfig(cmd, args)
    },
}

func init() {
    rootCmd.AddCommand(mergeCmd)
    
    // Here you will define your flags and configuration settings.
    
    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // mergeCmd.PersistentFlags().String("foo", "", "A help for foo")
    
    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    mergeCmd.Flags().StringP("output", "o", "", "File to output the merged config to. If not passed the output will be sent to stdout")
}

func mergeKubeconfig(cmd *cobra.Command, args []string) {
    base := util.ParseYaml(args[0])

    cluster_nm := 0
    context_nm := 0
    user_nm := 0

    for _, file := range args[1:] {
        tmp := util.ParseYaml(string(file))
        namesChanged := make(map[string]string)

        for _, cluster := range tmp.Clusters {
            for _, i := range base.Clusters {
                if i.Name == cluster.Name {
                    fmt.Printf("Found a duplicate cluster name. Renaming %s from file %s.\n", i.Name, file)
                    cluster_nm ++
                    var updated_name string = cluster.Name+strconv.Itoa(cluster_nm)
                    namesChanged[cluster.Name] = updated_name
                    cluster.Name = updated_name
                    context_nm ++
                    user_nm ++
                }
            }
            base.Clusters = append(base.Clusters, cluster)
        }
        for _, user := range tmp.Users {
            for _, i := range base.Users {
                if i.Name == user.Name {
                    fmt.Printf("Found a duplicate user name, Renaming %s from file %s.\n", i.Name, file)
                    var updated_user_name string = user.Name+strconv.Itoa(user_nm)
                    namesChanged[user.Name] = updated_user_name
                    user.Name = updated_user_name
                }
            }
            base.Users = append(base.Users, user)
        }
        for _, context := range tmp.Contexts {
            for _, i := range base.Contexts {
                if i.Name == context.Name {
                    fmt.Printf("Found a duplicate context name, Renaming %s from file %s.\n", i.Name, file)
                    var updated_context_name string = context.Name+strconv.Itoa(context_nm)
                    context.Name = updated_context_name
                }
            }

            // Ensure user and cluster are updated in the context if they were changed
            if val, ok := namesChanged[context.Context.Cluster]; ok {
                context.Context.Cluster = val
            }
            if val, ok := namesChanged[context.Context.User]; ok {
                context.Context.User = val
            }
            base.Contexts = append(base.Contexts, context)
        }
    }

    util.OutputYaml(base, cmd.Flags().Lookup("output").Value.String())
}
