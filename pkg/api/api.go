package api

import (
	"bytes"
	"io"

	"github.com/alecjacobs5401/kube-client-wrapper/pkg/types"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	// utilities for kubernetes integration
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

// Client is a wrapper around a Kubernetes Interface
type Client struct {
	client    k8s.Interface
	namespace string
}

// NewClient takes in a ClientConfig and generates a new Client to interface with the
// Kubernetes Cluster as outlined in the ClientConfig
func NewClient(config types.ClientConfig) (*Client, error) {
	c := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: config.ConfigFile},
		&clientcmd.ConfigOverrides{CurrentContext: config.Context, Context: clientcmdapi.Context{Namespace: config.Namespace}},
	)
	cc, err := c.ClientConfig()
	if err != nil {
		return nil, errors.Wrapf(err, "building config from kube config located at %s", config.ConfigFile)
	}

	namespace := ""
	if !config.AllNamespaces {
		n, _, err := c.Namespace()
		if err != nil {
			return nil, errors.Wrap(err, "getting namespace for client")
		}
		namespace = n
	}

	client, err := k8s.NewForConfig(cc)
	if err != nil {
		return nil, errors.Wrap(err, "building kubernetes client")
	}

	return &Client{client: client, namespace: namespace}, nil
}

// Pods takes a PodSelectors and returns an array of Kubernetes Pods based on those selectors
func (c *Client) Pods(selectors types.PodSelectors) ([]corev1.Pod, error) {
	podsAPI := c.client.CoreV1().Pods(c.namespace)

	pods := []corev1.Pod{}
	if len(selectors.Names) == 0 {
		podList, err := podsAPI.List(metav1.ListOptions{LabelSelector: selectors.Label, FieldSelector: selectors.Field})
		if err != nil {
			return nil, errors.Wrapf(err, "listing pods with LabelSelector: %s and FieldSelector: %s", selectors.Label, selectors.Field)
		}
		pods = podList.Items
	} else {
		for _, name := range selectors.Names {
			pod, err := podsAPI.Get(name, metav1.GetOptions{})
			if err != nil && !k8sErrors.IsNotFound(err) {
				return nil, errors.Wrapf(err, "Pod %s failed to delete!", name)
			}

			if err == nil || !k8sErrors.IsNotFound(err) {
				pods = append(pods, *pod)
			}
		}
	}

	return pods, nil
}

// Events takes a EventSelectors and returns an array of Kubernetes Events based on those selectors
func (c *Client) Events(selectors types.EventSelectors) ([]corev1.Event, error) {
	eventList, err := c.client.CoreV1().Events(c.namespace).List(metav1.ListOptions{LabelSelector: selectors.Label, FieldSelector: selectors.Field})
	if err != nil {
		return nil, errors.Wrapf(err, "listing events with LabelSelector: %s and FieldSelector: %s", selectors.Label, selectors.Field)
	}

	return eventList.Items, nil
}

// PodLogs grabs the logs for a specific Pod Container. If container is empty string, the default Pod
// Container will be used.
func (c *Client) PodLogs(pod corev1.Pod, container string) (string, error) {
	req := c.client.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{})
	podLogs, err := req.Stream()
	if err != nil {
		return "", errors.Wrap(err, "streaming log results")
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", errors.Wrap(err, "copying streamed log contents to buffer")
	}

	return buf.String(), nil
}
