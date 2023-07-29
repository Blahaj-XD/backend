package backend

import "context"

type BankAccountInfoOutput struct {
	//
}

func (d *Dependency) BankAccountInfo(ctx context.Context, parentID int) (BankAccountInfoOutput, error) {
	panic("not implemented")
}
