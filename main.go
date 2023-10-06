package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// GitHubConfig represents the GitHub source configuration.
type GitHubConfig struct {
	Token       string   `form:"github-token"`
	User        string   `form:"user"`
	Username    string   `form:"username"`
	Password    string   `form:"password"`
	SSH         bool     `form:"ssh"`
	SSHKey      string   `form:"sshkey"`
	Exclude     []string `form:"exclude[]"`
	Include     []string `form:"include[]"`
	ExcludeOrgs []string `form:"excludeorgs[]"`
	IncludeOrgs []string `form:"includeorgs[]"`
	Wiki        bool     `form:"wiki"`
	Starred     bool     `form:"starred"`
	Filter      struct {
		Stars           int      `form:"filter.stars"`
		LastActivity    string   `form:"filter.lastactivity"`
		ExcludeArchived bool     `form:"filter.excludearchived"`
		Languages       []string `form:"filter.languages[]"`
		ExcludeForks    bool     `form:"filter.excludeforks"`
	}
}

var htmlTemplates = template.Must(template.New("").ParseGlob("templates/*"))

func main() {
	r := gin.Default()

	// Configure the router to use your templates.
	r.LoadHTMLGlob("templates/*")

	// Serve static files (CSS)
	r.Static("/static", "./static")

	// Define a route for the GitHub configuration form
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Define a route to handle form submissions
	r.POST("/generate", func(c *gin.Context) {
		var githubConfig GitHubConfig

		// Bind the form data to the GitHubConfig struct
		if err := c.ShouldBind(&githubConfig); err != nil {
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"Error": err.Error()})
			return
		}

		// Generate the config.yaml content based on githubConfig
		configContent := generateConfigYAML(githubConfig)

		// Save the config.yaml content to a file
		// Here, you should implement logic to save it to the desired location
		// For simplicity, we'll print it for demonstration purposes.
		fmt.Println(configContent)

		// Save the config.yaml content to a file
		configFilename := "config.yaml" // Set the desired filename
		if err := writeConfigToFile(configContent, configFilename); err != nil {
			c.HTML(http.StatusInternalServerError, "index.html", gin.H{"Error": err.Error()})
			return
		}

		c.String(http.StatusOK, "Config.yaml generated and saved successfully")

	})

	// Start the server
	r.Run(":8080")
}

// generateConfigYAML generates a config.yaml string based on the provided GitHubConfig.
func generateConfigYAML(config GitHubConfig) string {
	// You should implement logic here to generate the config.yaml content
	// based on the provided GitHubConfig struct.
	// This is a simplified example.
	return fmt.Sprintf(`github:
  - token: %s
      user: %s
      username: %s
      password: %s
      ssh: %v
      sshkey: %s
      exclude:
  %s
      include:
%s
      excludeorgs:
%s
      includeorgs:
%s
      wiki: %v
      starred: %v
      filter:
        stars: %d
        lastactivity: %s
        excludearchived: %v
        languages:
%s
        excludeforks: %v`,
		config.Token, config.User, config.Username, config.Password, config.SSH, config.SSHKey,
		strings.Join(config.Exclude, "\n"), strings.Join(config.Include, "\n"),
		strings.Join(config.ExcludeOrgs, "\n"), strings.Join(config.IncludeOrgs, "\n"),
		config.Wiki, config.Starred, config.Filter.Stars, config.Filter.LastActivity,
		config.Filter.ExcludeArchived, strings.Join(config.Filter.Languages, "\n"), config.Filter.ExcludeForks)
}

func writeConfigToFile(configContent string, filename string) error {
	// Open the file for writing, create it if it doesn't exist, and truncate it if it does.
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the YAML content to the file
	_, err = file.WriteString(configContent)
	if err != nil {
		return err
	}

	return nil
}
