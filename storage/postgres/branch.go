package postgres

import (
	"database/sql"
	"fmt"
	"init/models"
	
	"strconv"
	"time"

	"github.com/google/uuid"
)


type branchRepo struct {
	db *sql.DB
}


func NewBranchRepo(db *sql.DB) *branchRepo{
	return &branchRepo{
		db:db,
	}
}






func (r *branchRepo) Create(req *models.CreateBranch) (*models.Branch, error) {

	var (
		branchId = uuid.New().String()
		query      = `
			INSERT INTO "branches"(
				"id",
    			"name",
    			"phone",
    			"photo",
    			"work_start_hour",
    			"work_end_hour",
    			"address"
			) VALUES ($1, $2, $3,$4,$5,$6,$7)`
	)

	_, err := r.db.Exec(
		query,
		branchId,
		req.Name,
		req.Phone,
		req.Photo,
		req.WorkStart,
		req.WorkEnd,
		req.Address,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.BranchPrimaryKey{Id: branchId})
}

func (r *branchRepo) GetByID(req *models.BranchPrimaryKey) (*models.Branch, error) {

	var (
		branch models.Branch
		query    = `
			SELECT
				"id",
				"name",
				"phone",
				"photo",
				"work_start_hour",
				"work_end_hour",
				"address"
			FROM "branches"
			WHERE "id" = $1
		`
	)

	err := r.db.QueryRow(query, req.Id).Scan(
		&branch.Id,
		&branch.Name,
		&branch.Phone,
		&branch.Photo,
		&branch.WorkStart,
		&branch.WorkEnd,
		&branch.Address,
	)

	if err != nil {
		return nil, err
	}

	return &branch, nil
}

func (r *branchRepo) GetList(req *models.GetListBranchRequest) (*models.GetListBranchResponse, error) {
	var (
		resp   models.GetListBranchResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
		sort   = " ORDER BY created_at DESC"
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	hour:= strconv.Itoa(time.Now().Hour())
	minute:=strconv.Itoa(time.Now().Minute())
	currentTime:= hour +":" + minute
	if len(req.Search) > 0 {
		where += " AND name ILIKE" + " '%" + req.Search + "%'" + " AND " +"'" + currentTime + "'" +" BETWEEN work_start_hour AND work_end_hour"
 	}

	if len(req.Query) > 0 {
		where += req.Query
	}
	
	
	
	
	

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
    		"name",
    		"phone",
    		"photo",
    		"work_start_hour",
    		"work_end_hour",
    		"address",
			"delivery_price",
			"status"
		FROM "branches"
	`

	query += where + sort + offset + limit
	fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			branch models.Branch
			
		)

		err = rows.Scan(
			&resp.Count,
			&branch.Id,
			&branch.Name,
			&branch.Phone,
			&branch.Photo,
			&branch.WorkStart,
			&branch.WorkEnd,
			&branch.Address,
			&branch.DeliveryPrice,
			&branch.Status,
		)
		if err != nil {
			return nil, err
		}

		resp.Branches = append(resp.Branches, &branch)
	}
	fmt.Println(&resp)
	return &resp, nil
}

func (r *branchRepo) Update(req *models.UpdateBranch) (int64, error) {

	query := `
		UPDATE branches
			SET
			"work_start_hour" = $2,
			"work_end_hour" =$3,
			"delivery_price"=$4,
			"phone" =$5,
			"photo" =$6
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.Id,
		req.WorkStart,
		req.WorkEnd,
		req.DeliveryPrice,
		req.Phone,
		req.Photo,
	)
	fmt.Println(query)
	if err != nil {
		return 0, err
	}

	fmt.Println("ok")
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *branchRepo) Delete(req *models.BranchPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM branches WHERE id = $1", req.Id)
	return err
}
