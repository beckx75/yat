/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"beckx.online/yat/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cw := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	cw.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[ %-6s]", i))
	}
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(cw)

	log.Info().Msg("yat started...")
	cmd.Execute()
}
