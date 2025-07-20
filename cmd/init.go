/*
Copyright © 2020 Raphael Pour <info@raphaelpour.de>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/RaphaelPour/blogctl/internal/config"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new blog environment",
	Long:  "",
	RunE: func(cmd *cobra.Command, args []string) error {

		if _, err := os.Stat(BlogPath); err == nil {
			return fmt.Errorf("Blog environment already exists\n")
		}

		if err := os.MkdirAll(BlogPath, os.ModePerm); err != nil {
			return fmt.Errorf("Error creating blog environment: %s\n", err)
		}

		configPath := path.Join(BlogPath, "blog.json")
		marshalledConfig, err := json.MarshalIndent(config.CreateDefaultConfig(), "", "    ")
		if err != nil {
			return fmt.Errorf("Internal error: default config not marshallable")
		}

		if err = os.WriteFile(configPath, marshalledConfig, 0777); err != nil {
			return fmt.Errorf("Error writing the default config file to %s", configPath)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
