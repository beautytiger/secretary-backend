package k8s

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"sync"
)

var clientset *kubernetes.Clientset

func GetClient() *kubernetes.Clientset {
	var once = sync.Once{}
	once.Do(func() {
		var err error
		var kubeconfig = filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err := rest.InClusterConfig()
		if err != nil {
			fmt.Printf("not in k8s cluster: %s\n", err.Error())
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("can not find kuberconfig: %s\n", err.Error())
			panic(err.Error())
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	})
	return clientset
}
