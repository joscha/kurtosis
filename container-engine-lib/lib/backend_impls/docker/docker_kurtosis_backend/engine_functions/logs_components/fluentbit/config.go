package fluentbit

import (
	"fmt"
	"github.com/kurtosis-tech/container-engine-lib/lib/backend_impls/docker/object_attributes_provider/label_key_consts"
	"strings"
)

const (
	filterRulesSeparator  = "\n	"
	outputLabelsSeparator = ", "
)

type Config struct{
	Service *Service
	Input *Input
	Filter *Filter
	Output *Output

}

type Service struct {
	LogLevel string
	HttpServerEnabled string
	HttpServerHost string
	HttpServerPort uint16
}

type Input struct {
	Name string
	Listen string
	Port uint16
}

type Filter struct {
	Name string
	Match string
	Rules []string
}

type Output struct {
	Name string
	Match string
	Host string
	Port uint16
	Labels []string
	LineFormat string
	TenantIDKey string
}

func newDefaultConfigForKurtosisCentralizedLogsForDocker(
	lokiHost string,
	lokiPort uint16,
) *Config {
	return &Config{
		Service: &Service {
			LogLevel: logLevel,
			HttpServerEnabled: httpServerEnabledValue,
			HttpServerHost: httpServerLocalhost,
			HttpServerPort: httpPortNumber,
		},
		Input: &Input{
			Name: inputName,
			Listen: inputListenIP,
			Port: tcpPortNumber,
		},
		Filter: &Filter{
			Name:   modifyFilterName,
			Match:  matchAllRegex,
			Rules : getModifyFilterRulesKurtosisLabels(),
		},
		Output: &Output{
			Name:        lokiOutputTypeName,
			Match:       matchAllRegex,
			Host:        lokiHost,
			Port:        lokiPort,
			Labels:      getOutputKurtosisLabelsForLogs(),
			LineFormat:  jsonLineFormat,
			TenantIDKey: getTenantIdKeyFromKurtosisLabels(),
		},
	}
}

func (filter *Filter) GetRulesStr() string {
	return strings.Join(filter.Rules, filterRulesSeparator)
}

func (output *Output) GetLabelsStr() string {
	return strings.Join(output.Labels, outputLabelsSeparator)
}

func getModifyFilterRulesKurtosisLabels() []string {

	modifyFilterRules := []string{}

	modifyFilterRuleAction := "rename"

	kurtosisLabelsForLogs := getTrackedKurtosisLabelsForLogs()
	for _, kurtosisLabel := range kurtosisLabelsForLogs {
		validFormatLabelValue := newValidFormatLabelValue(kurtosisLabel)
		modifyFilterRule := fmt.Sprintf("%v %v %v", modifyFilterRuleAction, kurtosisLabel, validFormatLabelValue)
		modifyFilterRules = append(modifyFilterRules, modifyFilterRule)
	}
	return modifyFilterRules
}

func getOutputKurtosisLabelsForLogs() []string {
	outputLabels := []string{}
	labelsVarPrefix := "$"

	kurtosisLabelsForLogs := getTrackedKurtosisLabelsForLogs()
	for _, kurtosisLabel := range kurtosisLabelsForLogs {
		validFormatLabelValue := newValidFormatLabelValue(kurtosisLabel)
		outputLabel := fmt.Sprintf("%v%v", labelsVarPrefix, validFormatLabelValue)
		outputLabels = append(outputLabels, outputLabel)
	}
	return outputLabels
}

func getTrackedKurtosisLabelsForLogs() []string {
	kurtosisLabelsForLogs := []string{
		label_key_consts.GUIDDockerLabelKey.GetString(),
		label_key_consts.ContainerTypeDockerLabelKey.GetString(),
	}
	return kurtosisLabelsForLogs
}

func getTenantIdKeyFromKurtosisLabels() string {
	return label_key_consts.EnclaveIDDockerLabelKey.GetString()
}


func newValidFormatLabelValue(stringToModify string) string {
	notAllowedCharInLabels := " .-"
	noSeparationChar := ""
	lowerString := strings.ToLower(stringToModify)
	shouldChangeNextCharToUpperCase := false
	shouldChangeCharToUpperCase := false
	var newString string
	for _, currenChar := range strings.Split(lowerString, noSeparationChar) {
		newChar := currenChar
		shouldChangeCharToUpperCase = shouldChangeNextCharToUpperCase
		if shouldChangeCharToUpperCase {
			newChar = strings.ToUpper(newChar)
		}
		if strings.ContainsAny(currenChar, notAllowedCharInLabels) {
			shouldChangeNextCharToUpperCase = true
		} else {
			shouldChangeNextCharToUpperCase = false
			newString = newString + newChar
		}
	}
	return newString
}
