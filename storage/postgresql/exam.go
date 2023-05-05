package postgresql

import (
	"app/api/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type ExamRepo struct {
	db *pgxpool.Pool
}

func NewExamRepo(db *pgxpool.Pool) *ExamRepo {
	return &ExamRepo{db: db}
}

//2-exam
func (r ExamRepo) EachStaff(ctx context.Context, req *models.Date) (res []models.StaffDate, err error) {
	query := `SELECT
    staffs.first_name || ' ' || staffs.last_name AS "employee",  categories.category_name AS "category",
       products.product_name AS "product",   order_items.quantity AS "quantity",   order_items.list_price * order_items.quantity AS "total"
FROM orders
         JOIN order_items ON orders.order_id = order_items.order_id
         JOIN products ON order_items.product_id = products.product_id
         JOIN categories ON products.category_id = categories.category_id
         JOIN staffs ON orders.staff_id = staffs.staff_id
WHERE orders.order_date = $1`

	var year string

	if req.Day == "" {
		dt := time.Now()
		year = dt.Format("2006-01-02")
	} else {
		year = req.Day
	}

	date, error := time.Parse("2006-01-02", year)
	if error != nil {
		fmt.Println(error)
		return
	}

	rows, err := r.db.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var s models.StaffDate
		err = rows.Scan(
			&s.StaffName,
			&s.Category,
			&s.Product,
			&s.Quantity,
			&s.Summ,
		)
		res = append(res, s)
		if err != nil {
			return res, err
		}
	}
	return res, nil

}

//4 - exam
func (r ExamRepo) Total(ctx context.Context, req *models.Id) (res models.Dis, err error) {

	query := `select order_id, sum(list_price) AS "list_price" , sum(discount) AS "discount"
from order_items
WHERE order_id = $1 GROUP BY  order_id`

	err = r.db.QueryRow(ctx, query, req.Order_id).Scan(
		&res.Order_id,
		&res.List_price,
		&res.Discount,
	)

	if err != nil {
		return res, err
	}

	if req.Promo_Code == "" {
		return res, nil
	}

	res.List_price -= res.Discount

	return res, nil
}

// 3-exam
func (r ExamRepo) Create(ctx context.Context, req *models.CreatePromo) (int, error) {

	query := `INSERT INTO promo("id","name",
	  "discount",
	  "discount_type",
	  "order_limit_price"
	  ) 
	  VALUES((SELECT MAX(id) + 1 FROM promocode),$1,$2,$3,$4) RETURNING id`
	id := 0
	err := r.db.QueryRow(ctx, query, req.Name, req.Discount, req.Type, req.Limitt).Scan(&id)
	if err != nil {
	  return 0, err
	}
	return id, nil
  }
  
  func (r ExamRepo) GetByID(ctx context.Context, req *models.PromocodePrimaryKey) (*models.Promocode, error) {
	var (
	  query     string
	  promocode models.Promocode
	)
  
	query = `SELECT * FROM promocode WHERE id =$1`
  
	err := r.db.QueryRow(ctx, query, req.PromocodeId).Scan(
	  &promocode.Id, &promocode.Name, &promocode.Discount, &promocode.DiscountType, &promocode.OrderLimitPrice)
	if err != nil {
	  return nil, err
	}
  
	return &promocode, nil
  }
  
  func (r ExamRepo) GetList(ctx context.Context, req *models.GetListBrandRequest) (resp *models.GetListPromocodeResponse, err error) {
  
	resp = &models.GetListPromocodeResponse{}
  
	var (
	  query  string
	  filter = " WHERE TRUE "
	  offset = " OFFSET 0"
	  limit  = " LIMIT 10"
	)
  
	query = `
	  SELECT
		COUNT(*) OVER(),
		id, 
		name, 
		discount,
		discount_type,
		order_limit_price
	  FROM promocode
	`
  
	if len(req.Search) > 0 {
	  filter += " AND name ILIKE '%'  '" + req.Search + "'  '%' "
	}
  
	if req.Offset > 0 {
	  offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}
  
	if req.Limit > 0 {
	  limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}
  
	query += filter + offset + limit
  
	rows, err := r.db.Query(ctx, query)
	if err != nil {
	  return nil, err
	}
	defer rows.Close()
  
	for rows.Next() {
	  var promocode models.Promocode
	  err = rows.Scan(
		&resp.Count,
		&promocode.Id,
		&promocode.Name,
		&promocode.Discount,
		&promocode.DiscountType,
		&promocode.OrderLimitPrice,
	  )
	  if err != nil {
		return nil, err
	  }
  
	  resp.Promocodes = append(resp.Promocodes, &promocode)
	}
  
	return resp, nil
  }
  
  func (r ExamRepo) Delete(ctx context.Context, req *models.PromocodePrimaryKey) (int64, error) {
	query := `
	  DELETE 
	  FROM promocode
	  WHERE id = $1
	`
  
	result, err := r.db.Exec(ctx, query, req.PromocodeId)
	if err != nil {
	  return 0, err
	}
	return result.RowsAffected(), nil
  }

//1-exam
// func (r stockRepo) SendStock(ctx context.Context, req *models.SendStockRequest) error {
// 	var (
// 		query  string
// 		params map[string]interface{}
// 	)

// 	sourceStock, err := r.GetStock(ctx, req.SourceStockID, req.ProductID)
// 	if err != nil {
// 		return err
// 	}
// 	if sourceStock.Quantity < req.Quantity {
// 		return errors.New("insufficient quantity in the source stock")
// 	}

// 	sourceStock.Quantity -= req.Quantity
// 	if sourceStock.Quantity == 0 {
// 		err = r.DeleteStock(ctx, sourceStock.ID)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		query = `
// 		UPDATE
// 		stocks
// 		SET
// 			quantity = :quantity
// 		WHERE id = :id
// 	`
// 		params = map[string]interface{}{
// 			"id":       sourceStock.ID,
// 			"quantity": sourceStock.Quantity,
// 		}
// 		query, args := helper.ReplaceQueryParams(query, params)
// 		_, err = r.db.Exec(ctx, query, args...)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	// Add the quantity to the destination stock
// 	destinationStock, err := r.GetStock(ctx, req.DestinationStockID, req.ProductID)
// 	if err != nil {
// 		destinationStock = &models.Stock{
// 			StoreID:   req.DestinationStockID,
// 			ProductID: req.ProductID,
// 			Quantity:  req.Quantity,
// 		}
// 		err = r.InsertStock(ctx, destinationStock)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		destinationStock.Quantity += req.Quantity
// 		query = `
// 		UPDATE
// 		stocks
// 		SET
// 			quantity = :quantity
// 		WHERE id = :id
// 	`
// 		params = map[string]interface{}{
// 			"id":       destinationStock.ID,
// 			"quantity": destinationStock.Quantity,
// 		}
// 		query, args := helper.ReplaceQueryParams(query, params)
// 		_, err = r.db.Exec(ctx, query, args...)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
// func (r *productRepo) GetStocks(ctx context.Context, storeID int) ([]*models.Stock, error) {
// 	var stocks []*models.Stock

// 	query := `
// 		SELECT
// 			s.product_id,
// 			p.product_name,
// 			s.quantity,
// 			p.list_price
// 		FROM
// 			stocks s
// 			INNER JOIN products p ON s.product_id = p.product_id
// 		WHERE
// 			s.store_id = $1
// 	`

// 	rows, err := r.db.Query(ctx, query, storeID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		stock := &models.Stock{}
// 		err := rows.Scan(
// 			&stock.ProductID,
// 			&stock.ProductName,
// 			&stock.Quantity,
// 			&stock.ListPrice,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		stocks = append(stocks, stock)
// 	}

// 	return stocks, nil
// }
// func (r *stockRepo) DeleteStock(ctx context.Context, storeID int64, productID int64) error {
// 	const query = `
// 		DELETE FROM stocks
// 		WHERE store_id = $1 AND product_id = $2
// 	`
// 	_, err := r.db.ExecContext(ctx, query, storeID, productID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }