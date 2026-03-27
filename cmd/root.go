package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"mysql-sync-tool/pkg/database"
	"mysql-sync-tool/pkg/sqlfile"

	"github.com/spf13/cobra"
)

var (
	sqlFile        string
	host           string
	port           int
	user           string
	password       string
	databaseName   string
	charset        string
	timeout        int
	noTransaction  bool
)

var rootCmd = &cobra.Command{
	Use:   "mysql-sync-tool",
	Short: "MySQL data synchronization tool",
	Long: `A lightweight MySQL data synchronization tool that executes SQL files
to synchronize data between databases. Supports cross-platform compilation
for Windows (x86/x64), Linux (arm64/x86/loongarch64).`,
	Run: func(cmd *cobra.Command, args []string) {
		if sqlFile == "" {
			fmt.Println("Error: --sql-file is required")
			os.Exit(1)
		}

		if err := sqlfile.FileExists(sqlFile); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		var db *sql.DB
		var err error
		
		if databaseName == "" {
			fmt.Println("Warning: No database specified, connecting to MySQL server...")
			db, err = database.Connect(host, port, user, password, "", charset, timeout)
			if err != nil {
				fmt.Printf("Error connecting to MySQL server: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()
			fmt.Println("Connected to MySQL server successfully")
		} else {
			db, err = database.ConnectOrCreate(host, port, user, password, databaseName, charset, timeout)
			if err != nil {
				fmt.Printf("Error connecting to database: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()
			fmt.Println("Connected to database successfully")
		}

		statements, err := sqlfile.ParseSQLFile(sqlFile)
		if err != nil {
			fmt.Printf("Error reading SQL file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Found %d SQL statements\n", len(statements))
		fmt.Println("Optimizing INSERT statements...")
		statements = sqlfile.OptimizeStatements(statements)
		fmt.Printf("Optimized to %d statements\n", len(statements))

		err = sqlfile.ExecuteStatements(db, statements, !noTransaction)
		if err != nil {
			fmt.Printf("Error executing SQL: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("SQL file executed successfully")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&sqlFile, "sql-file", "f", "", "SQL file path (required)")
	rootCmd.Flags().StringVar(&host, "host", "127.0.0.1", "MySQL host address")
	rootCmd.Flags().IntVar(&port, "port", 3306, "MySQL port")
	rootCmd.Flags().StringVar(&user, "user", "root", "MySQL username")
	rootCmd.Flags().StringVar(&password, "password", "", "MySQL password")
	rootCmd.Flags().StringVarP(&databaseName, "database", "d", "", "MySQL database name (optional, will create if not exists)")
	rootCmd.Flags().StringVar(&charset, "charset", "utf8mb4", "Character set")
	rootCmd.Flags().IntVar(&timeout, "timeout", 300, "Execution timeout in seconds")
	rootCmd.Flags().BoolVar(&noTransaction, "no-transaction", false, "Disable transaction mode (auto-commit)")
}
