package k8s

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type PodCPU struct {
	Name string
	CPU  string
}

type PodMemory struct {
	Name   string
	Memory string
}

func GetTopPodsByCPU(ctx context.Context, namespace string, count int) ([]PodCPU, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var pods []PodCPU
	for _, m := range podMetrics.Items {
		totalCPU := int64(0)
		for _, c := range m.Containers {
			cpuQuantity := c.Usage.Cpu().MilliValue()
			totalCPU += cpuQuantity
		}
		pods = append(pods, PodCPU{Name: m.Name, CPU: fmt.Sprintf("%dm", totalCPU)})
	}

	sort.Slice(pods, func(i, j int) bool {
		cpuI, _ := strconv.Atoi(strings.TrimSuffix(pods[i].CPU, "m"))
		cpuJ, _ := strconv.Atoi(strings.TrimSuffix(pods[j].CPU, "m"))
		return cpuI > cpuJ
	})

	if count > len(pods) {
		count = len(pods)
	}
	return pods[:count], nil
}

func GetTopPodsByMemory(ctx context.Context, namespace string, count int) ([]PodMemory, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics client: %w", err)
	}

	podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var pods []PodMemory
	for _, m := range podMetrics.Items {
		totalMem := int64(0)
		for _, c := range m.Containers {
			memQuantity := c.Usage.Memory().Value()
			totalMem += memQuantity
		}
		pods = append(pods, PodMemory{Name: m.Name, Memory: fmt.Sprintf("%dMi", totalMem/(1024*1024))})
	}

	sort.Slice(pods, func(i, j int) bool {
		memI, _ := strconv.Atoi(strings.TrimSuffix(pods[i].Memory, "Mi"))
		memJ, _ := strconv.Atoi(strings.TrimSuffix(pods[j].Memory, "Mi"))
		return memI > memJ
	})

	if count > len(pods) {
		count = len(pods)
	}
	return pods[:count], nil
}
