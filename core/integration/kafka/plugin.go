/*
 * === This file is part of ALICE O² ===
 *
 * Copyright 2022 CERN and copyright holders of ALICE O².
 * Author: Piotr Konopka <piotr.jan.konopka@cern.ch>
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

//go:generate protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. protos/kafka.proto

package kafka

import (
	"fmt"
	"github.com/AliceO2Group/Control/common/utils/uid"
	"github.com/AliceO2Group/Control/core/integration"
	kafkapb "github.com/AliceO2Group/Control/core/integration/kafka/protos"
	"github.com/AliceO2Group/Control/core/workflow/callable"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
)

type KafkaEnvInfo kafkapb.EnvInfo

type Plugin struct {
	kafkaBroker   string
	kafkaProducer *kafka.Producer
}

func NewPlugin(endpoint string) integration.Plugin {
	_, err := url.Parse(endpoint)
	if err != nil {
		log.WithField("endpoint", endpoint).
			WithError(err).
			Error("bad service endpoint")
		return nil
	}

	return &Plugin{
		kafkaBroker: endpoint,
	}
}

func (p *Plugin) GetName() string {
	return "kafka"
}

func (p *Plugin) GetPrettyName() string {
	return "Kafka FSM transition information service"
}

func (p *Plugin) GetEndpoint() string {
	return viper.GetString("kafkaEndpoint")
}

func (p *Plugin) GetConnectionState() string {
	return "READY"
}

func (p *Plugin) GetData(environmentIds []uid.ID) string {
	return ""
}

func (p *Plugin) Init(instanceId string) error {
	var err error
	p.kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return fmt.Errorf("Failed to initialize kafka producer with broker '" + p.kafkaBroker + "'. Details: " + err.Error())
	}
	log.Info("Successfully created a kafka producer with broker '" + p.kafkaBroker + "'")
	return nil
}

func (p *Plugin) ObjectStack(_ map[string]string) (stack map[string]interface{}) {
	stack = make(map[string]interface{})
	return stack
}

func (p *Plugin) RetrieveRunState(varStack map[string]string) {
	envId, ok := varStack["environment_id"]
	if !ok {
		log.Error("cannot acquire environment ID")
		return
	}

	_, err := strconv.ParseUint(varStack["run_number"], 10, 32)
	if err != nil {
		log.WithError(err).
			WithField("partition", envId).
			Error("cannot acquire run number for Run Start")
		return
	}

	return
}

func (p *Plugin) CallStack(data interface{}) (stack map[string]interface{}) {
	call, ok := data.(*callable.Call)
	if !ok {
		return
	}
	varStack := call.VarStack
	_, ok = varStack["environment_id"]
	if !ok {
		log.Error("cannot acquire environment ID")
		return
	}

	stack = make(map[string]interface{})
	stack["Standby"] = func() (out string) {
		return
	}
	stack["Deployed"] = func() (out string) {
		return
	}
	stack["Configured"] = func() (out string) {
		return
	}
	stack["Running"] = func() (out string) {

		return
	}
	stack["Done"] = func() (out string) {
		return
	}
	stack["Error"] = func() (out string) {
		return
	}
	return
}

func (p *Plugin) Destroy() error {
	p.kafkaProducer.Close()
	return nil
}
