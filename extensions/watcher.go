package extensions

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/pkg/labels"
	"k8s.io/client-go/pkg/util/flowcontrol"
	"k8s.io/client-go/rest"
)

// Watcher is the extension that watches for kubernetes services changes
type Watcher struct {
	config        *viper.Viper
	interval      time.Duration
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
	w.interval = time.Duration(w.config.GetInt("watcher.intervalms"))
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

// Start starts the watcher, this call is blocking!
func (w *Watcher) Start() {
	l := log.WithField("interval", w.interval)
	l.Info("starting mystack watcher")
	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(0.1, 1)
	for {
		rateLimiter.Accept()
		services, err := w.getMyStackServices()
		if err != nil {
			log.Error(err)
		} else {
			log.WithField("services", services.Items).Info("got items")
		}
		time.Sleep(w.interval)
	}

}
