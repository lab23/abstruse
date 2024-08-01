package auth

import (
	"fmt"

	"github.com/lab23/abstruse/pkg/fs"
)

func generateHtpasswdFile(filePath, user, password string) error {
	passwd, err := HashPassword(Password{Password: password, Cost: 1})
	if err != nil {
		return err
	}
	creds := fmt.Sprintf("%s:%s\n", user, passwd)
	return fs.WriteFile(filePath, creds)
}
