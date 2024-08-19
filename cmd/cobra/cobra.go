package cobra

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// default command, executed when no command passed to application
var rootCmd = &cobra.Command{
	Use:                   "docker",
	Long:                  "\nDocker is a tool for managing containers",
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help() //just print help
	},
}

var ps = &cobra.Command{
	Use:     "ps",
	Short:   "List containers",
	Aliases: []string{"list"},
	Run: func(cmd *cobra.Command, args []string) {
		f := cmd.Flag("global")
		//check if global was provided
		if f.Changed {
			log.Printf("global flag: %v", f.Value)
		}
		log.Println("printing containers")
	},
}

// RunCobraExample tiny example how to use cobra package to create commands
func RunCobraExample() {
	rootCmd.AddCommand(ps)
	rootCmd.PersistentFlags().StringP("global", "g", "", "Global flag available to all commands")

	//flag available to ps command only
	ps.PersistentFlags().StringP("format", "f", "", "Format the output")

	//can also bind flags to viper so can be accessed outside cmd scope
	err := viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatalf("error binding flags: %v", err)
	}
	cobra.CheckErr(rootCmd.Execute())
}
