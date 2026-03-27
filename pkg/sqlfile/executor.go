package sqlfile

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

func ExecuteStatements(db *sql.DB, statements []string, useTransaction bool) error {
	if useTransaction {
		return executeWithTransaction(db, statements)
	}
	return executeWithoutTransactionParallel(db, statements)
}

func executeWithTransaction(db *sql.DB, statements []string) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	successCount := 0
	batchSize := 500

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute statement %d: %w", i+1, err)
		}
		successCount++
		if successCount%batchSize == 0 {
			fmt.Printf("Executed %d statements...\n", successCount)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	fmt.Printf("Successfully executed %d statements\n", successCount)
	return nil
}

func executeWithoutTransaction(db *sql.DB, statements []string) error {
	successCount := 0
	batchSize := 500

	for i, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		_, err := db.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute statement %d: %w", i+1, err)
		}
		successCount++
		if successCount%batchSize == 0 {
			fmt.Printf("Executed %d statements...\n", successCount)
		}
	}
	fmt.Printf("Successfully executed %d statements\n", successCount)
	return nil
}

func executeWithoutTransactionParallel(db *sql.DB, statements []string) error {
	const numWorkers = 10
	const batchSize = 1000
	
	var ddlStatements []string
	var dmlStatements []string
	
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		
		upperStmt := strings.ToUpper(stmt)
		if strings.HasPrefix(upperStmt, "CREATE") || 
		   strings.HasPrefix(upperStmt, "ALTER") ||
		   strings.HasPrefix(upperStmt, "DROP") ||
		   strings.HasPrefix(upperStmt, "TRUNCATE") {
			ddlStatements = append(ddlStatements, stmt)
		} else {
			dmlStatements = append(dmlStatements, stmt)
		}
	}
	
	if len(ddlStatements) > 0 {
		fmt.Printf("Executing %d DDL statements sequentially...\n", len(ddlStatements))
		for i, stmt := range ddlStatements {
			_, err := db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("failed to execute DDL statement %d: %w", i+1, err)
			}
			if (i+1)%100 == 0 {
				fmt.Printf("DDL: %d/%d completed\n", i+1, len(ddlStatements))
			}
		}
		fmt.Printf("DDL statements completed\n")
	}
	
	if len(dmlStatements) == 0 {
		return nil
	}
	
	fmt.Printf("Executing %d DML statements with %d workers...\n", len(dmlStatements), numWorkers)
	
	var wg sync.WaitGroup
	var mu sync.Mutex
	successCount := 0
	errorChan := make(chan error, numWorkers)
	
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			
			start := workerID * len(dmlStatements) / numWorkers
			end := (workerID + 1) * len(dmlStatements) / numWorkers
			if workerID == numWorkers-1 {
				end = len(dmlStatements)
			}
			
			for j := start; j < end; j++ {
				stmt := strings.TrimSpace(dmlStatements[j])
				if stmt == "" {
					continue
				}
				
				_, err := db.Exec(stmt)
				if err != nil {
					errorChan <- fmt.Errorf("worker %d failed at statement %d: %w", workerID, j+1, err)
					return
				}
				
				mu.Lock()
				successCount++
				if successCount%batchSize == 0 {
					fmt.Printf("Executed %d DML statements...\n", successCount)
				}
				mu.Unlock()
			}
		}(i)
	}
	
	go func() {
		wg.Wait()
		close(errorChan)
	}()
	
	if err := <-errorChan; err != nil {
		return err
	}
	
	fmt.Printf("Successfully executed %d statements\n", len(statements))
	return nil
}

func OptimizeStatements(statements []string) []string {
	var optimized []string
	var insertBuffer strings.Builder
	var currentTable string
	var currentColumns string
	var valueCount int
	const maxValuesPerBatch = 1000

	insertRegex := regexp.MustCompile(`(?i)^INSERT\s+INTO\s+` + "`" + `?(\w+)` + "`" + `?\s*\(([^)]+)\)\s*VALUES\s*`)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		matches := insertRegex.FindStringSubmatch(stmt)
		if matches != nil {
			table := matches[1]
			columns := strings.TrimSpace(matches[2])

			if table == currentTable && columns == currentColumns && valueCount < maxValuesPerBatch {
				valuesPart := strings.TrimSpace(strings.TrimPrefix(stmt, matches[0]))
				if strings.HasPrefix(valuesPart, "(") {
					if insertBuffer.Len() > 0 {
						insertBuffer.WriteString(",")
					}
					insertBuffer.WriteString(valuesPart)
					valueCount++
					continue
				}
			}

			if insertBuffer.Len() > 0 {
				optimized = append(optimized, fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;", currentTable, currentColumns, insertBuffer.String()))
				insertBuffer.Reset()
				valueCount = 0
			}

			currentTable = table
			currentColumns = columns
			valuesPart := strings.TrimSpace(strings.TrimPrefix(stmt, matches[0]))
			if strings.HasPrefix(valuesPart, "(") {
				insertBuffer.WriteString(valuesPart)
				valueCount = 1
				continue
			}
		}

		if insertBuffer.Len() > 0 {
			optimized = append(optimized, fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;", currentTable, currentColumns, insertBuffer.String()))
			insertBuffer.Reset()
			valueCount = 0
			currentTable = ""
			currentColumns = ""
		}

		optimized = append(optimized, stmt)
	}

	if insertBuffer.Len() > 0 {
		optimized = append(optimized, fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s;", currentTable, currentColumns, insertBuffer.String()))
	}

	return optimized
}
