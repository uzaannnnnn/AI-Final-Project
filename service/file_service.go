package service

import (
    "encoding/csv"
    "errors"
    "strings"
	
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
)   

type FileService struct {
    Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
    reader := csv.NewReader(strings.NewReader(fileContent))
    rows, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    if len(rows) < 1 {
        return nil, errors.New("csv file is empty or missing header")
    }

    headers := rows[0]
    table := make(map[string][]string)
    for _, header := range headers {
        table[strings.TrimSpace(header)] = []string{}
    }

    for _, row := range rows[1:] {
        if len(row) != len(headers) {
            return nil, errors.New("csv row has incorrect number of fields")
        }
        for i, value := range row {
            table[strings.TrimSpace(headers[i])] = append(table[strings.TrimSpace(headers[i])], strings.TrimSpace(value))
        }
    }

    return table, nil
}
