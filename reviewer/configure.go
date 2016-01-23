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
	"github.com/spf13/viper"
	"log"
	"fmt"
	"errors"
)

func Configure() {

	err := CheckFile()
	if err != nil {
		log.Fatal(err)
	}

	err = CheckRepositories()
	if err != nil {
		log.Fatal(err)
	}

	resp, err2 := CheckRepositoriesData()
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Printf("%s",resp)

}

// ConfigFileUsed contains filename used for configuration.
var ConfigFileUsed = viper.ConfigFileUsed
var Sub = viper.Sub
var IsSet = viper.IsSet
var GetString = viper.GetString

func CheckFile() error {
	configFile := ConfigFileUsed()
	if configFile == "" {
		return errors.New("Config file not defined or empty")
	}
	return nil
}

func CheckRepositories() error {
	if(!IsSet("repositories")){
		return errors.New("Repositories not set")
	}
	return nil
}

func CheckRepositoriesData() (s string, err error){
	repo := Sub("repositories")
	keys := repo.AllKeys()
	var response = ""

	for _,v := range keys{
		username := GetString("repositories."+v+".username")
		status := GetString("repositories."+v+".status")
		required := GetString("repositories."+v+".required")
		if(username == "" || status == "" || required == "") {
			return "", errors.New("Fields not set")
		}
		mode := "ENABLED"
		if status == "false" {
			mode = "DISABLED"
		}
		response += fmt.Sprintf("- %s / %s %s +1:%s\n", username, v, mode, required )
	}
	return response, nil
}