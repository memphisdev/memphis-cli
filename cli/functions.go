package cli

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

//go:embed function_templates
var templatesFS embed.FS

// //go:embed function_templates/nodejs/*
// var templatesNodeJs embed.FS
// //go:embed function_templates/go/*
// var templatesNodeGo embed.FS

// //go:embed function_templates/python/*
// var templatesNodePython embed.FS

type Tag struct {
	Tag string `yaml:"tag"`
}

type Input struct {
	Name string `yaml:"name"`
}

type MemphisYaml struct {
	FunctionName string  `yaml:"function_name"`
	Description  string  `yaml:"description"`
	Runtime      string  `yaml:"runtime"`
	Handler      string  `yaml:"handler,omitempty"`
	Dependencies string  `yaml:"dependencies,omitempty"`
	Tags         []Tag   `yaml:"tags"`
	Inputs       []Input `yaml:"inputs"`
}

func copyFile(src, dst string) error {
	sourceFile, err := templatesFS.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if strings.HasSuffix(dst, ".mod1") || strings.HasSuffix(dst, ".go1") || strings.HasSuffix(dst, ".sum1") {
		dst = dst[:len(dst)-1]
	}
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func copyDir(srcDir, dstDir string) error {
	files, err := templatesFS.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcFilePath := filepath.Join(srcDir, file.Name())
		dstFilePath := filepath.Join(dstDir, file.Name())

		if file.IsDir() {
			continue
		}

		if err := copyFile(srcFilePath, dstFilePath); err != nil {
			return err
		}
	}
	return nil
}

var funcInitCmd = &cobra.Command{
	Use:     "init [function name]",
	Short:   "Generates a memphis function template",
	Args:    cobra.ExactArgs(1),
	Example: "func init myFunc --lang nodejs",
	Run: func(cmd *cobra.Command, args []string) {
		runtime, _ := cmd.Flags().GetString("lang")
		if runtime != "go" && runtime != "nodejs" && runtime != "python" {
			fmt.Println("Unsupported language: " + runtime + ". Supported languages are: nodejs, python, go")
			return
		}

		fmt.Println("Generating function template for " + runtime)
		functionName := args[0]
		if strings.Contains(functionName, "/") || strings.Contains(functionName, " ") {
			fmt.Println("Function name cannot contain '/' or spaces")
			return
		}

		err := os.Mkdir(functionName, 0755)
		if err != nil && !os.IsExist(err) {
			fmt.Println(err)
			return
		}

		err = copyDir("function_templates/"+runtime, functionName)
		if err != nil {
			fmt.Println(err)
			return
		}

		config := MemphisYaml{
			FunctionName: functionName,
			Description:  "",
			Runtime:      runtime,
			Tags: []Tag{
				{Tag: "json"},
			},
			Inputs: []Input{
				{Name: "field_to_ingest"},
			},
		}
		switch runtime {
		case "nodejs":
			config.Handler = "index.handler"
		case "python":
			config.Handler = "main.handler"
			config.Dependencies = "requirements.txt"
		}

		yamlData, err := yaml.Marshal(&config)
		if err != nil {
			fmt.Printf("Error marshalling YAML: %v\n", err)
			return
		}

		filePath := functionName + "/memphis.yaml"
		err = os.WriteFile(filePath, yamlData, 0644)
		if err != nil {
			fmt.Printf("Error writing YAML file: %v\n", err)
			return
		}

	},
}

var funcCmd = &cobra.Command{
	Use:     "function",
	Aliases: []string{"func"},
}

func init() {
	funcInitCmd.Flags().String("lang", "nodejs", "the desired language for the function")
	funcCmd.AddCommand(funcInitCmd)
	rootCmd.AddCommand(funcCmd)
}
