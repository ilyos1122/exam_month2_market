package postgres

import (
	"database/sql"
	"fmt"
	"init/models"


	"github.com/google/uuid"
)

type clientRepo struct {
	db *sql.DB
}



func NewClientRepo(db *sql.DB) *clientRepo {
	return &clientRepo{
		db: db,
	}
}

func (r *clientRepo) Create(req *models.CreateClient) (*models.Client, error) {
	var (
		clientId = uuid.New().String()
		query    = `
			INSERT INTO "clients"(
				"id",
				"first_name",
				"last_name",
				"phone",
				"photo",
				"date_of_birth"
			)VALUES ($1,$2,$3,$4,$5,$6)
		
		`
	)

	_, err := r.db.Exec(
		query,
		clientId,
		req.FirstName,
		req.LastName,
		req.Phone,
		req.Photo,
		req.DateOfBirth,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.ClientPrimaryKey{Id: clientId})
}


func (r *clientRepo) GetByID(req *models.ClientPrimaryKey)(*models.Client,error){
	var (
		client models.Client
		query = `SELECT 
		"id",
		"first_name",
		"last_name",
		"phone",
		"photo",
		"date_of_birth",
		"created_at"
		FROM "clients"
		WHERE "id" = $1
		
		`
	)
	fmt.Println("ok")
	err := r.db.QueryRow(query,req.Id).Scan(
		&client.Id,
		&client.FirstName,
		&client.LastName,
		&client.Phone,
		&client.Photo,
		&client.DateOfBirth,
		&client.CreatedAT,
	)
	if err != nil {
		return nil, err
	}
	fmt.Println("ok")
	return &client,nil
}



func (r *clientRepo) GetList(req *models.GetListClientRequest) (*models.GetListClientResponse,error){
	var (
		resp models.GetListClientResponse
		where = " WHERE TRUE "
		offset = "OFFSET 0"
		limit  = " LIMIT 10"
	)
	

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
	fmt.Println("Search: ",req.Search)
	if len(req.Search) > 0 {
		where += " AND first_name ILIKE " + "'%" + req.Search + "%'" + " OR last_name ILIKE" + "'%" + req.Search + "%'" + " OR phone ILIKE" + "'%" + req.Search + "%'"
	}
	var query = `SELECT 
		COUNT(*) OVER(),
		"id",
		"first_name",
		"last_name",
		"phone",
		"photo",
		"date_of_birth",
		"created_at"
		FROM clients
		`

	query += where + offset + limit 
	// fmt.Println(where)
	 fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next(){
		var (
			id					sql.NullString
			first_name			sql.NullString	
			last_name			sql.NullString	
			phone				sql.NullString
			photo				sql.NullString
			date_of_birth		sql.NullString	
			created_at			sql.NullString
		)
		err := rows.Scan(
			&resp.Count,
			&id,
			&first_name,
			&last_name,
			&phone,
			&photo,
			&date_of_birth,
			&created_at,
		)
		
		if err != nil {
			return nil,err
		}
		
		resp.Clients = append(resp.Clients,&models.Client{
			Id:          id.String,
			FirstName:   first_name.String,
			LastName:    last_name.String,
			Phone:       phone.String,
			Photo:       photo.String,
			DateOfBirth: date_of_birth.String,
			CreatedAT:   created_at.String,
		})
	}

	return &resp,nil


}


func (r *clientRepo) Update(req *models.UpdateClient) (int64,error){
	query := `
		UPDATE clients
		SET 
			phone =$1,
			photo = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(
		query,
		req.Phone,
		req.Photo,
		req.Id,
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected,nil
}

func (r *clientRepo) Delete(req *models.ClientPrimaryKey) error {
	_,err := r.db.Exec("DELETE FROM clients WHERE id = $1",req.Id)
	return err
}