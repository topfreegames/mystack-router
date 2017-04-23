// mystack
// https://github.com/topfreegames/mystack-router
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2017 Top Free Games <backend@tfgco.com>

package extensions

import (
	"fmt"
	"os"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/fields"
	"k8s.io/client-go/pkg/labels"
	"k8s.io/client-go/pkg/util/flowcontrol"
	"k8s.io/client-go/rest"

	_ "github.com/lib/pq"
	"github.com/topfreegames/mystack-router/models"
	"github.com/topfreegames/mystack-router/nginx"
)

const (
	nginxConfigDir      = "/etc/nginx"
	nginxConfigFilePath = nginxConfigDir + "/nginx.conf"
)

// Watcher is the extension that watches for kubernetes services changes
type Watcher struct {
	tokenPerSec      float32 // Token per second on Token-Bucket algorithm
	burst            int     // Bucket size on Token-Bucket algorithm
	KubeClientSet    kubernetes.Interface
	KubeDomainSuffix string
	DB               models.DB
}

//NewWatcher creates a new watcher with chosen clientset
//If clientset is nil, creates a inCluster clientset
func NewWatcher(config *viper.Viper, clientset kubernetes.Interface) (*Watcher, error) {
	w := &Watcher{}
	w.configureProps(config)
	err := w.configureDB(config)
	if err != nil {
		return nil, err
	}

	if clientset == nil {
		err := w.configureClient()
		return w, err
	}

	w.KubeClientSet = clientset
	return w, nil
}

func (w *Watcher) configureDB(config *viper.Viper) error {
	db, err := getDB(config)
	if err != nil {
		return err
	}

	w.DB = db
	return nil
}

func getDB(config *viper.Viper) (*sqlx.DB, error) {
	host := config.GetString("postgres.host")
	user := config.GetString("postgres.user")
	dbName := config.GetString("postgres.dbname")
	password := config.GetString("postgres.password")
	port := config.GetInt("postgres.port")
	connectionTimeoutMS := config.GetInt("postgres.connectionTimeoutMS")

	db, err := models.GetDB(
		host, user, dbName, password,
		port, connectionTimeoutMS,
	)

	fmt.Printf("%#v, %#v", db, err)
	fmt.Printf("Host %s, port %d", host, port)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func (w *Watcher) configureProps(config *viper.Viper) {
	key := "watcher.router-refresh-min-interval-s"
	w.burst = 1
	w.tokenPerSec = float32(w.burst) / float32(config.GetFloat64(key))
	w.KubeDomainSuffix = config.GetString("kubernetes.service-domain-suffix")
}

func (w *Watcher) configureClient() error {
	kubeConfig, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	w.KubeClientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	return err
}

//GetMyStackServices return list of services running on k8s
func (w *Watcher) GetMyStackServices() (*v1.ServiceList, error) {
	labelMap := labels.Set{"mystack/routable": "true"}
	listOptions := v1.ListOptions{
		LabelSelector: labelMap.AsSelector().String(),
		FieldSelector: fields.Everything().String(),
	}
	services, err := w.KubeClientSet.CoreV1().Services(api.NamespaceAll).List(listOptions)
	return services, err
}

//Build construct the routerConfig of cluster
func (w *Watcher) Build() (*models.RouterConfig, error) {
	appServices, err := w.GetMyStackServices()
	if err != nil {
		return nil, err
	}

	routerConfig := models.NewRouterConfig(w.KubeDomainSuffix)
	for _, appService := range appServices.Items {
		appConfig := models.BuildAppConfig(
			&appService,
			w.KubeDomainSuffix,
		)
		routerConfig.AppConfigs = append(routerConfig.AppConfigs, appConfig)
	}

	return routerConfig, nil
}

//CreateConfigFile make nginx directory (if not exists) and create nginx config file.
func (w *Watcher) CreateConfigFile() error {
	err := os.MkdirAll(nginxConfigDir, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(nginxConfigFilePath)
	return err
}

// Start starts the watcher, this call is blocking!
func (w *Watcher) Start(fs models.FileSystem) error {
	l := log.WithFields(log.Fields{
		"tokenPerSecond": w.tokenPerSec,
		"burst":          w.burst,
	})
	l.Info("starting mystack watcher")
	rateLimiter := flowcontrol.NewTokenBucketRateLimiter(w.tokenPerSec, w.burst)
	known := &models.RouterConfig{}

	// Remove this dir because it overwrites our conf if exists
	err := os.RemoveAll(fmt.Sprintf("%s/conf.d", nginxConfigDir))
	if err != nil {
		return err
	}

	for {
		rateLimiter.Accept()
		routerConfig, err := w.Build()
		if err != nil {
			return err
		}
		// Generate new RouterConfig with Build calling getMyStackServices
		// If DeepEquals to known, call continue to loop
		// else, calls reload and save new known
		if reflect.DeepEqual(routerConfig, known) {
			continue
		}
		l.Info("new config found")
		err = routerConfig.AddCustomDomains(w.DB)
		if err != nil {
			log.Printf("Failed to get custom domains: %v", err)
			continue
		}
		err = nginx.WriteConfig(routerConfig, fs, nginxConfigFilePath)
		if err != nil {
			log.Printf("Failed to write new nginx configuration; continuing with existing configuration: %v", err)
			continue
		}
		err = nginx.Reload(l)
		if err != nil {
			return err
		}
		l.Info("nginx reloaded")
		err = nginx.WriteConfig(routerConfig, fs, nginxConfigFilePath)
		known = routerConfig
	}
}
