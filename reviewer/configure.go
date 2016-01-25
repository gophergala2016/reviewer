// Copyright © 2016 See CONTRIBUTORS <alvaro.cabanas@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reviewer

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

// ConfigRepositoriesChecker is an interface for checking Viper's keys or getting their values
type ConfigRepositoriesChecker interface {
	AllKeys() []string
	GetString(key string) string
}

// Config is a wrapper around Viper
type Config struct {
	config *viper.Viper
}

// NewConfig is the constructor for Config.
func NewConfig(viper *viper.Viper) *Config {
	return &Config{
		config: viper,
	}
}

// AllKeys returns all keys regardless where they are set
func (c *Config) AllKeys() []string {
	return c.config.AllKeys()
}

// GetString returns the value associated with the key as a string
func (c *Config) GetString(key string) string {
	return c.config.GetString(key)
}

// GetBool returns the value associated with the key as a bool
func (c *Config) GetBool(key string) bool {
	return c.config.GetBool(key)
}

// GetInt returns the value associated with the key as an integer
func (c *Config) GetInt(key string) int {
	return c.config.GetInt(key)
}

// GetStringSlice returns the value associated with the key as an slice of strings
func (c *Config) GetStringSlice(key string) []string {
	return c.config.GetStringSlice(key)
}

// IsSet contains the function used to check if key is set
var IsSet = viper.IsSet

// ConfigFileUsed contains the function used to get Configuration file used
var ConfigFileUsed = viper.ConfigFileUsed

// Configure runs all the checks needed for the config file to be set up correctly and prints the repositories available
func Configure() {

	err := CheckFile()
	if err != nil {
		log.Fatal(err)
	}

	err = CheckRepositories()
	if err != nil {
		log.Fatal(err)
	}

	config := NewConfig(viper.Sub("repositories"))
	resp, err2 := CheckRepositoriesData(config)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Printf("%s", resp)

}

// CheckFile checks if the configuration file exists and is not empty
func CheckFile() error {
	configFile := ConfigFileUsed()
	if configFile == "" {
		return errors.New("Config file not defined or empty")
	}
	return nil
}

// CheckRepositories checks if the configuration file has repositories
func CheckRepositories() error {
	if !IsSet("repositories") {
		return errors.New("Repositories not set")
	}
	return nil
}

// CheckRepositoriesData checks if all the repositories have all the parameters set
func CheckRepositoriesData(config ConfigRepositoriesChecker) (s string, err error) {

	keys := config.AllKeys()
	var response = ""
	for _, v := range keys {
		username := config.GetString(v + ".username")
		status := config.GetString(v + ".status")
		required := config.GetString(v + ".required")
		if username == "" || status == "" || required == "" {
			return "", errors.New("Fields not set")
		}
		mode := "ENABLED"
		if status == "false" {
			mode = "DISABLED"
		}
		response += fmt.Sprintf("- %s / %s %s +1:%s\n", username, v, mode, required)
	}
	return response, nil
}
