package environment

import (
	"bb/api"
	"bb/util"
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var VariablesCmd = &cobra.Command{
	Use:     "variables [environment]",
	Short:   "Manage variables for specific environment",
	Long:    "Manage variables for specific environment. If variable is secured only *** is displayed",
	Aliases: []string{"var"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := viper.GetString("repo")
		envName := args[0]
		env := api.GetEnvironment(repo, envName)
		if env == nil {
			cobra.CheckErr(fmt.Sprintf("Environment \"%s\" not found", envName))
		}

		var variables []api.EnvironmentVariable
		for variable := range api.GetEnvironmentVariables(repo, envName) {
			variables = append(variables, variable)
		}

		setVars, _ := cmd.Flags().GetStringArray("set")
		upsertEnvironmentVariables(repo, env.UUID, setVars, variables, false)
		setSecureVars, _ := cmd.Flags().GetStringArray("set-secure")
		upsertEnvironmentVariables(repo, env.UUID, setSecureVars, variables, true)

		deleteVars, _ := cmd.Flags().GetStringArray("delete")
		if len(deleteVars) > 0 {
			for _, toDelete := range deleteVars {
				for _, ev := range variables {
					if ev.Key == toDelete {
						api.DeleteEnvironmentVariable(repo, env.UUID, ev.UUID)
						util.Printf("\033[1;31mDeleted\033[m \"%s\"\n", ev.Key)
						break
					}
				}
			}
		}

		if len(setVars) == 0 && len(deleteVars) == 0 && len(setSecureVars) == 0 {
			for _, variable := range variables {
				if variable.Secured {
					util.Printf("%s=\033[37m***\033[m", variable.Key)
				} else {
					util.Printf("%s=\033[37m%s\033[m", variable.Key, variable.Value)
				}
				fmt.Println()
			}
		}
	},
}

func init() {
	VariablesCmd.Flags().StringArrayP("set", "s", []string{}, `set one or multiple variables. If the variable doesn't exist one is created.
	Variables must be in the format KEY=VALUE. If '-' is provided, variables are read from stdin.`)
	VariablesCmd.Flags().StringArrayP("set-secure", "S", []string{}, `set one or multiple secure variables. If the variable doesn't exist one is created.
	Variables must be in the format KEY=VALUE. If '-' is provided, variables are read from stdin.`)
	VariablesCmd.Flags().StringArrayP("delete", "d", []string{}, "delete one or multiple variables")
}

func upsertEnvironmentVariables(repo string, envUUID string, setVars []string, variables []api.EnvironmentVariable, secure bool) {
	varRegex := regexp.MustCompile(`([^=]+)=(.*)`)

	process := func(v string) {
		keyVal := varRegex.FindStringSubmatch(v)
		if len(keyVal) != 3 {
			cobra.CheckErr(fmt.Sprintf("Variable \"%s\" must be in the format \"KEY=VALUE\"", v))
		}
		updated := false
		for _, ev := range variables {
			if ev.Key == keyVal[1] {
				updatedVar := <-api.UpdateEnvironmentVariable(repo, envUUID, ev.UUID, keyVal[1], keyVal[2], secure)
				util.Printf("\033[1;34mUpdated\033[m \"%s=%s\"\n", updatedVar.Key, updatedVar.Value)
				updated = true
				break
			}
		}
		if !updated {
			createdVar := <-api.CreateEnvironmentVariable(repo, envUUID, keyVal[1], keyVal[2], secure)
			util.Printf("\033[1;32mCreated\033[m \"%s=%s\"\n", createdVar.Key, createdVar.Value)
		}
	}

	if len(setVars) > 0 {
		if len(setVars) == 1 && setVars[0] == "-" {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := scanner.Text()
				if line != "" {
					process(line)
				}
			}
		} else {
			for _, v := range setVars {
				process(v)
			}
		}
	}
}
