package instance

import (
	"context"
	"io/ioutil"

	"path/filepath"

	utility "github.com/gigiozzz/depiy/common-libs/utilities"
	"github.com/gigiozzz/depiy/operators/bundle-operator/api/v1alpha1"
	"github.com/gigiozzz/depiy/operators/bundle-operator/controllers/applyer"

	common "github.com/gigiozzz/depiy/common-libs/commons"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Manifest struct {
	Base *common.BaseK8sStructure
}

func NewManifest(base *common.BaseK8sStructure) *Manifest {
	return &Manifest{
		Base: base,
	}
}

func (d *Manifest) ApplyManifest(ctx context.Context, cr *v1alpha1.EntandoBundleInstanceV2,
	scheme *runtime.Scheme,
	manifestPath string) error {
	log := d.Base.Log
	// read yaml
	yfile, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return err
	}
	var config *rest.Config
	config, err = rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			var kubeconfig string
			if home := homedir.HomeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			}

			var internalError error
			config, internalError = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if internalError != nil {
				return err
			}
			log.Info("Use kube config")
		}
	} else {
		log.Info("Use incluster config")
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	// apply yaml
	ns, _ := utility.GetWatchNamespace()
	applyOptions := applyer.NewApplyOptions(dynamicClient, discoveryClient)
	if err := applyOptions.Apply(context.TODO(), ns, []byte(yfile)); err != nil {
		return err
	}

	return nil
}
