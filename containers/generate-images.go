package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Environment represents the structure for each Environment configuration in the YAML file
type Environment struct {
	Name				string   `yaml:"name"`
	Type				string   `yaml:"type"`
	DockerfileExt		string   `yaml:"dockerfile_ext"`
	Base				string   `yaml:"base"`
	Maintainer			string   `yaml:"maintainer"`
	SteamAppID			string   `yaml:"steam_app_id"`
	AdditionalDeps		[]string `yaml:"additional_dependencies"`
}

// Config represents the structure of the entire YAML config file
type Config struct {
	Environment []Environment `yaml:"environment"`
}

// TemplateData is used to pass data to the Dockerfile template
type TemplateData struct {
	EnvironmentName		string
	Type				string
	BaseDockerfile    	string
	Maintainer         	string
	SteamAppID  	   	string
	AdditionalDeps     	[]string
}

func main() {
	// Parse command-line flags
	baseVersion := flag.String("base-version", "0.2.3", "Base image version")
	tagSuffix := flag.String("tag-suffix", "", "Tag suffix (e.g., -alpha, -beta)")
	flag.Parse()

	// Load the YAML file
	data, err := ioutil.ReadFile("environment.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	// Parse the YAML file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error unmarshalling YAML: %v\n", err)
		return
	}

	// Create the images directory if it doesn't exist
	err = os.MkdirAll("./images", os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating images directory: %v\n", err)
		return
	}

	// Load the Dockerfile template from file
	templateContent, err := ioutil.ReadFile("Dockerfile.j2")
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		return
	}

	// Parse the template
	tmpl, err := template.New("dockerfile").Parse(string(templateContent))
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		return
	}

	// Generate Dockerfiles for each environment
	for _, environment := range config.Environment {
		baseDockerfile := fmt.Sprintf("kubelize/game-servers:%s-%s%s", *baseVersion, environment.Base, *tagSuffix)
		outputFileName := fmt.Sprintf("Dockerfile.%s", environment.DockerfileExt)

		// Create the output file
		outputFile, err := os.Create(fmt.Sprintf("./images/%s", outputFileName))
		if err != nil {
			fmt.Printf("Error creating Dockerfile for %s: %v\n", environment.Name, err)
			continue
		}
		defer outputFile.Close()

		// Prepare the template data
		templateData := TemplateData{
			EnvironmentName:	environment.Name,
			Type:				environment.Type,
			BaseDockerfile: 	baseDockerfile,
			Maintainer:     	environment.Maintainer,
			SteamAppID: 		environment.SteamAppID,
			AdditionalDeps: 	environment.AdditionalDeps,
		}

		// Render the template and write to the file
		err = tmpl.Execute(outputFile, templateData)
		if err != nil {
			fmt.Printf("Error writing Dockerfile for %s: %v\n", environment.Name, err)
			continue
		}

		fmt.Printf("Dockerfile for %s created at ./images/%s\n", environment.Name, outputFileName)
	}
}
