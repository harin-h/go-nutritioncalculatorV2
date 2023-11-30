package repository

import "github.com/jmoiron/sqlx"

type menuRepositoryDB struct {
	db *sqlx.DB
}

func NewMenuRositoryDB(db *sqlx.DB) menuRepositoryDB {
	return menuRepositoryDB{db: db}
}

func (r menuRepositoryDB) CreateMenu(menu Menu) (*Menu, error) {
	var menuId int
	err := r.db.QueryRow("INSERT INTO nutritioncalculator_menu (name,protein,fat,carb,creator_id,status,created_timestamp) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		menu.Name,
		menu.Protein,
		menu.Fat,
		menu.Carb,
		menu.CreatorId,
		menu.Status,
		menu.CreatedTimestamp).Scan(&menuId)
	if err != nil {
		return nil, err
	}
	menu.Id = menuId
	return &menu, nil
}

func (r menuRepositoryDB) GetAllMenues() ([]Menu, error) {
	var menues []Menu
	err := r.db.Select(&menues,
		`SELECT menu.id, menu.name, menu.protein , menu.fat, menu.carb , menu.creator_id , u1.username AS creator_name, menu.status, menu.created_timestamp, menu.status, COUNT(u2.user_id) AS count_like
		FROM (nutritioncalculator_menu AS menu INNER JOIN nutritioncalculator_user AS u1 ON menu.creator_id = u1.user_id ) 
		LEFT JOIN nutritioncalculator_user AS u2 ON CAST(menu.id AS text) = ANY(regexp_split_to_array(u2.favorite_menues,','))
		GROUP BY 1, 2,3,4,5,6,7,8,9, 10`)
	if err != nil {
		return nil, err
	}
	return menues, nil
}

func (r menuRepositoryDB) GetMenuById(id int) (*Menu, error) {
	var menu Menu
	err := r.db.Get(&menu,
		`SELECT menu.id, menu.name, menu.protein , menu.fat, menu.carb , menu.creator_id , u1.username AS creator_name, menu.status, menu.created_timestamp, menu.status, COUNT(u2.user_id) AS count_like
		FROM (nutritioncalculator_menu AS menu INNER JOIN nutritioncalculator_user AS u1 ON menu.creator_id = u1.user_id ) 
		LEFT JOIN nutritioncalculator_user AS u2 ON CAST(menu.id AS text) = ANY(regexp_split_to_array(u2.favorite_menues,','))
		WHERE menu.id = $1
		GROUP BY 1, 2,3,4,5,6,7,8,9, 10`,
		id)
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r menuRepositoryDB) UpdateMenu(menu Menu) error {
	tx := r.db.MustBegin()
	tx.MustExec("UPDATE nutritioncalculator_menu SET status=0 WHERE id=$1",
		menu.Id)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
