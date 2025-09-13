package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsNotFoundErr(err error) bool {
	return apierrors.IsNotFound(err)
}

func IsNamespaceNotFoundErr(err error) bool {
	return apierrors.IsNotFound(err)
}

func GetPodLogs(ctx context.Context, podName, namespace string) ([]byte, error) {
	clientset, err := GetClientSet()
	if err != nil {
		return nil, fmt.Errorf("failed to build kube client: %w", err)
	}

	_, err = clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})
	stream, err := req.Stream(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open log stream: %w", err)
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stream)
	if err != nil {
		return nil, fmt.Errorf("failed to read log stream: %w", err)
	}

	return buf.Bytes(), nil
}
