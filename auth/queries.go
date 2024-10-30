package main

var (
	AddUserQuery string = `
		INSERT INTO users (username, password_hash, email, activation_code) 
		VALUES ($1, crypt($2, gen_salt('bf')), $3, CAST(1000000000 + floor(random() * 9000000000) AS bigint))
		RETURNING activation_code;`

	ChangePasswordQuery string = `
		UPDATE users 
		SET password_hash = crypt($1, gen_salt('bf')) 
		WHERE id = $2;`

	AttemptLoginQuery string = `
		SELECT (password_hash = crypt($1, password_hash) AND activated)
    	AS password_match 
		FROM users 
		WHERE username = $2;`

	GetUserNameQuery string = `
		SELECT username
		FROM users
		WHERE id = $1;`

	SetActivatedQuery string = `
		UPDATE users
		set activated = TRUE
		WHERE username = $1;`

	CheckActivationCodeQuery string = `
		SELECT (activation_code = $1 AND activation_tries < 5)
		AS code_match
		FROM users
		WHERE username = $2;`

	IncreaseActivationTriesQuery string = `
		UPDATE users
		SET activation_tries = activation_tries + 1
		WHERE username = $1;`
)

func (as *authServer) registerUser(username, password, email string) (string, error) {
	var activationCode string
	err := as.db.Database.QueryRow(AddUserQuery, username, password, email).Scan(&activationCode)

	if err != nil {
		return "", err
	}

	return activationCode, nil
}

func (as *authServer) changePassword(id, newPassword string) error {
	_, err := as.db.Database.Exec(ChangePasswordQuery, newPassword, id)

	return err
}

func (as *authServer) attemptLogin(username, password string) (bool, error) {
	var result string
	err := as.db.Database.QueryRow(AttemptLoginQuery, password, username).Scan(&result)

	if err != nil {
		return false, err
	}

	return result == "t" || result == "true", nil
}

func (as *authServer) getUserName(id string) (string, error) {
	var result string
	err := as.db.Database.QueryRow(GetUserNameQuery).Scan(&result)

	if err != nil {
		return "", err
	}

	return result, nil
}

func (as *authServer) activateUser(username string) error {
	_, err := as.db.Database.Exec(SetActivatedQuery, username)

	return err
}

func (as *authServer) checkActivationCode(activationCode string, username string) (bool, error) {
	var result string
	err := as.db.Database.QueryRow(CheckActivationCodeQuery, activationCode, username).Scan(&result)

	if err != nil {
		return false, err
	}

	return result == "t" || result == "true", nil
}

func (as *authServer) increaseActivationTries(username string) error {
	_, err := as.db.Database.Exec(IncreaseActivationTriesQuery, username)

	return err
}
