/*
 * === This file is part of ALICE O² ===
 *
 * Copyright 2021-2022 CERN and copyright holders of ALICE O².
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

//go:generate protoc --go_out=. --go-grpc_out=require_unimplemented_servers=false:. protos/odc.proto

package odc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/AliceO2Group/Control/apricot"
	"github.com/AliceO2Group/Control/common/gera"
	"github.com/AliceO2Group/Control/common/logger/infologger"
	"github.com/AliceO2Group/Control/common/utils/uid"
	"github.com/AliceO2Group/Control/core/integration"
	odc "github.com/AliceO2Group/Control/core/integration/odc/protos"
	"github.com/AliceO2Group/Control/core/workflow/callable"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

const(
	ODC_DIAL_TIMEOUT = 2 * time.Second
	ODC_GENERAL_OP_TIMEOUT = 5 * time.Second
	ODC_CONFIGURE_TIMEOUT = 60 * time.Second
	ODC_START_TIMEOUT = 15 * time.Second
	ODC_STOP_TIMEOUT = 15 * time.Second
	ODC_RESET_TIMEOUT = 30 * time.Second
)


type Plugin struct {
	odcHost string
	odcPort int

	odcClient *RpcClient
}

func NewPlugin(endpoint string) integration.Plugin {
	u, err := url.Parse(endpoint)
	if err != nil {
		log.WithField("endpoint", endpoint).
			WithError(err).
			Error("bad service endpoint")
		return nil
	}

	portNumber, _ := strconv.Atoi(u.Port())

	return &Plugin{
		odcHost: u.Hostname(),
		odcPort: portNumber,
		odcClient:   nil,
	}
}

func (p *Plugin) GetName() string {
	return "odc"
}

func (p *Plugin) GetPrettyName() string {
	return "ODC"
}

func (p *Plugin) GetEndpoint() string {
	return viper.GetString("odcEndpoint")
}

func (p *Plugin) GetConnectionState() string {
	if p == nil || p.odcClient == nil {
		return "UNKNOWN"
	}
	return p.odcClient.conn.GetState().String()
}

func (p *Plugin) GetData(environmentIds []uid.ID) string {
	if p == nil || p.odcClient == nil {
		return ""
	}

	partitionStates := make(map[string]string)

	for _, envId := range environmentIds {
		in := odc.StateRequest{
			Partitionid: envId.String(),
			Path:        "",
			Detailed:    false,
		}
		state, err := p.odcClient.GetState(context.Background(), &in, grpc.EmptyCallOption{})
		if err != nil {
			continue
		}
		if state == nil || state.Reply == nil {
			continue
		}
		partitionStates[envId.String()] = state.Reply.State
	}

	out, err := json.Marshal(partitionStates)
	if err != nil {
		return ""
	}
	return string(out[:])
}

func (p *Plugin) Init(_ string) error {
	if p.odcClient == nil {
		cxt, cancel := context.WithCancel(context.Background())
		p.odcClient = NewClient(cxt, cancel, viper.GetString("odcEndpoint"))
		if p.odcClient == nil {
			return fmt.Errorf("failed to connect to ODC service on %s", viper.GetString("ddSchedulerEndpoint"))
		}
		log.Debug("ODC plugin initialized")
	}
	return nil
}

func (p *Plugin) ObjectStack(varStack map[string]string) (stack map[string]interface{}) {
	envId, ok := varStack["environment_id"]
	if !ok {
		log.Error("cannot acquire environment ID")
		return
	}

	var csErr error
	configStack := apricot.Instance().GetDefaults()
	configStack, csErr = gera.MakeStringMapWithMap(apricot.Instance().GetVars()).WrappedAndFlattened(gera.MakeStringMapWithMap(configStack))
	if csErr != nil {
		log.Error("cannot access AliECS workflow configuration defaults")
		return
	}

	stack = make(map[string]interface{})
	stack["GenerateEPNWorkflowScript"] = func() (out string) {
		/*
		OCTRL-558 example:
		GEN_TOPO_HASH=[0/1] GEN_TOPO_SOURCE=[...] DDMODE=[TfBuilder Mode] GEN_TOPO_LIBRARY_FILE=[...]
		GEN_TOPO_WORKFLOW_NAME=[...] WORKFLOW_DETECTORS=[...] WORKFLOW_DETECTORS_QC=[...]
		WORKFLOW_DETECTORS_CALIB=[...] WORKFLOW_PARAMETERS=[...] RECO_NUM_NODES_OVERRIDE=[...]
		MULTIPLICITY_FACTOR_RAWDECODERS=[...] MULTIPLICITY_FACTOR_CTFENCODERS=[...]
		MULTIPLICITY_FACTOR_REST=[...] GEN_TOPO_WIPE_CACHE=[0/1] BEAMTYPE=[PbPb/pp/pPb/cosmic/technical]
		NHBPERTF=[...] GEN_TOPO_PARTITION=[...] GEN_TOPO_ONTHEFLY=1 [Extra environment variables]
		/home/epn/pdp/gen_topo.sh

		R3C-710:
		`pdp_o2pdpsuite_version` is a new field. Its content should be sent in the string as `OVERRIDE_PDPSUITE_VERSION=[...]`.
			In case it is set to `default`, instead of the string `default` the preconfigured default version in consul should be sent.
		`pdp_qcjson_version`: similar to avove, new field. please send as `SET_QCJSON_VERSION`.
			If set to the string `default`, please sent the default version configured in consul instead.
		`pdp_o2_data_processing_hash`: if set to the string `default`, sent the default hash configured in consul instead.
		`odc_n_epns_max_fail` : new field. Please send as `RECO_MAX_FAIL_NODES_OVERRIDE=[...]`.
		`epn_store_raw_data_fraction` new field, please send as `DD_DISK_FRACTION=[...]`.
		`pdp_nr_compute_nodes` removed this field since no longer needed.
			Please send the value of `odc_n_epns` directly as `RECO_NUM_NODES_OVERRIDE=[...]`.
		`pdp_epn_shmid`: new field, please send as `SHM_MANAGER_SHMID=[...]`
		`pdp_epn_shm_recreate`: new field, please send as `SHM_MANAGER_SHM_RECREATE=[0|1]`
		*/

		var (
			pdpConfigOption, o2DPSource, tfbDDMode string
			pdpLibraryFile, pdpLibWorkflowName string
			pdpDetectorList, pdpDetectorListQc, pdpDetectorListCalib string
			pdpWorkflowParams string
			pdpRawDecoderMultiFactor, pdpCtfEncoderMultiFactor, pdpRecoProcessMultiFactor string
			pdpWipeWorkflowCache, pdpBeamType, pdpNHbfPerTf string
			pdpExtraEnvVars, pdpGeneratorScriptPath string
			odcNEpns string
			ok bool
			accumulator []string
			pdpO2PdpSuiteVersion, pdpQcJsonVersion string
			odcNEpnsMaxFail, epnStoreRawDataFraction string
			pdpEpnShmId, pdpEpnShmRecreate string
		)
		accumulator = make([]string, 0)

		pdpConfigOption, ok = varStack["pdp_config_option"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow configuration mode")
			return
		}

		switch pdpConfigOption {
		case "Repository hash":
			o2DPSource, ok = varStack["pdp_o2_data_processing_hash"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire PDP Repository hash")
				return
			}
			if strings.TrimSpace(o2DPSource) == "default" {	// if UI sends 'default', we look in Consul
				pdpO2PdpSuiteVersion, ok = configStack["pdp_o2_data_processing_hash"]
				if !ok {
					log.WithField("partition", envId).
						WithField("call", "GenerateEPNWorkflowScript").
						Error("cannot acquire PDP Repository hash default")
					return
				}
			}
			accumulator = append(accumulator, "GEN_TOPO_HASH=1")

		case "Repository path":
			o2DPSource, ok = varStack["pdp_o2_data_processing_path"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire PDP Repository path")
				return
			}
			accumulator = append(accumulator, "GEN_TOPO_HASH=0")

		case "Manual XML": fallthrough
		default:
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("GEN_TOPO_SOURCE='%s'", strings.TrimSpace(o2DPSource)))

		tfbDDMode, ok = varStack["tfb_dd_mode"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire TF Builder mode")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("DDMODE='%s'", strings.TrimSpace(tfbDDMode)))

		pdpLibraryFile, ok = varStack["pdp_topology_description_library_file"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire topology description library file")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("GEN_TOPO_LIBRARY_FILE='%s'", strings.TrimSpace(pdpLibraryFile)))

		pdpLibWorkflowName, ok = varStack["pdp_workflow_name"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow name in topology library file")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("GEN_TOPO_WORKFLOW_NAME='%s'", strings.TrimSpace(pdpLibWorkflowName)))

		pdpDetectorList, ok = varStack["pdp_detector_list_global"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow name in topology library file")
			return
		}
		if strings.TrimSpace(pdpDetectorList) == "default" {
			pdpDetectorList, ok = varStack["detectors"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire general detector list from varStack")
				return
			}
			detectorsSlice, err := p.parseDetectors(pdpDetectorList)
			if err != nil {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot parse general detector list")
				return
			}
			pdpDetectorList = strings.Join(detectorsSlice, ",")
		}
		accumulator = append(accumulator, fmt.Sprintf("WORKFLOW_DETECTORS='%s'", strings.TrimSpace(pdpDetectorList)))

		pdpDetectorListQc, ok = varStack["pdp_detector_list_qc"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow name in topology library file")
			return
		}
		if strings.TrimSpace(pdpDetectorListQc) == "default" {
			pdpDetectorListQc, ok = varStack["detectors"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire general detector list from varStack")
				return
			}
			detectorsSlice, err := p.parseDetectors(pdpDetectorListQc)
			if err != nil {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot parse general detector list")
				return
			}
			pdpDetectorListQc = strings.Join(detectorsSlice, ",")
		}
		accumulator = append(accumulator, fmt.Sprintf("WORKFLOW_DETECTORS_QC='%s'", strings.TrimSpace(pdpDetectorListQc)))

		pdpDetectorListCalib, ok = varStack["pdp_detector_list_calib"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow name in topology library file")
			return
		}
		if strings.TrimSpace(pdpDetectorListCalib) == "default" {
			pdpDetectorListCalib, ok = varStack["detectors"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire general detector list from varStack")
				return
			}
			detectorsSlice, err := p.parseDetectors(pdpDetectorListCalib)
			if err != nil {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot parse general detector list")
				return
			}
			pdpDetectorListCalib = strings.Join(detectorsSlice, ",")
		}
		accumulator = append(accumulator, fmt.Sprintf("WORKFLOW_DETECTORS_CALIB='%s'", strings.TrimSpace(pdpDetectorListCalib)))

		pdpWorkflowParams, ok = varStack["pdp_workflow_parameters"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow parameters")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("WORKFLOW_PARAMETERS='%s'", strings.TrimSpace(pdpWorkflowParams)))

		odcNEpns, ok = varStack["odc_n_epns"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire ODC number of EPNs")
			return
		}
		odcNEpnsI, err := strconv.Atoi(odcNEpns)
		if err != nil {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot parse ODC number of EPNs")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("RECO_NUM_NODES_OVERRIDE=%d", odcNEpnsI))

		odcNEpnsMaxFail, ok = varStack["odc_n_epns_max_fail"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire ODC number of EPNs max fail")
			return
		}
		odcNEpnsMaxFailI, err := strconv.Atoi(odcNEpnsMaxFail)
		if err != nil {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot parse ODC number of EPNs max fail")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("RECO_MAX_FAIL_NODES_OVERRIDE=%d", odcNEpnsMaxFailI))

		pdpRawDecoderMultiFactor, ok = varStack["pdp_raw_decoder_multi_factor"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP number of raw decoder processing instances")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("MULTIPLICITY_FACTOR_RAWDECODERS=%s", strings.TrimSpace(pdpRawDecoderMultiFactor)))

		pdpCtfEncoderMultiFactor, ok = varStack["pdp_ctf_encoder_multi_factor"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP number of CTF encoder processing instances")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("MULTIPLICITY_FACTOR_CTFENCODERS=%s", strings.TrimSpace(pdpCtfEncoderMultiFactor)))

		pdpRecoProcessMultiFactor, ok = varStack["pdp_reco_process_multi_factor"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP number of other reconstruction processing instances")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("MULTIPLICITY_FACTOR_REST=%s", strings.TrimSpace(pdpRecoProcessMultiFactor)))

		pdpWipeWorkflowCache, ok = varStack["pdp_wipe_workflow_cache"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP workflow cache wipe option")
			return
		}
		pdpWipeWorkflowCacheB, err := strconv.ParseBool(pdpWipeWorkflowCache)
		if err != nil {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot parse PDP workflow cache wipe option")
			pdpWipeWorkflowCacheB = false
		}
		pdpWipeWorkflowCacheI := 0
		if pdpWipeWorkflowCacheB {
			pdpWipeWorkflowCacheI = 1
		}
		accumulator = append(accumulator, fmt.Sprintf("GEN_TOPO_WIPE_CACHE=%d", pdpWipeWorkflowCacheI))

		pdpBeamType, ok = varStack["pdp_beam_type"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire beam type")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("BEAMTYPE='%s'", strings.TrimSpace(pdpBeamType)))

		pdpNHbfPerTf, ok = varStack["pdp_n_hbf_per_tf"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire number of HBFs per TF")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("NHBPERTF=%s", strings.TrimSpace(pdpNHbfPerTf)))

		// envId
		accumulator = append(accumulator, fmt.Sprintf("GEN_TOPO_PARTITION='%s'", envId))

		accumulator = append(accumulator, "GEN_TOPO_ONTHEFLY=1")

		pdpO2PdpSuiteVersion, ok = varStack["pdp_o2pdpsuite_version"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP Suite version")
			return
		}
		if strings.TrimSpace(pdpO2PdpSuiteVersion) == "default" {	// if UI sends 'default', we look in Consul
			pdpO2PdpSuiteVersion, ok = configStack["pdp_o2pdpsuite_version"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire PDP Suite version default")
				return
			}
		}
		accumulator = append(accumulator, fmt.Sprintf("OVERRIDE_PDPSUITE_VERSION='%s'", pdpO2PdpSuiteVersion))

		pdpQcJsonVersion, ok = varStack["pdp_qcjson_version"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP QCJson version")
			return
		}
		if strings.TrimSpace(pdpQcJsonVersion) == "default" {	// if UI sends 'default', we look in Consul
			pdpQcJsonVersion, ok = configStack["pdp_qcjson_version"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "GenerateEPNWorkflowScript").
					Error("cannot acquire PDP QCJson version default")
				return
			}
		}
		accumulator = append(accumulator, fmt.Sprintf("SET_QCJSON_VERSION='%s'", pdpQcJsonVersion))

		epnStoreRawDataFraction, ok = varStack["epn_store_raw_data_fraction"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire EPN DD disk raw data fraction")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("DD_DISK_FRACTION='%s'", epnStoreRawDataFraction))

		pdpEpnShmId, ok = varStack["pdp_epn_shmid"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP EPN SHMID")
			return
		}
		accumulator = append(accumulator, fmt.Sprintf("SHM_MANAGER_SHMID='%s'", pdpEpnShmId))

		pdpEpnShmRecreate, ok = varStack["pdp_epn_shm_recreate"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP EPN SHM recreate")
			return
		}
		pdpEpnShmRecreateB, err := strconv.ParseBool(pdpEpnShmRecreate)
		if err != nil {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot parse PDP EPN SHM recreate")
			pdpEpnShmRecreateB = false
		}
		pdpEpnShmRecreateI := 0
		if pdpEpnShmRecreateB {
			pdpEpnShmRecreateI = 1
		}
		accumulator = append(accumulator, fmt.Sprintf("SHM_MANAGER_SHM_RECREATE=%d", pdpEpnShmRecreateI))

		pdpExtraEnvVars, ok = varStack["pdp_extra_env_vars"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP extra environment variables")
			return
		}
		accumulator = append(accumulator, strings.TrimSpace(pdpExtraEnvVars))

		pdpGeneratorScriptPath, ok = varStack["pdp_generator_script_path"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "GenerateEPNWorkflowScript").
				Error("cannot acquire PDP generator script path")
			return
		}
		accumulator = append(accumulator, strings.TrimSpace(pdpGeneratorScriptPath))

		out = strings.Join(accumulator, " ")
		return
	}
	return stack
}

func (p *Plugin) CallStack(data interface{}) (stack map[string]interface{}) {
	call, ok := data.(*callable.Call)
	if !ok {
		return
	}
	varStack := call.VarStack
	envId, ok := varStack["environment_id"]
	if !ok {
		log.Error("cannot acquire environment ID")
		return
	}

	stack = make(map[string]interface{})
	stack["Configure"] = func() (out string) {
		// ODC Run + SetProperties + Configure

		var (
			pdpConfigOption, script, topology, plugin, resources string
		)
		ok := false
		isManualXml := false

		pdpConfigOption, ok = varStack["pdp_config_option"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "Configure").
				Error("cannot acquire PDP workflow configuration mode")
			return
		}
		switch pdpConfigOption {
		case "Repository hash": fallthrough
		case "Repository path":
			script, ok = varStack["odc_script"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "Configure").
					Error("cannot acquire ODC script, make sure GenerateEPNWorkflowScript is called and its " +
						"output is written to odc_script")
				return
			}

		case "Manual XML":
			topology, ok = varStack["odc_topology"]
			if !ok {
				log.WithField("partition", envId).
					WithField("call", "Configure").
					Error("cannot acquire ODC topology")
				return
			}
			isManualXml = true

		default:
			log.WithField("partition", envId).
				WithField("call", "Configure").
				WithField("value", pdpConfigOption).
				Error("cannot acquire valid PDP workflow configuration mode value")
			return
		}

		plugin, ok = varStack["odc_plugin"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "Configure").
				Error("cannot acquire ODC RMS plugin declaration")
			return
		}

		resources, ok = varStack["odc_resources"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "Configure").
				Error("cannot acquire ODC resources declaration")
			return
		}

		timeout := callable.AcquireTimeout(ODC_CONFIGURE_TIMEOUT, varStack, "Configure", envId)

		arguments := make(map[string]string)
		arguments["environment_id"] = envId

		// FIXME: this only copies over vars prefixed with "odc_"
		// Figure out a better way!
		for k, v := range varStack {
			if strings.HasPrefix(k, "odc_") &&
				k != "odc_enabled" &&
				k != "odc_resources" &&
				k != "odc_plugin" &&
				k != "odc_script" &&
				k != "odc_topology" {
				arguments[strings.TrimPrefix(k, "odc_")] = v
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := handleConfigure(ctx, p.odcClient, arguments, isManualXml, topology, script, plugin, resources, envId)
		if err != nil {
			log.WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "Configure").
				WithError(err).Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "Configure").
				Error("EPN Configure call failed")
		}
		return
	}
	stack["Start"] = func() (out string) {	// must formally return string even when we return nothing
		// ODC SetProperties + Start

		rn, ok := varStack["run_number"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "Start").
				Warn("cannot acquire run number for ODC")
		}
		var (
			runNumberu64 uint64
			err error
		)

		runNumberu64, err = strconv.ParseUint(rn, 10, 32)
		if err != nil {
			log.WithField("partition", envId).
				WithError(err).
				Error("cannot acquire run number for DCS SOR")
			runNumberu64 = 0
		}

		timeout := callable.AcquireTimeout(ODC_START_TIMEOUT, varStack, "Start", envId)

		arguments := make(map[string]string)
		arguments["run_number"] = rn
		arguments["runNumber"] = rn

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err = handleStart(ctx, p.odcClient, arguments, envId, runNumberu64)
		if err != nil {
			log.WithError(err).
				WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "Start").
				Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "Start").
				Error("EPN Start call failed")
		}
		return
	}
	stack["Stop"] = func() (out string) {
		// ODC Stop

		rn, ok := varStack["run_number"]
		if !ok {
			log.WithField("partition", envId).
				WithField("call", "Start").
				Warn("cannot acquire run number for ODC")
		}
		var (
			runNumberu64 uint64
			err error
		)

		runNumberu64, err = strconv.ParseUint(rn, 10, 32)
		if err != nil {
			log.WithField("partition", envId).
				WithError(err).
				Error("cannot acquire run number for DCS SOR")
			runNumberu64 = 0
		}

		timeout := callable.AcquireTimeout(ODC_STOP_TIMEOUT, varStack, "Stop", envId)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err = handleStop(ctx, p.odcClient, nil, envId, runNumberu64)
		if err != nil {
			log.WithError(err).
				WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "Stop").
				Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "Stop").
				Error("EPN Stop call failed")
		}
		return
	}
	stack["Reset"] = func() (out string) {
		// ODC Reset + Terminate + Shutdown

		timeout := callable.AcquireTimeout(ODC_RESET_TIMEOUT, varStack, "Reset", envId)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := handleReset(ctx, p.odcClient, nil, envId)
		if err != nil {
			log.WithError(err).
				WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "Reset").
				Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "Reset").
				Error("EPN Reset call failed")
		}
		return
	}
	stack["EnsureCleanupLegacy"] = func() (out string) {
		// ODC Reset + Terminate + Shutdown for current env

		timeout := callable.AcquireTimeout(ODC_GENERAL_OP_TIMEOUT, varStack, "EnsureCleanupLegacy", envId)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := handleCleanupLegacy(ctx, p.odcClient, nil, envId)
		if err != nil {
			log.WithError(err).
				WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "EnsureCleanupLegacy").

				Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "EnsureCleanupLegacy").
				Error("EPN Cleanup sequence failed")
		}
		return
	}
	stack["EnsureCleanup"] = func() (out string) {
		// ODC Shutdown for current env + all orphans

		timeout := callable.AcquireTimeout(ODC_GENERAL_OP_TIMEOUT, varStack, "EnsureCleanup", envId)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := handleCleanup(ctx, p.odcClient, nil, envId)
		if err != nil {
			log.WithError(err).
				WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "EnsureCleanup").

				Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "EnsureCleanup").
				Error("EPN Cleanup sequence failed")
		}
		return
	}
	stack["PreDeploymentCleanup"] = func() (out string) {
		// ODC Shutdown for all orphans

		timeout := callable.AcquireTimeout(ODC_GENERAL_OP_TIMEOUT, varStack, "PreDeploymentCleanup", envId)

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		err := handleCleanup(ctx, p.odcClient, nil, "")
		if err != nil {
			log.WithError(err).
				WithField("level", infologger.IL_Support).
				WithField("partition", envId).
				WithField("call", "PreDeploymentCleanup").

				Error("ODC error")
			log.WithField("partition", envId).
				WithField("call", "PreDeploymentCleanup").
				Error("EPN PreDeploymentCleanup sequence failed")
		}
		return
	}
	return
}

func (p *Plugin) parseDetectors(detectorsParam string) (detectors []string, err error) {
	detectorsSlice := make([]string, 0)
	bytes := []byte(detectorsParam)
	err = json.Unmarshal(bytes, &detectorsSlice)
	if err != nil {
		log.WithError(err).
			Error("error processing EPN/PDP detectors list")
		return
	}
	detectors = detectorsSlice
	return
}

func (p *Plugin) Destroy() error {
	return p.odcClient.Close()
}
