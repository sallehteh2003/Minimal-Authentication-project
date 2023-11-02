package Database

import "golang.org/x/crypto/bcrypt"

func (gdb *DB) GetUserByEmail(email string) (*User, error) {
	u := User{}
	err := gdb.sql.Where(User{Email: email}).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (gdb *DB) CheckUserDuplicateByEmail(email string) (bool, error) {
	u := User{}
	err := gdb.sql.Where(User{Email: email}).First(&u).Error
	if err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (gdb *DB) CreateNewUser(u *User) error {
	if pw, err := bcrypt.GenerateFromPassword([]byte(u.Password), 0); err != nil {
		return err
	} else {
		u.Password = string(pw)
	}

	err := gdb.sql.Create(u).Error
	return err
}
