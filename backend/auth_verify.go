package backend

import "github.com/BlahajXD/backend/logic"

func (d *Dependency) AuthVerify(accessToken string) bool {
	_, _, err := logic.VerifyJWT(accessToken)
	if err != nil {
		return false
	}

	return true
}
