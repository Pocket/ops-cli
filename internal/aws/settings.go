package aws

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"io/ioutil"
	"os"
)

type Settings struct {
	Parameters       []cloudformation.Parameter  `json:"parameters,ommitempty"`
	Tags             []cloudformation.Tag        `json:"tags,ommitempty"`
	StackName        *string                     `json:"stack_name,ommitempty"`
	OnFailure        cloudformation.OnFailure    `json:"on_failure,ommitempty"`
	Capabilities     []cloudformation.Capability `json:"capabilities,ommitempty"`
	TemplateBody     *string                     `json:"template_body,ommitempty"`
	TimeoutInMinutes *int64                      `json:"timeout_in_minutes,ommitempty"`
	ECSCluster       *string                     `json:"ecs_cluster,ommitempty"`
}

func NewSettings(jsonPath string) *Settings {
	// Open our jsonFile
	jsonFile, err := os.Open(jsonPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic("Error reading parameters file, " + jsonPath + ", " + err.Error())
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var settings Settings

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'pocketParams' which we defined above
	json.Unmarshal(byteValue, &settings)

	return &settings
}

func NewSettingsParams(paramFilePath string, stackName *string, templatefilePath string, gitSHA *string) *Settings {
	settings := NewSettings(paramFilePath)

	if stackName != nil {
		branchName := *stackName
		settings.setBranchName(&branchName)
	}

	if settings.StackName != nil && stackName != nil {
		stackName2 := *settings.StackName + *stackName
		stackName = &stackName2
	}
	if stackName != nil {
		settings.setName(stackName)
	}
	settings.setFilePath(templatefilePath)

	if gitSHA != nil {
		settings.setGitSHA(gitSHA)
	}

	return settings
}

func (settings *Settings) setGitSHA(gitSHA *string) {
	settings.replaceTag("GitSHA", gitSHA)
	settings.replaceParameter("GitSHA", gitSHA)
}

func (settings *Settings) setBranchName(branchName *string) {
	settings.replaceTag("BranchName", branchName)
	settings.replaceParameter("BranchName", branchName)
}

func (settings *Settings) replaceParameter(key string, value *string) {
	var parameters []cloudformation.Parameter

	for _, parameter := range settings.Parameters {
		if *parameter.ParameterKey == key {
			parameter.ParameterValue = value
		}
		parameters = append(parameters, parameter)
	}
	settings.Parameters = parameters
}

func (settings *Settings) getParameter(key string) *string {
	for _, parameter := range settings.Parameters {
		if *parameter.ParameterKey == key {
			return parameter.ParameterValue
		}
	}
	return nil
}

func (settings *Settings) replaceTag(key string, value *string) {
	var tags []cloudformation.Tag

	for _, tag := range settings.Tags {
		if *tag.Key == key {
			tag.Value = value
		}
		tags = append(tags, tag)
	}
	settings.Tags = tags
}

func (settings *Settings) setName(stackName *string) {
	settings.StackName = stackName
}

func (settings *Settings) getBaseUrl() *string {
	return settings.getParameter("DomainBase")
}

func (settings *Settings) GetDeployUrl() *string {
	base := settings.getParameter("DomainBase")
	formattedBranch := settings.getParameter("BranchName")
	if base == nil || formattedBranch == nil {
		return nil
	}

	url := "https://" + *formattedBranch + "." + *base
	return &url
}

func (settings *Settings) setFilePath(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Error reading template file, " + filePath + ", " + err.Error())
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	template := string(byteValue)
	settings.TemplateBody = &template
}
