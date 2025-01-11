package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "start",
	Short: "CLI for in-toto",
	Long:  `CLI for in-toto`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		fmt.Println("test my app")
	},
}

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Prints a hello message",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		if verbose {
			fmt.Println("This is a verbose message.")
		}
		fmt.Println("Hello, world!")
	},
}

func Execute() {
	// 添加子命令
	//rootCmd.AddCommand(helloCmd)
	//// 为 hello 命令添加标志
	//helloCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
	rootCmd.Execute()
	log.Println("start execute program")

}
