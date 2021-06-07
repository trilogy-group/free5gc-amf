/*
 * AMF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/free5gc/amf/logger"
)

var AmfConfig Config

var cfgSyncCh chan interface{}

func InitConfigFileFactory(f string) error {
	if content, err := ioutil.ReadFile(f); err != nil {
		return err
	} else {
		AmfConfig = Config{}

		if yamlErr := yaml.Unmarshal(content, &AmfConfig); yamlErr != nil {
			return yamlErr
		}
	}
	return nil
}

func InitConfigFactory(f string) error {
	cfgSyncCh = make(chan interface{})

	CfgMgrStart()

	select {
	case <-cfgSyncCh:
		fmt.Println("XXX cfgMgr init config done")
		return nil
	case <-time.After(time.Second * 5):
		fmt.Println("XXX cfgMgr timeout read from config file")
		return InitConfigFileFactory(f)
	}

	return nil
}

func CheckConfigVersion() error {
	currentVersion := AmfConfig.GetVersion()

	if currentVersion != AMF_EXPECTED_CONFIG_VERSION {
		return fmt.Errorf("config version is [%s], but expected is [%s].",
			currentVersion, AMF_EXPECTED_CONFIG_VERSION)
	}

	logger.CfgLog.Infof("config version [%s]", currentVersion)

	return nil
}
