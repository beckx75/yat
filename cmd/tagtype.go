/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"beckx.online/yat/yat"
	// "beckx.online/yat/ataglib"
	"github.com/spf13/cobra"
)

// tagtypeCmd represents the tagtype command
var tagtypeCmd = &cobra.Command{
	Use:   "tagtype",
	Short: "looks for the Tag-Type in Audiofiles",
	Long: `to be defined...`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
		yd, err := yat.NewYatData(args)
		if err != nil {
			panic(err)
		}
		fmt.Println("Got files: ", len(yd.Files))
		err = yd.ReadAudioMetadata(true)
		if err != nil {
			panic(err)
		}
		ttOccurence := make(map[string][]string)
		for _, amd := range yd.AudioMetadatas {
			t := amd.TagVersion
			_, ok := ttOccurence[t]
			if ok {
				ttOccurence[t] = append(ttOccurence[t], amd.Filepath)
			} else {
				ttOccurence[t] = []string{amd.Filepath}
			}
		}
		for tt, fps := range ttOccurence {
			fmt.Printf("%s: %d\n", tt, len(fps))
		}
	},
}

func init() {
	rootCmd.AddCommand(tagtypeCmd)
	
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagtypeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagtypeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
