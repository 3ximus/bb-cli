package pipeline

import (
	"bb/api"
	"bb/util"
	"fmt"
	"regexp"

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

		setVars, _ := cmd.Flags().GetStringArray("set")
		upsertVariables(repo, setVars, variables, false)
		setSecureVars, _ := cmd.Flags().GetStringArray("set-secure")
		upsertVariables(repo, setSecureVars, variables, true)

		deleteVars, _ := cmd.Flags().GetStringArray("delete")
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

		if len(setVars) == 0 && len(deleteVars) == 0 && len(setSecureVars) == 0 {
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
	VariablesCmd.Flags().StringArrayP("set", "s", []string{}, `set one or multiple variables. If the variable doesn't exist one is created.
	Variables must be in the format KEY=VALUE`)
	VariablesCmd.Flags().StringArrayP("set-secure", "S", []string{}, `set one or multiple secure variables. If the variable doesn't exist one is created.
	Variables must be in the format KEY=VALUE`)
	VariablesCmd.Flags().StringArrayP("delete", "d", []string{}, "delete one or multiple variables")
}

func upsertVariables(repo string, setVars []string, variables []api.EnvironmentVariable, secure bool) {
	varRegex := regexp.MustCompile(`([^=]+)=(.*)`)
	if len(setVars) > 0 {
		for _, v := range setVars {
			keyVal := varRegex.FindStringSubmatch(v)
			if len(keyVal) != 3 {
				cobra.CheckErr(fmt.Sprintf("Variable \"%s\" must be in the format \"KEY=VALUE\"", v))
			}
			updated := false
			for _, ev := range variables {
				if ev.Key == keyVal[1] {
					updatedVar := <-api.UpdatePipelineVariable(repo, ev.UUID, keyVal[1], keyVal[2], secure)
					util.Printf("\033[1;34mUpdated\033[m \"%s=%s\"\n", updatedVar.Key, updatedVar.Value)
					updated = true
					break
				}
			}
			if !updated {
				createdVar := <-api.CreatePipelineVariable(repo, keyVal[1], keyVal[2], secure)
				util.Printf("\033[1;32mCreated\033[m \"%s=%s\"\n", createdVar.Key, createdVar.Value)
			}
		}
	}
}
