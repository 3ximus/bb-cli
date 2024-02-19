package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var VariablesCmd = &cobra.Command{
	Use:     "variables",
	Short:   "Manage pipeline variables",
	Long:    "Manage pipeline variables. If variable is secured only *** is displayed",
	Aliases: []string{"var"},
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		variables := <-api.GetPipelineVariables(repo)

		setVars, _ := cmd.Flags().GetStringSlice("set")
		if len(setVars) > 0 {
			for _, v := range setVars {
				keyVal := strings.Split(v, "=")
				if len(keyVal) != 2 {
					cobra.CheckErr(fmt.Sprintf("Variable \"%s\" must be in the format \"KEY=VALUE\"", v))
				}
				updated := false
				for _, ev := range variables {
					if ev.Key == keyVal[0] {
						updatedVar := <-api.UpdatePipelineVariable(repo, ev.UUID, keyVal[0], keyVal[1])
						util.Printf("\033[1;34mUpdated\033[m \"%s=%s\"\n", updatedVar.Key, updatedVar.Value)
						updated = true
						break
					}
				}
				if !updated {
					createdVar := <-api.CreatePipelineVariable(repo, keyVal[0], keyVal[1])
					util.Printf("\033[1;32mCreated\033[m \"%s=%s\"\n", createdVar.Key, createdVar.Value)
				}
			}
		}

		deleteVars, _ := cmd.Flags().GetStringSlice("delete")
		if len(deleteVars) > 0 {
			for _, toDelete := range deleteVars {
				for _, ev := range variables {
					if ev.Key == toDelete {
						api.DeletePipelineVariable(repo, ev.UUID)
						util.Printf("\033[1;31mDeleted\033[m \"%s\"\n", ev.Key)
						break
					}
				}
			}

		}

		if len(setVars) == 0 && len(deleteVars) == 0 {
			for _, variable := range variables {
				if variable.Secured {
					util.Printf("%s = \033[37m***\033[m", variable.Key)
				} else {
					util.Printf("%s = \033[37m%s\033[m", variable.Key, variable.Value)
				}
				fmt.Println()
			}
		}
	},
}

func init() {
	// TODO make it possible to set secure variables with a special flag. Maybe A=B
	// TODO Commas are used to separate variables within the same flag so "-s A=B,X" will not set variable A to the value B,X
	// So we must not use StringSlice and use String array instead ?
	// TODO It cannot also handle = signs
	VariablesCmd.Flags().StringSliceP("set", "s", []string{}, `set one or multiple variables. If the variable doesn't exist one is created.
	Variables must be in the format KEY=VALUE`)
	VariablesCmd.Flags().StringSliceP("delete", "d", []string{}, "delete one or multiple variables")
}
