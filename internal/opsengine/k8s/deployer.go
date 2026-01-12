package k8s

import (
	"context"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Deployer struct {
	clientset *kubernetes.Clientset
}

// NewDeployer 初始化K8s部署器
func NewDeployer() (*Deployer, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("kubeconfig not found")
	}
	// 加载配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// 创建Clientset (K8s 操作入口)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Deployer{clientset: clientset}, nil
}

func (d *Deployer) UpdateImage(ctx context.Context, namespace, deploymentName, newImage string) error {
	fmt.Printf("☸️  [K8s] 正在更新 Deployment [%s] -> %s\n", deploymentName, newImage)

	deploymentsClient := d.clientset.AppsV1().Deployments(namespace)

	// 1. 获取当前的 Deployment 对象 (Get)
	deployment, err := deploymentsClient.Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get deployment: %v", err)
	}

	// 2. 修改内存中的对象 (Modify)
	containers := deployment.Spec.Template.Spec.Containers
	if len(containers) == 0 {
		return fmt.Errorf("deployment has no containers")
	}

	oldImage := containers[0].Image
	if oldImage == newImage {
		fmt.Printf("⚠️  [K8s] 镜像未发生变化，跳过更新\n")
		return nil
	}
	containers[0].Image = newImage

	// 3. 提交更新到集群 (Update)
	_, err = deploymentsClient.Update(ctx, deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update deployment: %v", err)
	}
	fmt.Printf("✅ [K8s] 更新成功！从 %s 变更为 %s\n", oldImage, newImage)
	return nil
}
