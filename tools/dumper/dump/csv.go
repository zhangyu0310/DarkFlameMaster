package dump

import (
	"DarkFlameMaster/ticket"
	"encoding/csv"
	"fmt"
	"os"
)

func ToCSV(filePath string, tk []*ticket.Ticket) error {
	// create dump file and write data
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)
	err = writer.Write([]string{"ID", "CustomerProof", "Row", "Column", "CreateTime"})
	if err != nil {
		return err
	}
	for _, t := range tk {
		err = writer.Write([]string{
			t.ID,
			t.CustomerProof,
			fmt.Sprintf("%d", t.Row),
			fmt.Sprintf("%d", t.Column),
			t.CreateTime.Format("2006-01-02 15:04:05")})
		if err != nil {
			return err
		}
	}
	writer.Flush()
	_ = file.Close()
	return nil
}
