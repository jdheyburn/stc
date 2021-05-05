package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger, _ = zap.NewDevelopment()
var fromStation, toStation string

func init() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ = config.Build()
	rootCmd.AddCommand(calcCmd)
	calcCmd.Flags().StringVarP(&fromStation, "from", "f", "", "Where the season ticket should start from")
	calcCmd.Flags().StringVarP(&toStation, "to", "t", "", "Where the season ticket should end at")
	calcCmd.MarkFlagRequired("from")
	calcCmd.MarkFlagRequired("to")
}

var calcCmd = &cobra.Command{
	Use:   "calc",
	Short: "Calculate a season ticket",
	Long:  `TBC`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("Arguments", zap.String("from", fromStation), zap.String("to", toStation))
		if err := calc(fromStation, toStation); err != nil {
			logger.Error("error running calc", zap.Error(err))
			os.Exit(1)
		}
	},
}

func calc(fromStation, toStation string) error {

}
