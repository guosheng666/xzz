package sqlfile

import (
	"fmt"
	"os"
	"strings"
)

func FileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("SQL file does not exist: %s", path)
	}
	return nil
}

func ParseSQLFile(path string) ([]string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var statements []string
	var currentStatement strings.Builder
	var delimiter = ";"

	lines := strings.Split(string(content), "\n")
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		if trimmedLine == "" {
			continue
		}
		
		if strings.HasPrefix(trimmedLine, "--") || strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		
		if strings.HasPrefix(strings.ToUpper(trimmedLine), "DELIMITER") {
			parts := strings.Fields(trimmedLine)
			if len(parts) >= 2 {
				delimiter = parts[1]
			}
			continue
		}
		
		currentStatement.WriteString(line)
		currentStatement.WriteString("\n")
		
		if strings.HasSuffix(trimmedLine, delimiter) {
			stmt := strings.TrimSpace(currentStatement.String())
			if stmt != "" {
				stmt = strings.TrimSuffix(stmt, delimiter)
				stmt = strings.TrimSpace(stmt)
				if stmt != "" {
					statements = append(statements, stmt)
				}
			}
			currentStatement.Reset()
		}
	}

	if currentStatement.Len() > 0 {
		stmt := strings.TrimSpace(currentStatement.String())
		stmt = strings.TrimSuffix(stmt, delimiter)
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements, nil
}
