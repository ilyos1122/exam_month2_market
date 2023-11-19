package postgres

import (
	"database/sql"
	"fmt"
	"init/config"
	"init/storage"

	_ "github.com/lib/pq"
)

type Store struct {
	db       *sql.DB
	category storage.CategoryRepoI
	product  storage.ProductRepoI
	client storage.ClientRepoI
	branch storage.BranchRepoI
	order storage.OrderRepoI
	orderproduct storage.OrderProductRepoI
}

func NewConnectionPostgres(cfg *config.Config) (storage.StorageI, error) {

	connect := fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresDatabase,
		cfg.PostgresPassword,
		cfg.PostgresPort,
	)

	db, err := sql.Open("postgres", connect)
	if err != nil {
		panic(err)
	}

	return &Store{
		db: db,
	}, nil
}

func (s *Store) Category() storage.CategoryRepoI {

	if s.category == nil {
		s.category = NewCategoryRepo(s.db)
	}

	return s.category
}

func (s *Store) Product() storage.ProductRepoI {

	if s.product == nil {
		s.product = NewProductRepo(s.db)
	}

	return s.product
}


func (s *Store) Client() storage.ClientRepoI{
	if s.client == nil {
		s.client = NewClientRepo(s.db)
	}
	return s.client
}


func (s *Store) Branch() storage.BranchRepoI{
	if s.branch == nil {
		s.branch = NewBranchRepo(s.db)
	}
	return s.branch
}



func (s *Store) Order() storage.OrderRepoI{
	if s.order == nil {
		s.order = NewOrderRepo(s.db)
	}
	return s.order
}

func (s *Store) OrderProduct() storage.OrderProductRepoI{
	if s.orderproduct == nil {
		s.orderproduct = NewOrderProductRepo(s.db)
	}
	return s.orderproduct
}