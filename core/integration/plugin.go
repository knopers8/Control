/*
 * === This file is part of ALICE O² ===
 *
 * Copyright 2021 CERN and copyright holders of ALICE O².
 * Author: Teo Mrnjavac <teo.mrnjavac@cern.ch>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * In applying this license CERN does not waive the privileges and
 * immunities granted to it by virtue of its status as an
 * Intergovernmental Organization or submit itself to any jurisdiction.
 */

package integration

import (
	"sync"

	"github.com/AliceO2Group/Control/common/logger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var log = logger.New(logrus.StandardLogger(), "integration")

var(
	once     sync.Once
	pluginLoaders map[string]func() Plugin
	instance Plugins
)

type Plugins []Plugin

type Plugin interface {
	GetName() string
	Init(instanceId string) error
	ObjectStack(data interface{}) map[string]interface{}
	Destroy() error
}

type NewFunc func(endpoint string) Plugin

func RegisterPlugin(pluginName string, endpointArgumentName string, newFunc NewFunc) {
	once.Do(func() {
		pluginLoaders = make(map[string]func() Plugin)
	})
	pluginLoaders[pluginName] = func() Plugin {
		if viper.IsSet(endpointArgumentName) {
			endpoint := viper.GetString(endpointArgumentName)
			return newFunc(endpoint)
		}
		return nil
	}
}

func (p Plugins) InitAll(fid string) {
	for _, plugin := range p {
		initErr := plugin.Init(fid)
		if initErr != nil {
			log.WithError(initErr).
				WithField("plugin", plugin.GetName()).
				Error("workflow plugin failed to initialize")
		}
	}
}

func (p Plugins) DestroyAll() {
	for _, plugin := range p {
		err := plugin.Destroy()
		if err != nil {
			log.WithError(err).
				WithField("plugin", plugin.GetName()).
				Error("workflow plugin failed to destroy")
		}
	}
}

func (p Plugins) ObjectStack(data interface{}) (stack map[string]interface{}) {
	stack = make(map[string]interface{})
	for _, plugin := range p {
		s := plugin.ObjectStack(data)
		stack[plugin.GetName()] = s
	}
	return
}

func PluginsInstance() Plugins {
	once.Do(func() {
		instance = Plugins{}
		pluginList := viper.GetStringSlice("integrationPlugins")

		for _, pluginName := range pluginList {
			if pluginLoaders == nil {
				log.WithField("plugin", pluginName).
					Error("requested plugin unavailable")
				continue
			}
			pluginLoader, ok := pluginLoaders[pluginName]
			if !ok {
				log.WithField("plugin", pluginName).
					Error("requested plugin unavailable")
				continue
			}
			newPlugin := pluginLoader()
			if newPlugin == nil {
				log.WithField("plugin", pluginName).
					Error("plugin loader failed")
				continue
			}
			instance = append(instance, newPlugin)
		}
	})
	return instance
}
