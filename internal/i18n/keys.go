package i18n

// This file contains constants for all translation keys
// Benefits: Type safety, autocomplete, easy refactoring, prevents typos

const (
	// Root command
	KeyRootShort = "root.short"
	KeyRootLong  = "root.long"

	// Add command
	KeyAddShort           = "add.short"
	KeyAddLong            = "add.long"
	KeyAddSuccess         = "add.success"
	KeyAddSuccessMerged   = "add.success_merged"
	KeyAddFlagHours       = "add.flag.hours"
	KeyAddFlagDescription = "add.flag.description"
	KeyAddFlagProject     = "add.flag.project"
	KeyAddFlagClient      = "add.flag.client"
	KeyAddFlagConsultant  = "add.flag.consultant"
	KeyAddFlagRate        = "add.flag.rate"
	KeyAddFlagDate        = "add.flag.date"

	// Add output labels
	KeyAddOutputDate        = "add.output.date"
	KeyAddOutputConsultant  = "add.output.consultant"
	KeyAddOutputHours       = "add.output.hours"
	KeyAddOutputRate        = "add.output.rate"
	KeyAddOutputCost        = "add.output.cost"
	KeyAddOutputProject     = "add.output.project"
	KeyAddOutputCustomer    = "add.output.customer"
	KeyAddOutputDescription = "add.output.description"

	// Get command
	KeyGetShort          = "get.short"
	KeyGetLong           = "get.long"
	KeyGetNoResults      = "get.no_results"
	KeyGetTotalHours     = "get.total_hours"
	KeyGetTotalCost      = "get.total_cost"
	KeyGetFlagConsultant = "get.flag.consultant"
	KeyGetFlagProject    = "get.flag.project"
	KeyGetFlagCustomer   = "get.flag.customer"
	KeyGetFlagMonth      = "get.flag.month"
	KeyGetFlagFromDate   = "get.flag.from_date"
	KeyGetFlagToDate     = "get.flag.to_date"
	KeyGetFlagDate       = "get.flag.date"
	KeyGetFlagToday      = "get.flag.today"
	KeyGetFlagWeek       = "get.flag.week"
	KeyGetFlagYear       = "get.flag.year"
	KeyGetFlagOutput     = "get.flag.output"
	KeyGetFlagOutputFile = "get.flag.output_file"

	// Get table headers
	KeyGetHeaderDate        = "get.header.date"
	KeyGetHeaderConsultant  = "get.header.consultant"
	KeyGetHeaderHours       = "get.header.hours"
	KeyGetHeaderRate        = "get.header.rate"
	KeyGetHeaderCost        = "get.header.cost"
	KeyGetHeaderProject     = "get.header.project"
	KeyGetHeaderCustomer    = "get.header.customer"
	KeyGetHeaderDescription = "get.header.description"

	// Export (used by get command output)
	KeyExportSuccess = "export.success"
	KeyExportTotal   = "export.total"

	// Config command
	KeyConfigShort               = "config.short"
	KeyConfigLong                = "config.long"
	KeyConfigTitle               = "config.title"
	KeyConfigNoDefaults          = "config.no_defaults"
	KeyConfigSetInstruction      = "config.set_instruction"
	KeyConfigDefaultConsultant   = "config.default_consultant"
	KeyConfigDefaultClient       = "config.default_client"
	KeyConfigDefaultProject      = "config.default_project"
	KeyConfigDefaultRate         = "config.default_rate"
	KeyConfigLanguage            = "config.language"
	KeyConfigSaved               = "config.saved"
	KeyConfigCleared             = "config.cleared"
	KeyConfigFlagConsultant      = "config.flag.consultant"
	KeyConfigFlagClient          = "config.flag.client"
	KeyConfigFlagProject         = "config.flag.project"
	KeyConfigFlagRate            = "config.flag.rate"
	KeyConfigFlagLanguage        = "config.flag.language"
	KeyConfigFlagDatabaseHost    = "config.flag.database_host"
	KeyConfigFlagDatabasePort    = "config.flag.database_port"
	KeyConfigFlagDatabaseUser    = "config.flag.database_user"
	KeyConfigFlagDatabasePass    = "config.flag.database_password"
	KeyConfigFlagDatabaseName    = "config.flag.database_name"

	// Config set subcommand
	KeyConfigSetShort = "config.set.short"
	KeyConfigSetLong  = "config.set.long"

	// Config clear subcommand
	KeyConfigClearShort = "config.clear.short"
	KeyConfigClearLong  = "config.clear.long"

	// Config database section
	KeyConfigDatabaseTitle = "config.database.title"
	KeyConfigDatabaseHost  = "config.database.host"
	KeyConfigDatabasePort  = "config.database.port"
	KeyConfigDatabaseUser  = "config.database.user"
	KeyConfigDatabaseName  = "config.database.name"

	// Error messages - general
	KeyErrReadConfig         = "error.read_config"
	KeyErrLoadConfig         = "error.load_config"
	KeyErrMustSpecifyValue   = "error.must_specify_value"
	KeyErrInvalidDateFormat  = "error.invalid_date_format"
	KeyErrHoursMustBePositive = "error.hours_must_be_positive"
	KeyErrWeekRange          = "error.week_range"
	KeyErrMonthRange         = "error.month_range"

	// Error messages - add command
	KeyErrConsultantRequired  = "error.consultant_required"
	KeyErrCustomerRequired    = "error.customer_required"
	KeyErrProjectRequired     = "error.project_required"
	KeyErrRateRequired        = "error.rate_required"
	KeyErrGetCreateConsultant = "error.get_create_consultant"
	KeyErrGetCreateCustomer   = "error.get_create_customer"
	KeyErrGetCreateProject    = "error.get_create_project"
	KeyErrCheckExistingEntry  = "error.check_existing_entry"
	KeyErrUpdateWorkLog       = "error.update_worklog"
	KeyErrSaveWorkLog         = "error.save_worklog"

	// Error messages - get command
	KeyErrFetchWorkLogs = "error.fetch_worklogs"

	// Error messages - CSV export
	KeyErrCreateOutputFile = "error.create_output_file"
	KeyErrWriteCSVHeader   = "error.write_csv_header"
	KeyErrWriteCSVRow      = "error.write_csv_row"
	KeyErrWriteTotalsRow   = "error.write_totals_row"

	// Error messages - database
	KeyErrDatabaseHostRequired     = "error.database.host_required"
	KeyErrDatabasePortRequired     = "error.database.port_required"
	KeyErrDatabaseUserRequired     = "error.database.user_required"
	KeyErrDatabasePasswordRequired = "error.database.password_required"
	KeyErrDatabaseNameRequired     = "error.database.name_required"
	KeyErrDatabaseConnect          = "error.database.connect"
	KeyErrDatabaseMigrate          = "error.database.migrate"

	// Error messages - initialization
	KeyErrInitConfig   = "error.init.config"
	KeyErrInitDatabase = "error.init.database"
	KeyErrInitI18n     = "error.init.i18n"
)
