/*
Copyright 2022 The K2 Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"context"
	"fmt"

	"github.com/poneding/k2/pkg/kube"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate kubeconfig from serviceaccount",
	Long: `Generate kubeconfig from serviceaccount. For example:
	k2 config gen --from default -n default`,
	Run: func(cmd *cobra.Command, args []string) {
		genConfig(cmd, args)
	},
}

var genOptions = struct {
	kubeconfig string
	from       string
	to         string
	namespace  string
}{}

func init() {
	configCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&genOptions.kubeconfig, "kubeconfig", "", clientcmd.RecommendedHomeFile, "the path of kubeconfig")
	genCmd.Flags().StringVarP(&genOptions.from, "from", "f", "", "serviceaccount name")
	genCmd.Flags().StringVarP(&genOptions.to, "to", "t", "", "write generated config to file")
	genCmd.Flags().StringVarP(&genOptions.namespace, "namespace", "n", "default", "namespaced name")

	validateOptions()

	genCmd.MarkFlagRequired("from")
}

func validateOptions() {
	// write validations here
}

func genConfig(cmd *cobra.Command, args []string) {
	kubecfg, err := kube.Config(genOptions.kubeconfig)
	if err != nil {
		panic(err)
	}
	kubecli, err := kube.Client(genOptions.kubeconfig)
	if err != nil {
		panic(err)
	}

	sa, err := kubecli.CoreV1().ServiceAccounts(genOptions.namespace).Get(context.TODO(), genOptions.from, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	secret, err := kubecli.CoreV1().Secrets(genOptions.namespace).Get(context.TODO(), sa.Secrets[0].Name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	cfg := buildConfigFromSecret(kubecfg, secret)

	if len(genOptions.to) != 0 {
		clientcmd.WriteToFile(cfg, genOptions.to)
	} else {
		cfgContent, err := clientcmd.Write(cfg)
		if err != nil {
			panic(err)
		}
		fmt.Printf("---\n# Generated by k2 from serviceaccount.\n# - serviceaccount: %s\n# - namespace: %s\n", genOptions.from, genOptions.namespace)
		fmt.Println(string(cfgContent))
	}
}

func buildConfigFromSecret(kubecfg *rest.Config, secret *corev1.Secret) clientcmdapi.Config {
	cfgContext := "context-" + rand.String(5)
	cfgCluster := "cluster-" + rand.String(5)
	cfgUser := "user-" + rand.String(5)
	cfg := clientcmdapi.Config{
		APIVersion:     "v1",
		Kind:           "Config",
		CurrentContext: cfgContext,
		Contexts: map[string]*clientcmdapi.Context{
			cfgContext: {
				Cluster:   cfgCluster,
				AuthInfo:  cfgUser,
				Namespace: string(secret.Data["namespace"]),
			},
		},
		Clusters: map[string]*clientcmdapi.Cluster{
			cfgCluster: {
				Server:                   kubecfg.Host,
				CertificateAuthorityData: secret.Data["ca.crt"],
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			cfgUser: {
				Token: string(secret.Data["token"]),
			},
		},
	}
	return cfg
}
