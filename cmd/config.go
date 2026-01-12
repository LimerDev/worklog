package cmd

import (
	"fmt"

	"github.com/LimerDev/worklog/internal/config"
	"github.com/LimerDev/worklog/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	configConsultant string
	configClient     string
	configProject    string
	configRate       float64
	configLanguage   string
	configDBHost     string
	configDBPort     string
	configDBUser     string
	configDBPassword string
	configDBName     string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "",
	Long:  "",
	RunE:  runConfigShow,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "",
	Long:  "",
	RunE:  runConfigSet,
}

var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "",
	Long:  "",
	RunE:  runConfigClear,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configClearCmd)

	configSetCmd.Flags().StringVarP(&configConsultant, "consultant", "n", "", "")
	configSetCmd.Flags().StringVarP(&configClient, "client", "c", "", "")
	configSetCmd.Flags().StringVarP(&configProject, "project", "p", "", "")
	configSetCmd.Flags().Float64VarP(&configRate, "rate", "r", 0, "")
	configSetCmd.Flags().StringVarP(&configLanguage, "language", "l", "", "")
	configSetCmd.Flags().StringVar(&configDBHost, "db-host", "", "")
	configSetCmd.Flags().StringVar(&configDBPort, "db-port", "", "")
	configSetCmd.Flags().StringVar(&configDBUser, "db-user", "", "")
	configSetCmd.Flags().StringVar(&configDBPassword, "db-password", "", "")
	configSetCmd.Flags().StringVar(&configDBName, "db-name", "", "")
}

func localizeConfigCommand() {
	configCmd.Short = i18n.T(i18n.KeyConfigShort)
	configCmd.Long = i18n.T(i18n.KeyConfigLong)

	configSetCmd.Short = i18n.T(i18n.KeyConfigSetShort)
	configSetCmd.Long = i18n.T(i18n.KeyConfigSetLong)

	configClearCmd.Short = i18n.T(i18n.KeyConfigClearShort)
	configClearCmd.Long = i18n.T(i18n.KeyConfigClearLong)

	configSetCmd.Flags().Lookup("consultant").Usage = i18n.T(i18n.KeyConfigFlagConsultant)
	configSetCmd.Flags().Lookup("client").Usage = i18n.T(i18n.KeyConfigFlagClient)
	configSetCmd.Flags().Lookup("project").Usage = i18n.T(i18n.KeyConfigFlagProject)
	configSetCmd.Flags().Lookup("rate").Usage = i18n.T(i18n.KeyConfigFlagRate)
	configSetCmd.Flags().Lookup("language").Usage = i18n.T(i18n.KeyConfigFlagLanguage)
	configSetCmd.Flags().Lookup("db-host").Usage = i18n.T(i18n.KeyConfigFlagDatabaseHost)
	configSetCmd.Flags().Lookup("db-port").Usage = i18n.T(i18n.KeyConfigFlagDatabasePort)
	configSetCmd.Flags().Lookup("db-user").Usage = i18n.T(i18n.KeyConfigFlagDatabaseUser)
	configSetCmd.Flags().Lookup("db-password").Usage = i18n.T(i18n.KeyConfigFlagDatabasePass)
	configSetCmd.Flags().Lookup("db-name").Usage = i18n.T(i18n.KeyConfigFlagDatabaseName)
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("%s: %w", i18n.T(i18n.KeyErrLoadConfig), err)
	}

	fmt.Print(i18n.T(i18n.KeyConfigTitle))

	if cfg.DefaultConsultant == "" && cfg.DefaultClient == "" && cfg.DefaultProject == "" && cfg.DefaultRate == 0 {
		fmt.Println(i18n.T(i18n.KeyConfigNoDefaults))
		fmt.Println(i18n.T(i18n.KeyConfigSetInstruction))
		return nil
	}

	if cfg.DefaultConsultant != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDefaultConsultant)+"\n", cfg.DefaultConsultant)
	}
	if cfg.DefaultClient != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDefaultClient)+"\n", cfg.DefaultClient)
	}
	if cfg.DefaultProject != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDefaultProject)+"\n", cfg.DefaultProject)
	}
	if cfg.DefaultRate > 0 {
		fmt.Printf(i18n.T(i18n.KeyConfigDefaultRate)+"\n", cfg.DefaultRate)
	}
	if cfg.Language != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigLanguage)+"\n", cfg.Language)
	}

	fmt.Print(i18n.T(i18n.KeyConfigDatabaseTitle))
	if cfg.Database.Host != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDatabaseHost)+"\n", cfg.Database.Host)
	}
	if cfg.Database.Port != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDatabasePort)+"\n", cfg.Database.Port)
	}
	if cfg.Database.User != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDatabaseUser)+"\n", cfg.Database.User)
	}
	if cfg.Database.Name != "" {
		fmt.Printf(i18n.T(i18n.KeyConfigDatabaseName)+"\n", cfg.Database.Name)
	}

	return nil
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	if configConsultant == "" && configClient == "" && configProject == "" && configRate == 0 && configLanguage == "" &&
		configDBHost == "" && configDBPort == "" && configDBUser == "" && configDBPassword == "" && configDBName == "" {
		return fmt.Errorf(i18n.T(i18n.KeyErrMustSpecifyValue))
	}

	if err := config.SaveDefaults(configConsultant, configClient, configProject, configRate, configLanguage); err != nil {
		return err
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	fmt.Println(i18n.T(i18n.KeyConfigSaved))
	fmt.Printf(i18n.T(i18n.KeyConfigDefaultConsultant)+"\n", cfg.DefaultConsultant)
	fmt.Printf(i18n.T(i18n.KeyConfigDefaultClient)+"\n", cfg.DefaultClient)
	fmt.Printf(i18n.T(i18n.KeyConfigDefaultProject)+"\n", cfg.DefaultProject)
	fmt.Printf(i18n.T(i18n.KeyConfigDefaultRate)+"\n", cfg.DefaultRate)

	return nil
}

func runConfigClear(cmd *cobra.Command, args []string) error {
	if err := config.ClearDefaults(); err != nil {
		return err
	}

	fmt.Println(i18n.T(i18n.KeyConfigCleared))

	return nil
}
