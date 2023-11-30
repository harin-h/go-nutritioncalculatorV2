package repository

import "github.com/jmoiron/sqlx"

type userRepositoryDB struct {
	db *sqlx.DB
}

func NewUserRepositoryDB(db *sqlx.DB) userRepositoryDB {
	return userRepositoryDB{db: db}
}

func (r userRepositoryDB) GetUserById(userId string) (*User, error) {
	user := User{}
	err := r.db.Get(&user,
		`SELECT 
		*
	FROM nutritioncalculator_user
	WHERE user_id=$1`,
		userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepositoryDB) GetUserByUsername(username string) (*User, error) {
	user := User{}
	err := r.db.Get(&user,
		`SELECT 
		*
	FROM nutritioncalculator_user
	WHERE username=$1`,
		username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r userRepositoryDB) CreateUser(user User) error {
	tx := r.db.MustBegin()
	tx.MustExec("INSERT INTO nutritioncalculator_user (user_id,password,username,weight,protein,fat,carb,favorite_menues,created_timestamp) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)",
		user.UserId,
		user.Password,
		user.Username,
		user.Weight,
		user.Protein,
		user.Fat,
		user.Carb,
		user.FavoriteMenues,
		user.CreatedTimestamp)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (r userRepositoryDB) UpdateUser(user User) error {
	tx := r.db.MustBegin()
	tx.MustExec("UPDATE nutritioncalculator_user SET password=$1,username=$2,weight=$3,protein=$4,fat=$5,carb=$6,favorite_menues=$7 WHERE user_id=$8",
		user.Password,
		user.Username,
		user.Weight,
		user.Protein,
		user.Fat,
		user.Carb,
		user.FavoriteMenues,
		user.UserId)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
