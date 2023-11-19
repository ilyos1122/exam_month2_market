package postgres

import (
	"database/sql"
	"fmt"
	"init/models"
	"init/pkg/helpers"

	"github.com/google/uuid"
)

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *categoryRepo {
	return &categoryRepo{
		db: db,
	}
}

func (r *categoryRepo) Create(req *models.CreateCategory) (*models.Category, error) {

	var (
		categoryId = uuid.New().String()
		query      = `
			INSERT INTO "category"(
				"id",
				"title",
				"image",
				"parent_id", 
				"updated_at"
			) VALUES ($1, $2, $3, $4 , NOW())`
	)

	_, err := r.db.Exec(
		query,
		categoryId,
		req.CategoryTitle,
		req.CategoryImage,
		helpers.NewNullString(req.ParentID),
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(&models.CategoryPrimaryKey{Id: categoryId})
}

func (r *categoryRepo) GetByID(req *models.CategoryPrimaryKey) (*models.Category, error) {

	var (
		category models.Category
		query    = `
			SELECT
				"id",
				"title",
				"image",
				COALESCE(CAST("parent_id" AS VARCHAR), ''),
				"created_at",
				"updated_at"	
			FROM "category"
			WHERE "id" = $1
		`
	)

	err := r.db.QueryRow(query, req.Id).Scan(
		&category.Id,
		&category.CategoryTitle,
		&category.CategoryImage,
		&category.ParentID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepo) GetList(req *models.GetListCategoryRequest) (*models.GetListCategoryResponse, error) {
	var (
		resp   models.GetListCategoryResponse
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

	if len(req.Search) > 0 {
		where += " AND title ILIKE" + " '%" + req.Search + "%'"
	}

	if len(req.Query) > 0 {
		where += req.Query
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"id",
			"category_title",
			"image",
			"parent_id",
			"created_at",
			"updated_at"
		FROM "category"
	`

	query += where + sort + offset + limit
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			category models.Category
			parentID sql.NullString
			categoryImage sql.NullString
		)

		err = rows.Scan(
			&resp.Count,
			&category.Id,
			&category.CategoryTitle,
			&categoryImage,
			&parentID,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		category.CategoryImage = categoryImage.String
		category.ParentID = parentID.String
		resp.Categories = append(resp.Categories, &category)
	}

	return &resp, nil
}

func (r *categoryRepo) Update(req *models.UpdateCategory) (int64, error) {

	query := `
		UPDATE category
			SET
				category_title = $2,
				parent_id = $3
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.Id,
		req.CategoryTitle,
		helpers.NewNullString(req.ParentID),
	)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (r *categoryRepo) Delete(req *models.CategoryPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM category WHERE id = $1", req.Id)
	return err
}
