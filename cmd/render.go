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
	"github.com/RaphaelPour/blogctl/internal/site"

	"github.com/spf13/cobra"
)

// renderCmd represents the render command
var renderCmd = &cobra.Command{
	Use:   "render",
	Short: "Renders blog to static website",
	Long:  "Collects all posts and renders the markdown using the metadata as static website",
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := site.New(site.Options{
			BlogPath: BlogPath,
			OutPath:  OutPath,
			Force:    Force,
		})
		if err != nil {
			return err
		}

		return s.Render()
	},
}

const (
	DEFAULT_OUT_PATH = "./out/"
)

var (
	OutPath string
	Force   bool
)

func init() {
	rootCmd.AddCommand(renderCmd)

	renderCmd.Flags().StringVarP(
		&OutPath,
		"out",
		"o",
		DEFAULT_OUT_PATH,
		"Output folder.",
	)

	renderCmd.Flags().BoolVarP(
		&Force,
		"force",
		"f",
		false,
		"Overwrites an existing output folder.",
	)
}
