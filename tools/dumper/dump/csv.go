package dump

import (
	"DarkFlameMaster/config"
	"DarkFlameMaster/ticket"
	"encoding/csv"
	"fmt"
	"io"
)

func ToCSV(w io.Writer, tk []*ticket.Ticket) error {
	cfg := config.GetGlobalConfig()
	proofName := cfg.ProofName
	additionalName := cfg.AdditionalName
	if proofName == "" {
		proofName = "CustomerProof"
	}
	if additionalName == "" {
		additionalName = "Additional"
	}

	writer := csv.NewWriter(w)
	err := writer.Write([]string{"ID", proofName, "Row", "Column", "CreateTime", additionalName})
	if err != nil {
		return err
	}
	for _, t := range tk {
		err = writer.Write([]string{
			t.ID,
			t.CustomerProof,
			fmt.Sprintf("%d", t.Row),
			fmt.Sprintf("%d", t.Column),
			t.CreateTime.Format("2006-01-02 15:04:05"),
			t.Additional})
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
