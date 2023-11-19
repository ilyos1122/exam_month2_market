package postgres

import (
	"database/sql"
	
	"init/models"

	"github.com/google/uuid"
)

type orderProductRepo struct {
	db *sql.DB
}



func NewOrderProductRepo(db *sql.DB) *orderProductRepo {
	return &orderProductRepo{
		db: db,
	}
}

func (r *orderProductRepo) Create(req *models.CreateOrderProduct) error {
	orderProductId := uuid.New().String()

	query := `INSERT INTO order_products(
					"order_id",
					"order_product_id",
					"product_id",
					"discount_type",
					"discount_amount",
					"quantity"
					
				)VALUES ($1,$2,$3,$4,$5,$6)
				
	`
	_, err := r.db.Exec(
		query,
		req.OrderID,
		orderProductId,
		req.ProductID,
		req.DiscountType,
		req.DiscountAmount,
		req.Quantity,
		
	)
	if err != nil {
		return err
	}

	return nil

}

func (r *orderProductRepo) Delete(req *models.OrderProductPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM order_products WHERE order_product_id = $1", req.OrderProductId)

	return err
}
