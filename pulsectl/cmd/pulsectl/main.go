package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "status":
		cmdStatus()
	case "config":
		if len(os.Args) < 3 || os.Args[2] != "validate" {
			printUsage()
			os.Exit(1)
		}
		cmdConfigValidate()
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("kullanim: pulsectl <komut>")
	fmt.Println("komutlar:")
	fmt.Println("  status            agent'in son durumunu gosterir")
	fmt.Println("  config validate   config dosyasinin gecerliligini kontrol eder")
}

func cmdStatus() {
	statusPath := "../agent/status.json"
	data, err := os.ReadFile(statusPath)
	if err != nil {
		fmt.Printf("durum dosyasi okunamadi (%s): %v\n", statusPath, err)
		fmt.Println("agent hic calismamis olabilir")
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func cmdConfigValidate() {
	configPath := "../agent/config.example.yaml"
	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("config dosyasi okunamadi (%s): %v\n", configPath, err)
		os.Exit(1)
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		fmt.Printf("config GECERSIZ: %v\n", err)
		os.Exit(1)
	}

	requiredKeys := []string{"collector", "host", "checks", "queue"}
	var missing []string
	for _, key := range requiredKeys {
		if _, ok := raw[key]; !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		fmt.Printf("config GECERSIZ: eksik alanlar: %v\n", missing)
		os.Exit(1)
	}

	fmt.Println("config GECERLI")
}
