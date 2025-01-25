package main

import (
	"fmt"
	"os"
	"text/template"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

// Game represents the structure for each game configuration in the YAML file
type Game struct {
	Name                string   `yaml:"name"`
	SteamAppID          int      `yaml:"steam_app_id"`
	DockerfileExt       string   `yaml:"dockerfile_ext"`
	Base                string   `yaml:"base"`
	Maintainer          string   `yaml:"maintainer"`
	AdditionalDeps      []string `yaml:"additional_dependencies"`
}

// Config represents the structure of the entire YAML config file
type Config struct {
	Games []Game `yaml:"games"`
}

// TemplateData is used to pass data to the Dockerfile template
type TemplateData struct {
	GameName           string
	BaseDockerfile     string
	Maintainer         string
	SteamAppID         int
	AdditionalDeps     []string
}

func main() {
	// Load the YAML file
	data, err := ioutil.ReadFile("games.yaml")
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

	// Create the game-servers directory if it doesn't exist
	err = os.MkdirAll("./game-servers", os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating game-servers directory: %v\n", err)
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

	// Generate Dockerfiles for each game
	for _, game := range config.Games {
		baseDockerfile := fmt.Sprintf("Dockerfile.%s", game.Base)
		outputFileName := fmt.Sprintf("Dockerfile.%s", game.DockerfileExt)

		// Create the output file
		outputFile, err := os.Create(fmt.Sprintf("./game-servers/%s", outputFileName))
		if err != nil {
			fmt.Printf("Error creating Dockerfile for %s: %v\n", game.Name, err)
			continue
		}
		defer outputFile.Close()

		// Prepare the template data
		templateData := TemplateData{
			GameName:       game.Name,
			BaseDockerfile: baseDockerfile,
			Maintainer:     game.Maintainer,
			SteamAppID:     game.SteamAppID,
			AdditionalDeps: game.AdditionalDeps,
		}

		// Render the template and write to the file
		err = tmpl.Execute(outputFile, templateData)
		if err != nil {
			fmt.Printf("Error writing Dockerfile for %s: %v\n", game.Name, err)
			continue
		}

		fmt.Printf("Dockerfile for %s created at ./game-servers/%s\n", game.Name, outputFileName)
	}
}
