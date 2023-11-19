package postgres

import (
	"database/sql"
	"fmt"
	"init/models"
	"init/pkg/helpers"

	"github.com/google/uuid"
)

type orderRepo struct {
	db *sql.DB
}

func NewOrderRepo(db *sql.DB) *orderRepo {
	return &orderRepo{
		db: db,
	}
}

var OrderCount = 1

func (r *orderRepo) Create(req *models.CreateOrder) (*models.Order, error) {
	var (
		id      = uuid.New().String()
		orderID = helpers.GenerateProductOrOrderID("O", OrderCount)
		query   = `INSERT INTO orders(
			"id",
			"order_id",
			"client_id",
			"branch_id",
			"delivery_address"
			
		) VALUES ($1,$2,$3,$4,$5)
		`
	)

	_, err := r.db.Exec(
		query,
		id,
		orderID,
		req.ClientID,
		req.BranchID,
		req.DeliveryAddress,
	)
	OrderCount++
	if err != nil {
		return nil, err
	}
	fmt.Println("okok")
	return r.GetByID(&models.OrderPrimaryKey{ID: id})
}

func (r *orderRepo) GetByID(req *models.OrderPrimaryKey) (*models.Order, error) {
	var (
		query2 = `SELECT 
					"order_product_id",
					"order_id",
					"product_id",
					"discount_type",
					"discount_amount",
					"quantity",
					"price",
					"sum",
					"created_at"
				FROM order_products WHERE order_id = $1
		`
		query = `SELECT 
			orders."id",
			orders."client_id",
			orders."branch_id",
			orders."delivery_address",
			orders."delivery_price",
			orders."total_count",
			orders."total_price",
			orders."status",
			orders."created_at",
			orders."updated_at"
		FROM orders 
		WHERE orders.id=$1
		`
	)
	var order models.Order;

	var (
		id               sql.NullString
		client_id        sql.NullString
		branch_id        sql.NullString
		delivery_address sql.NullString
		delivery_price   sql.NullInt64
		total_count      sql.NullInt64
		total_price      sql.NullInt64
		order_status     sql.NullString
		created_at       sql.NullString
		updated_at       sql.NullString
		
	)
	fmt.Println("OK2")
	err := r.db.QueryRow(query,req.ID).Scan(
		&id,
		&client_id,
		&branch_id,
		&delivery_address,
		&delivery_price,
		&total_count,
		&total_price,
		&order_status,
		&created_at,
		&updated_at,
		
	)
	if err != nil {

		return nil, err

	}
	rows,_ := r.db.Query(query2,req.ID)
	for rows.Next(){
		var (
		order_product_id sql.NullString
		order_id         sql.NullString
		product_id       sql.NullString
		discount_type    sql.NullString
		discount_amount  sql.NullInt64
		quantity         sql.NullInt64
		price            sql.NullInt64
		sum              sql.NullInt64
		opcreated_at     sql.NullString
		)

		err:= rows.Scan(
			&order_product_id,
			&order_id,
			&product_id,
			&discount_type,
			&discount_amount,
			&quantity,
			&price,
			&sum,
			&opcreated_at,
		)
		if err != nil {
			return nil, err
		}
		order.OrderProduct = append(order.OrderProduct, models.OrderProduct{
			OrderProductId: order_product_id.String,
			OrderID: order_id.String,
			ProductID: product_id.String,
			DiscountType: discount_type.String,
			DiscountAmount: int(discount_amount.Int64),
			Quantity: int(quantity.Int64),
			Price: int(price.Int64),
			Sum: int(sum.Int64),
			CreatedAt: opcreated_at.String,
		})
	}

	

	fmt.Println("OK1")

	return &models.Order{
		ID:              id.String,
		ClientID:        client_id.String,
		BranchID:        branch_id.String,
		DeliveryAddress: delivery_address.String,
		DeliveryPrice:   int(delivery_price.Int64),
		TotalCount:      int(total_count.Int64),
		TotalPrice:      int(total_price.Int64),
		OrderStatus:     order_status.String,
		CreatedAt:       created_at.String,
		UpdatedAt:       updated_at.String,
		OrderProduct: order.OrderProduct,
	}, nil

}

func (r *orderRepo) GetList(req *models.GetListOrderRequest) (*models.GetListOrderResponse, error) {
	var (
		resp   models.GetListOrderResponse
		order  models.Order
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
		where += " AND title ILIKE" + " '%" + req.Search + "%'"
	}

	var query = `
		SELECT
			COUNT(*) OVER(),
			"order_id",
			"client_id",
			"branch_id",
			"delivery_address",
			"deliver_price",
			"total_count",
			"total_price",
			"order_status",
			"created_at",
			"updated_at"
		FROM orders
	`

	query += where + offset + limit

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {

		err := rows.Scan(
			&resp.Count,
			&order.ID,
			&order.ClientID,
			&order.BranchID,
			&order.DeliveryAddress,
			&order.DeliveryPrice,
			&order.TotalCount,
			&order.TotalPrice,
			&order.OrderStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Orders = append(resp.Orders, &models.Order{
			ID:              order.ID,
			ClientID:        order.ClientID,
			BranchID:        order.BranchID,
			DeliveryAddress: order.DeliveryAddress,
			DeliveryPrice:   order.DeliveryPrice,
			TotalCount:      order.TotalCount,
			TotalPrice:      order.TotalPrice,
			OrderStatus:     order.OrderStatus,
			CreatedAt:       order.CreatedAt,
			UpdatedAt:       order.UpdatedAt,
		})
	}

	return &resp, nil

}

func (r *orderRepo) Update(req *models.UpdateOrder) (int64, error) {

	query := `
		UPDATE orders
			SET
			"delivery_address" = $2,
			"branch_id" = $3,
		WHERE id = $1
	`
	result, err := r.db.Exec(
		query,
		req.OrderID,
		req.DeliveryAddress,
		req.BranchID,
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

func (r *orderRepo) Delete(req *models.OrderPrimaryKey) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = $1", req.ID)
	return err
}

func (r *orderRepo) ChangeStatus(req *models.ChangeOrderStatus) (int64, error) {
	var query = `UPDATE orders SET order_status = $1 WHERE id = $2`
	var query1 = `Select "status" from orders WHERE id =$1`

	var order models.Order
	_ = r.db.QueryRow(query1, req.ID).Scan(
		&order.OrderStatus,
	)
	if req.Status == "canceled" || order.OrderStatus == "new" {
		result, err := r.db.Exec(query, req.Status, req.ID)
		if err != nil {
			return 0, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		return rowsAffected, nil
	}
	result, err := r.db.Exec(query, req.Status, req.ID)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
