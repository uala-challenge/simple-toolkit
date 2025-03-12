package viper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/app_profile"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/file_utils"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var _ Service = (*service)(nil)

func NewService() Service {
	once.Do(func() {
		instance = &service{
			propertyFiles: getPropertyFiles(),
			path:          getConfigPath(),
		}
	})
	return instance
}

func (s *service) Apply() (Config, error) {
	if err := s.validateRequiredFiles(); err != nil {
		return Config{}, err
	}

	mergedConfig, err := s.loadAndMergeConfigs()
	if err != nil {
		return Config{}, fmt.Errorf("error loading configuration - %w", err)
	}

	return s.mapConfigToStruct(mergedConfig)
}

func (s *service) validateRequiredFiles() error {
	files, err := file_utils.ListFiles(s.path)
	if err != nil {
		return err
	}

	missingFiles := getMissingFiles(s.propertyFiles, files)
	if len(missingFiles) > 0 {
		return fmt.Errorf("missing required files: %v", missingFiles)
	}
	return nil
}

func (s *service) loadAndMergeConfigs() (*viper.Viper, error) {
	baseConfig, err := loadConfig(s.path, "application")
	if err != nil {
		return nil, err
	}

	envConfig, err := loadConfig(s.path, s.getPropertyFileName())
	if err != nil {
		return nil, err
	}

	if err := baseConfig.MergeConfigMap(envConfig.AllSettings()); err != nil {
		return nil, fmt.Errorf("failed to merge configurations: %w", err)
	}

	return baseConfig, nil
}

func (s *service) mapConfigToStruct(v *viper.Viper) (Config, error) {
	configMap, err := unmarshalConfig(v)
	if err != nil {
		return Config{}, err
	}

	processedConfig, err := processConfigValues(configMap)
	if err != nil {
		return Config{}, err
	}

	return decodeToStruct(processedConfig)
}

func (s *service) getPropertyFileName() string {
	scopeFile := fmt.Sprintf("application-%s", app_profile.GetScopeValue())
	profileFile := fmt.Sprintf("application-%s", app_profile.GetProfileByScope())
	files, _ := file_utils.ListFiles(s.path)

	if contains(files, scopeFile) {
		return scopeFile
	}
	return profileFile
}

func getPropertyFiles() []string {
	requiredFiles := []string{
		"application.yaml",
		"application-local.yaml",
		"application-prod.yaml",
	}
	scopeFile := fmt.Sprintf("application-%s.yaml", app_profile.GetScopeValue())

	if scopeFile != "application-local.yaml" && scopeFile != "application-prod.yaml" {
		requiredFiles = append(requiredFiles, scopeFile)
	}

	return requiredFiles
}

func getConfigPath() string {
	if path := os.Getenv("CONF_DIR"); path != "" {
		return path
	}
	return "kit/config"
}

func loadConfig(path, filename string) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(filename)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if v.GetBool("enable_config_watch") {
		watchConfig(v)
	}

	return v, nil
}

func watchConfig(v *viper.Viper) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})
}

func unmarshalConfig(v *viper.Viper) (map[string]interface{}, error) {
	var configMap map[string]interface{}
	if err := v.Unmarshal(&configMap); err != nil {
		return nil, err
	}
	return configMap, nil
}

func processConfigValues(config map[string]interface{}) (map[string]interface{}, error) {
	for key, value := range config {
		switch v := value.(type) {
		case map[string]interface{}:
			processed, err := processConfigValues(v)
			if err != nil {
				return nil, err
			}
			config[key] = processed
		case []interface{}:
			processed, err := processSliceValues(v)
			if err != nil {
				return nil, err
			}
			config[key] = processed
		case string:
			config[key] = resolveEnvValue(v)
		}
	}
	return config, nil
}

func processSliceValues(slice []interface{}) ([]interface{}, error) {
	for i, elem := range slice {
		switch v := elem.(type) {
		case map[string]interface{}:
			processed, err := processConfigValues(v)
			if err != nil {
				return nil, err
			}
			slice[i] = processed
		case []interface{}:
			processed, err := processSliceValues(v)
			if err != nil {
				return nil, err
			}
			slice[i] = processed
		case string:
			slice[i] = resolveEnvValue(v)
		}
	}
	return slice, nil
}

func resolveEnvValue(value string) string {
	if strings.HasPrefix(value, "${") && strings.HasSuffix(value, "}") {
		trimmed := strings.Trim(value, "${}")
		parts := strings.SplitN(trimmed, ":-", 2)
		envValue := os.Getenv(parts[0])

		if envValue != "" {
			return envValue
		}
		if len(parts) > 1 {
			return parts[1]
		}
	}
	return value
}

func decodeToStruct(config map[string]interface{}) (Config, error) {
	var result Config
	if err := mapstructure.Decode(config, &result); err != nil {
		return Config{}, err
	}
	return result, nil
}

func getMissingFiles(required, available []string) []string {
	var missing []string
	for _, file := range required {
		if !contains(available, file) {
			missing = append(missing, file)
		}
	}
	return missing
}

func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
