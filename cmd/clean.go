// Copyright © 2016 Devin Pastoor <devin.pastoor@gmail.com>
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

package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/dpastoor/nonmemutils/utils"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	noFolders bool
	noFiles   bool
	inverse   bool
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "clean files and folders",
	Long: `
nmu clean ^run // anything beginning with the letters run
nmu clean ^run -v // print out files and folders that will be deleted 
nmu clean ^run --noFolders // only remove matching files 
nmu clean _est_ --noFiles  // only remove matching folders  
nmu clean _est_ --noFiles --simulateClean // show what output would be if clean occured but don't actually clean 
nmu clean "run009.[^mod]" // all matching run009.<ext> but not .mod files
nmu clean "run009.(mod|lst)$" // match run009.lst and run009.mod

can also clean via the opposite of a match with inverse

nmu clean ".modt{0,1}$" --noFolders --inverse // clean all files not matching .mod or .modt
 `,
	RunE: clean,
}

func clean(cmd *cobra.Command, args []string) error {
	if debug {
		viper.Debug()
	}

	AppFs := afero.NewOsFs()

	dir, _ := filepath.Abs(".")
	dirInfo, err := afero.ReadDir(AppFs, dir)
	if err != nil {
		return fmt.Errorf("error finding dir (%s)", err)
	}
	files := utils.ListFiles(dirInfo)
	folders := utils.ListDirNames(dirInfo)
	matchedFiles := []string{}
	matchedFolders := []string{}
	for _, expr := range args {
		if !noFolders {
			var matches []string
			if !inverse {
				matches, err = utils.ListMatchesByRegex(folders, expr)
			} else {
				matches, err = utils.ListNonMatchesByRegex(folders, expr)
			}
			if err != nil {
				return fmt.Errorf("error with regex (%s), err: (%s)", expr, err)
			}
			matchedFolders = append(matchedFolders, matches...)
		}
		if !noFiles {
			var matches []string
			if !inverse {
				matches, err = utils.ListMatchesByRegex(files, expr)
			} else {

				matches, err = utils.ListNonMatchesByRegex(files, expr)
			}
			if err != nil {
				return fmt.Errorf("error with regex (%s), err: (%s)", expr, err)
			}
			matchedFiles = append(matchedFiles, matches...)
		}
	}
	if verbose {
		fmt.Println("cleaning files: ", matchedFiles)
		fmt.Println("cleaning folders: ", matchedFolders)
	}
	if preview {
		fmt.Println("would clean files: ", matchedFiles)
		fmt.Println("would clean folders: ", matchedFolders)
		return nil
	}
	if !noFolders {
		for _, f := range matchedFolders {
			AppFs.RemoveAll(filepath.Join(dir, f))
		}
	}
	if !noFiles {
		for _, f := range matchedFiles {
			AppFs.Remove(filepath.Join(dir, f))
		}
	}
	return nil
}
func init() {
	RootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&noFolders, "noFolders", false, "exclude folders during cleaning")
	cleanCmd.Flags().BoolVar(&noFiles, "noFiles", false, "exclude files during cleaning")
	cleanCmd.Flags().BoolVar(&inverse, "inverse", false, "inverse selection from the given regex match criteria")
}

