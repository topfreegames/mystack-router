package extensions

import (
	"os"
	"os/exec"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/pkg/labels"
	"k8s.io/client-go/pkg/util/flowcontrol"
	"k8s.io/client-go/rest"

	"github.com/topfreegames/mystack-router/model"
	"github.com/topfreegames/mystack-router/nginx"
)

const (
	nginxConfigDir      = "/etc/nginx"
	nginxConfigFilePath = nginxConfigDir + "/nginx.conf"
)

// Watcher is the extension that watches for kubernetes services changes
type Watcher struct {
	config        *viper.Viper
	tokenPerSec   float32 // Token per second on Token-Bucket algorithm
	burst         int     // Bucket size on Token-Bucket algorithm
	kubeConfig    *rest.Config
	kubeClientSet *kubernetes.Clientset
}

// NewWatcher creates a new watcher instance
func NewWatcher(config *viper.Viper) (*Watcher, error) {
	w := &Watcher{
		config: config,
	}
	w.loadConfigurationDefaults()
	err := w.configure()
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Watcher) loadConfigurationDefaults() {
}

func (w *Watcher) configure() error {
	w.tokenPerSec = float32(w.config.GetFloat64("watcher.token-per-sec"))
	w.burst = w.config.GetInt("watcher.burst")

	var err error
	w.kubeConfig, err = rest.InClusterConfig()
	if err != nil {
		return err
	}

	w.kubeClientSet, err = kubernetes.NewForConfig(w.kubeConfig)
	return err
}

func (w *Watcher) getMyStackServices() (*v1.ServiceList, error) {
	labelMap := labels.Set{"router.mystack/routable": "true"}
	listOptions := v1.ListOptions{
		LabelSelector: labelMap.AsSelector().String(),
		FieldSelector: fields.Everything().String(),
	}
	services, err := w.kubeClientSet.CoreV1().Services(api.NamespaceAll).List(listOptions)
	return services, err
}

func (w *Watcher) build() (*model.RouterConfig, error) {
	appServices, err := w.getMyStackServices()
	if err != nil {
		return nil, err
	}

	routerConfig := model.NewRouterConfig()

	for _, appService := range appServices.Items {
		appConfig, err := model.BuildAppConfig(w.kubeClientSet, appService, routerConfig)
		if err != nil {
			return nil, err
		}
		if appConfig != nil {
			routerConfig.AppConfigs = append(routerConfig.AppConfigs, appConfig)
		}
	}

	return routerConfig, nil
}

// Start starts the watcher, this call is blocking!
func (w *Watcher) Start() error {
	l := log.WithFields(log.Fields{
		"tokenPerSecond": w.tokenPerSec,
		"burst":          w.burst,
	})
	l.Info("starting mystack watcher")
	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(w.tokenPerSec, w.burst)
	known := &model.RouterConfig{}

	err := os.MkdirAll(nginxConfigDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = exec.Command("touch", nginxConfigFilePath).Run()
	if err != nil {
		return err
	}

	err = nginx.Start(l)
	if err != nil {
		return err
	}

	for {
		rateLimiter.Accept()
		routerConfig, err := w.build()
		if err != nil {
			return err
		}
		// Generate new RouterConfig with Build calling getMyStackServices
		// If DeepEquals to known, call continue to loop
		// else, calls reload and save new known
		if reflect.DeepEqual(routerConfig, known) {
			continue
		}
		err = nginx.WriteConfig(routerConfig, nginxConfigFilePath)
		if err != nil {
			log.Printf("Failed to write new nginx configuration; continuing with existing configuration: %v", err)
			continue
		}
		err = nginx.Reload(l)
		if err != nil {
			return err
		}
		known = routerConfig
	}
}
