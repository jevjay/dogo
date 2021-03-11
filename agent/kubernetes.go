// Copyright 2021 tappythumbz development
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package agent

import (
	"bytes"
	"context"
	"flag"
	"io"
	"path/filepath"
	"strings"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type KubernetesAgentConfig struct {
	Name       string
	Namespace  string
	Image      string
	Kubeconfig string
	Command    string
}

type KubernetesAgent struct {
	Pod core.Pod
	Out string
}

func NewInternalClient() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func NewExternalClient(kubeconfig *string) (*kubernetes.Clientset, error) {

	// check if machine has home directory.
	if home := homedir.HomeDir(); home != "" {
		// read kubeconfig flag. if not provided use config file $HOME/.kube/config
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// build configuration from the config file.
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	// create kubernetes clientset.
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func createPodOjbect(config KubernetesAgentConfig) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      config.Name,
			Namespace: config.Namespace,
			Labels: map[string]string{
				"dogo-agent": "true",
				"app":        config.Name,
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            config.Name,
					Image:           config.Image,
					ImagePullPolicy: core.PullIfNotPresent,
					Command:         strings.Split(config.Command, " "),
				},
			},
		},
	}
}

func getPodLogs(pod core.Pod, client *kubernetes.Clientset, ctx context.Context) (string, error) {
	podLogOpts := core.PodLogOptions{}
	req := client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	str := buf.String()

	return str, nil
}

func ExecuteKubernetes(config KubernetesAgentConfig) (*KubernetesAgent, error) {
	ctx := context.Background()
	var a KubernetesAgent
	if len(config.Kubeconfig) > 0 {
		client, err := NewExternalClient(&config.Kubeconfig)
		if err != nil {
			return nil, err
		}
		// build the pod defination we want to deploy
		pod := createPodOjbect(config)

		// now create the pod in kubernetes cluster using the clientset
		pod, err = client.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		out, err := getPodLogs(*pod, client, ctx)
		if err != nil {
			return nil, err
		}

		a.Pod = *pod
		a.Out = out
	} else {
		client, err := NewInternalClient()
		if err != nil {
			return nil, err
		}

		// build the pod defination we want to deploy
		pod := createPodOjbect(config)

		// now create the pod in kubernetes cluster using the clientset
		pod, err = client.CoreV1().Pods(pod.Namespace).Create(ctx, pod, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}

		out, err := getPodLogs(*pod, client, ctx)
		if err != nil {
			return nil, err
		}

		a.Pod = *pod
		a.Out = out
	}

	return &a, nil
}
