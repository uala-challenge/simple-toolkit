package viper

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/app_profile"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/file_utils"
)

var _ Service = (*service)(nil)

func NewService(logger *logrus.Logger) Service {
	once.Do(func() {
		instance = &service{
			propertyFiles: getPropertyFiles(logger),
			path:          getConfigPath(logger),
			log:           logger,
		}
	})
	return instance
}

func (s *service) Apply() (Config, error) {
	if err := s.validateRequiredFiles(); err != nil {
		s.log.Error("Error validando archivos de configuración: ", err)
		return Config{}, err
	}

	mergedConfig, err := s.loadAndMergeConfigs()
	if err != nil {
		s.log.Error("Error cargando configuración: ", err)
		return Config{}, fmt.Errorf("error loading configuration - %w", err)
	}

	s.log.Info("Configuración cargada correctamente")
	return s.mapConfigToStruct(mergedConfig)
}

func (s *service) validateRequiredFiles() error {
	files, err := file_utils.ListFiles(s.path)
	if err != nil {
		s.log.Errorf("Error listando archivos en %s: %v", s.path, err)
		return err
	}

	missingFiles := getMissingFiles(s.propertyFiles, files)
	if len(missingFiles) > 0 {
		s.log.Errorf("Archivos de configuración faltantes: %v", missingFiles)
		return fmt.Errorf("faltan archivos de configuración: %v", missingFiles)
	}

	s.log.Debug("Todos los archivos de configuración requeridos están presentes")
	return nil
}

func (s *service) loadAndMergeConfigs() (*viper.Viper, error) {
	baseConfig, err := loadConfig(s.path, "application", s.log)
	if err != nil {
		return nil, err
	}

	envConfig, err := loadConfig(s.path, s.getPropertyFileName(), s.log)
	if err != nil {
		return nil, err
	}

	if err := baseConfig.MergeConfigMap(envConfig.AllSettings()); err != nil {
		return nil, fmt.Errorf("failed to merge configurations: %w", err)
	}

	s.log.Debug("Archivos de configuración combinados correctamente")
	return baseConfig, nil
}

func (s *service) mapConfigToStruct(v *viper.Viper) (Config, error) {
	configMap, err := unmarshalConfig(v)
	if err != nil {
		s.log.Error("Error al deserializar configuración", err)
		return Config{}, err
	}

	processedConfig, err := processConfigValues(configMap)
	if err != nil {
		s.log.Error("Error al procesar configuración", err)
		return Config{}, err
	}

	s.log.Debug("Configuración mapeada correctamente a la estructura de datos")
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

func getPropertyFiles(logger *logrus.Logger) []string {
	requiredFiles := []string{"application.yaml"}
	scopeFile := fmt.Sprintf("application-%s.yaml", app_profile.GetScopeValue())
	availableFiles, _ := file_utils.ListFiles(getConfigPath(logger))

	if contains(availableFiles, scopeFile) {
		requiredFiles = append(requiredFiles, scopeFile)
	}

	logger.Debugf("Archivos requeridos: %v", requiredFiles)
	return requiredFiles
}

func getConfigPath(logger *logrus.Logger) string {
	if path := os.Getenv("CONF_DIR"); path != "" {
		logger.Debugf("Usando CONF_DIR: %s", path)
		return path
	}

	if app_profile.IsLocalProfile() {
		logger.Debug("Usando perfil local para configuración")
		return "kit/config"
	}

	logger.Debug("Usando configuración por defecto en /app/kit/config")
	return "/app/kit/config"
}

func loadConfig(path, filename string, logger *logrus.Logger) (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(filename)
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		logger.Errorf("Error leyendo archivo de configuración %s/%s: %v", path, filename, err)
		return nil, err
	}

	if v.GetBool("enable_config_watch") {
		watchConfig(v, logger)
	}

	logger.Infof("Archivo de configuración %s cargado correctamente", filename)
	return v, nil
}

func watchConfig(v *viper.Viper, logger *logrus.Logger) {
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		logger.Warnf("Archivo de configuración cambiado: %s", e.Name)
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
