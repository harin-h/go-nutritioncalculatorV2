package repository

import "github.com/jmoiron/sqlx"

type recordRepositoryDB struct {
	db *sqlx.DB
}

func NewRecordRepositoryDB(db *sqlx.DB) recordRepositoryDB {
	return recordRepositoryDB{db: db}
}

func (r recordRepositoryDB) GetRecordsByUserId(userId string) ([]Record, error) {
	records := []Record{}
	err := r.db.Select(&records,
		`SELECT id, user_id , list, note, weight, status, created_timestamp , event_timestamp, string_agg(concat(menu_name, '-',nums,' ') ,',') AS menues, SUM(nums * protein ) AS protein, SUM(nums * fat ) AS fat, SUM(nums * carb ) AS carb, MIN(menu_status) AS is_updated
		FROM 
		(
		SELECT r.id, r.user_id, r.list, r.note, r.weight, r.status, r.created_timestamp, r.event_timestamp,
		cardinality(regexp_split_to_array(r.list,',')) - cardinality(array_remove(regexp_split_to_array(r.list,','),CAST(m.id AS text))) AS nums
		, m."name" AS menu_name, m.protein , m.fat, m.carb , m.status AS menu_status
		FROM nutritioncalculator_record AS r LEFT JOIN nutritioncalculator_menu AS m  
		ON CAST(m.id AS text) = ANY(regexp_split_to_array(r.list,','))
		WHERE r.user_id = $1 AND r.status = 1
		) AS t
		GROUP BY 1, 2, 3, 4, 5, 6, 7, 8`,
		userId)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (r recordRepositoryDB) GetRecordById(recordId int) (*Record, error) {
	record := Record{}
	err := r.db.Get(&record,
		`SELECT id, user_id , list, note, weight, status, created_timestamp , event_timestamp, string_agg(concat(menu_name, '-',nums,' ') ,',') AS menues, SUM(nums * protein ) AS protein, SUM(nums * fat ) AS fat, SUM(nums * carb ) AS carb, MIN(menu_status) AS is_updated
		FROM 
		(
		SELECT r.id, r.user_id, r.list, r.note, r.weight, r.status, r.created_timestamp , r.event_timestamp,
		cardinality(regexp_split_to_array(r.list,',')) - cardinality(array_remove(regexp_split_to_array(r.list,','),CAST(m.id AS text))) AS nums
		, m."name" AS menu_name, m.protein , m.fat, m.carb , m.status AS menu_status
		FROM nutritioncalculator_record AS r LEFT JOIN nutritioncalculator_menu AS m  
		ON CAST(m.id AS text) = ANY(regexp_split_to_array(r.list,','))
		WHERE r.id = $1 AND r.status = 1
		) AS t
		GROUP BY 1, 2, 3, 4, 5, 6, 7, 8`,
		recordId)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r recordRepositoryDB) CreateRecord(record Record) (*Record, error) {
	var recordId int
	err := r.db.QueryRow("INSERT INTO nutritioncalculator_record (user_id,List,weight,note,event_timestamp,status,created_timestamp) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id",
		record.UserId,
		record.List,
		record.Weight,
		record.Note,
		record.EventTimestamp,
		record.Status,
		record.CreatedTimestamp).Scan(&recordId)
	if err != nil {
		return nil, err
	}
	record.Id = recordId
	return &record, nil
}

func (r recordRepositoryDB) UpdateRecord(record Record) error {
	tx := r.db.MustBegin()
	tx.MustExec("UPDATE nutritioncalculator_record SET list=$1,note=$2,weight=$3,event_timestamp=$4,status=$5 WHERE id=$6",
		record.List,
		record.Note,
		record.Weight,
		record.EventTimestamp,
		record.Status,
		record.Id)
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
