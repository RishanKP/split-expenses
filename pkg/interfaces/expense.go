package interfaces

import (
	"errors"
	"split-expenses/library/utils"
	"split-expenses/pkg/models"
)

type CreateExpenseInput struct {
	Amount       float64              `json:"amount"`
	SplitType    string               `json:"splitType"`
	Description  string               `json:"description"`
	Participants []models.Participant `json:"participants"`
}

func (c CreateExpenseInput) AsExpense() (models.Expense, error) {
	if !utils.Contains(c.SplitType, []string{models.SPLIT_TYPE_EQUAL, models.SPLIT_TYPE_EXACT, models.SPLIT_TYPE_PERCENTAGE}) {
		return models.Expense{}, errors.New("invalid split type")
	}

	if c.SplitType == models.SPLIT_TYPE_PERCENTAGE {
		totalPercentage := 0
		for i := range c.Participants {
			c.Participants[i].Amount = utils.GetPercentageAmount(c.Amount, c.Participants[i].Percentage)
			totalPercentage += int(c.Participants[i].Percentage)
		}

		if totalPercentage != 100 {
			return models.Expense{}, errors.New("Percentage should add upto 100")
		}
	}

	if c.SplitType == models.SPLIT_TYPE_EQUAL {
		for i := range c.Participants {
			c.Participants[i].Amount = c.Amount / float64(len(c.Participants))
		}
	}

	return models.Expense{
		Amount:       c.Amount,
		Description:  c.Description,
		Participants: c.Participants,
	}, nil
}
