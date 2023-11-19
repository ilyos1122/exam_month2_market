package postgres

import (
	"database/sql"
	"fmt"
	"init/models"
	"init/pkg/helpers"

	"github.com/google/uuid"
)

var ProductCount = 1

type productRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *productRepo {
	return &productRepo{
		db: db,
	}
}

func (r *productRepo) Create(req *models.CreateProduct) (*models.Product, error) {

	var (
		productUUID = uuid.New().String()
		productId   = helpers.GenerateProductOrOrderID("P", ProductCount)
		query       = `
			INSERT INTO "product"(
				"id",
				"product_id",	
				"title",
				"description",
				"price",
				"product_image",
				"category_id"
			) VALUES ($1, $2, $3, $4, $5,$6,$7)
			`
	)

	_, err := r.db.Exec(
		query,
		productUUID,
		productId,
		req.Title,
		req.Description,
		req.Price,
		req.Image,
		helpers.NewNullString(req.CategoryId),
	) 
	if err != nil {
		return nil, err

	}
	ProductCount++
	fmt.Println(ProductCount)

	return r.GetByID(&models.ProductPrimaryKey{Id: productUUID})
}

func (r *productRepo) GetByID(req *models.ProductPrimaryKey) (*models.Product, error) {

	var (
		query = `
			SELECT
				p."id",
				p."product_id",
				p."title",
				p."description",
				p."price",
				p."product_image",
				p."category_id",
				p."created_at",
				p."updated_at",
				c."id",
				c."category_title",
				c."image",
				c."parent_id",
				c."created_at",
				c."updated_at"

			FROM "product" AS p
			JOIN "category" AS c ON c.id = p.category_id
			WHERE p."id" = $1
		`
	)

	var (
		id                  sql.NullString
		produt_id           sql.NullString
		title               sql.NullString
		description         sql.NullString
		price               sql.NullFloat64
		image               sql.NullString
		category_id         sql.NullString
		updated_at          sql.NullString
		created_at          sql.NullString
		category_uuid_id    sql.NullString
		category_title      sql.NullString
		category_image      sql.NullString
		category_parent_id  sql.NullString
		category_created_at sql.NullString
		category_updated_at sql.NullString
	)

	err := r.db.QueryRow(query, req.Id).Scan(
		&id,
		&produt_id,
		&title,
		&description,
		&price,
		&image,
		&category_id,
		&updated_at,
		&created_at,
		&category_uuid_id,
		&category_title,
		&category_image,
		&category_parent_id,
		&category_created_at,
		&category_updated_at,
	)
	if err != nil {
		return nil, err
	}
	

	return &models.Product{
		Id:          id.String,
		ProductID:   produt_id.String,
		Title:       title.String,
		Description: description.String,
		Price:       price.Float64,
		Image:       image.String,
		CategoryId:  category_id.String,
		UpdatedAt:   updated_at.String,
		CreatedAt:   created_at.String,
		Category:   models.Category{
			Id: category_uuid_id.String,
			CategoryTitle: category_title.String,
			CategoryImage: category_image.String,
			ParentID: category_parent_id.String,
			CreatedAt: category_created_at.String,
			UpdatedAt: category_updated_at.String,
			},
	}, nil
}

func (r *productRepo) GetList(req *models.GetListProductRequest) (*models.GetListProductResponse, error) {
	var (
		resp   models.GetListProductResponse
		where  = " WHERE TRUE"
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d;", req.Limit)
	}

	if len(req.Search) > 0 {
		where += " AND title ILIKE" + " '%" + req.Search + "%'" + " OR category.category_title ILIKE" + " '%" + req.Search + "%'"
	}

	var query = `
	SELECT
		COUNT(*) OVER(),
		product."id",
		product."product_id",
		product."title",
		product."description",
		product."price",
		product."product_image",
		product."category_id",
		product."created_at",
		product."updated_at",
		category."id",
		category."category_title",
		category."image",
		category."parent_id",
		category."created_at",
		category."updated_at"
	FROM "product"
	JOIN "category" ON category.id = product.category_id
	`

	query += where + offset + limit

	// fmt.Println(query)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			id                  sql.NullString
			produt_id           sql.NullString
			title               sql.NullString
			description         sql.NullString
			price               sql.NullFloat64
			image               sql.NullString
			category_id         sql.NullString
			updated_at          sql.NullString
			created_at          sql.NullString
			category_uuid_id    sql.NullString
			category_title      sql.NullString
			category_image      sql.NullString
			category_parent_id  sql.NullString
			category_created_at sql.NullString
			category_updated_at sql.NullString
		)

		err := rows.Scan(
			&resp.Count,
			&id,
			&produt_id,
			&title,
			&description,
			&price,
			&image,
			&category_id,
			&updated_at,
			&created_at,
			&category_uuid_id,
			&category_title,
			&category_image,
			&category_parent_id,
			&category_created_at,
			&category_updated_at,
		)
		if err != nil {
			return nil, err
		}
		resp.Products = append(resp.Products, &models.Product{
			Id:          id.String,
			ProductID:   produt_id.String,
			Title:       title.String,
			Description: description.String,
			Price:       price.Float64,
			Image:       image.String,
			CategoryId:  category_id.String,
			UpdatedAt:   updated_at.String,
			CreatedAt:   created_at.String,
			Category:  models.Category{
				Id: category_uuid_id.String,
				CategoryTitle: category_title.String,
				CategoryImage: category_image.String,
				ParentID: category_parent_id.String,
				CreatedAt: category_created_at.String,
				UpdatedAt: category_updated_at.String,
				}  ,
		})
	}
	// fmt.Println(resp)
	return &resp, nil

}

func (r *productRepo) Update(req *models.UpdateProduct) (int64, error) {

	query := `
		UPDATE product
			SET
				title = $2,
				description = $3,
				price = $4,
				product_image = $5,
				category_id = $6
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.Id,
		req.Title,
		req.Description,
		req.Price,
		req.Image,
		helpers.NewNullString(req.CategoryId),
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

func (r *productRepo) Delete(req *models.ProductPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM product WHERE id = $1", req.Id)
	return err
}
