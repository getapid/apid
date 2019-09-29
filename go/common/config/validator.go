package config

import "errors"

func Validate(c Config) error {
	if c == nil {
		return errors.New("The config yaml must not be empty")
	}
	if c.Version != "1" {
		return errors.New("The config yaml must contain property \"version\" with one of the following values: \"1\"")
	}
	if c.Transactions == nil || len(c.Transactions) == 0 {
		return errors.New("The config yaml must contain at least 1 transation")
	}
	for tran := range c.Transactions {
		if tran.ID == "" {
			// TODO: extract the list of IDs and ensure unique-ness, same for steps
			return errors.New("Each transaction in the config yaml must contain an \"id\" property")
		}
		if tran.Steps == nil || len(tran.Steps) == 0 {
			return errors.New("Each transaction must contain at least 1 step")
		}
	}
	return nil
}
