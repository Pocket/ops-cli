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
}

func NewSettings(jsonPath string) *Settings {
	// Open our jsonFile
	jsonFile, err := os.Open(jsonPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic("Error reading json file, " + jsonPath + ", " + err.Error())
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
	settings.setName(stackName)
	settings.setFilePath(templatefilePath)

	if gitSHA != nil {
		settings.setGitSHA(gitSHA)
	}

	return settings
}

func (settings *Settings) setGitSHA(gitSHA *string)  {
	for _, tag := range settings.Tags {
		if *tag.Key == "GitSHA" {
			tag.Value = gitSHA
		}
	}

	for _, parameter := range settings.Parameters {
		if *parameter.ParameterKey == "GitSHA" {
			parameter.ParameterValue = gitSHA
		}
	}
}

func (settings *Settings) setName(stackName *string)  {
	settings.StackName = stackName
}

func (settings *Settings) setFilePath(filePath string)  {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Error reading json file, " + filePath + ", " + err.Error())
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	template := string(byteValue)
	settings.TemplateBody = &template
}

