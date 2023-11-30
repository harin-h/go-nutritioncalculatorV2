package repository

import "github.com/jmoiron/sqlx"

type favListRepositoryDB struct {
	db *sqlx.DB
}

func NewFavListRepositoryDB(db *sqlx.DB) favListRepositoryDB {
	return favListRepositoryDB{db: db}
}

func (r favListRepositoryDB) GetFavListsByUserId(userId string) ([]FavList, error) {
	favLists := []FavList{}
	err := r.db.Select(&favLists,
		`SELECT id, user_id , name, list, status, created_timestamp , string_agg(concat(menu_name, '-',nums,' ') ,',') AS menues, SUM(nums * protein ) AS protein, SUM(nums * fat ) AS fat, SUM(nums * carb ) AS carb, MIN(menu_status) AS is_updated
		FROM 
		(
		SELECT fl.id, fl.user_id, fl.name, fl.list, fl.status, fl.created_timestamp,
		cardinality(regexp_split_to_array(fl.list,',')) - cardinality(array_remove(regexp_split_to_array(fl.list,','),CAST(m.id AS text))) AS nums
		, m."name" AS menu_name, m.protein , m.fat, m.carb , m.status AS menu_status
		FROM nutritioncalculator_favorite_list AS fl LEFT JOIN nutritioncalculator_menu AS m  
		ON CAST(m.id AS text) = ANY(regexp_split_to_array(fl.list,','))
		WHERE fl.user_id = $1 AND fl.status = 1
		) AS t
		GROUP BY 1, 2, 3, 4, 5, 6`,
		userId)
	if err != nil {
		return nil, err
	}
	return favLists, nil
}

func (r favListRepositoryDB) GetFavListById(favListId int) (*FavList, error) {
	favList := FavList{}
	err := r.db.Get(&favList,
		`SELECT id, user_id , name, list, status, created_timestamp , string_agg(concat(menu_name, '-',nums,' ') ,',') AS menues, SUM(nums * protein ) AS protein, SUM(nums * fat ) AS fat, SUM(nums * carb ) AS carb, MIN(menu_status) AS is_updated
		FROM 
		(
		SELECT fl.id, fl.user_id, fl.name, fl.list, fl.status, fl.created_timestamp,
		cardinality(regexp_split_to_array(fl.list,',')) - cardinality(array_remove(regexp_split_to_array(fl.list,','),CAST(m.id AS text))) AS nums
		, m."name" AS menu_name, m.protein , m.fat, m.carb , m.status AS menu_status
		FROM nutritioncalculator_favorite_list AS fl LEFT JOIN nutritioncalculator_menu AS m  
		ON CAST(m.id AS text) = ANY(regexp_split_to_array(fl.list,','))
		WHERE fl.id=$1 AND fl.status = 1
		) AS t
		GROUP BY 1, 2, 3, 4, 5, 6`,
		favListId)
	if err != nil {
		return nil, err
	}
	return &favList, nil
}

func (r favListRepositoryDB) CreateFavList(favList FavList) (*FavList, error) {
	var favListId int
	err := r.db.QueryRow("INSERT INTO nutritioncalculator_favorite_list (user_id,name,list,status,created_timestamp) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		favList.UserId,
		favList.Name,
		favList.List,
		favList.Status,
		favList.CreatedTimestamp).Scan(&favListId)
	if err != nil {
		return nil, err
	}
	favList.Id = favListId
	return &favList, nil
}

func (r favListRepositoryDB) UpdateFavList(favList FavList) error {
	tx := r.db.MustBegin()
	tx.MustExec("UPDATE nutritioncalculator_favorite_list SET name=$1,list=$2,status=$3 WHERE id=$4",
		favList.Name,
		favList.List,
		favList.Status,
		favList.Id)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
