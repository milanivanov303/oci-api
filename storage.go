package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccounts() ([]*Account, error)
	GetAccountByID(int) (*Account, error)
}

type MysqlStore struct {
	db *sql.DB
}

func NewMysqlStore(configuration Configuration) (*MysqlStore, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", configuration.DBUsername, configuration.DBPassword, configuration.DBHost, configuration.DBPort, configuration.DBDatabase)
	db, err := sql.Open("mysql", connStr)
    if err != nil {
        return nil, err
    }

	if err := db.Ping(); err != nil {
		return nil, err
	}

    return &MysqlStore{
		db: db,
	}, nil
}

func (s *MysqlStore) Init() error {
	return s.createAccountsTable()
}

func (s *MysqlStore) createAccountsTable() error {
	query := `create table if not exists accounts(
		id INT AUTO_INCREMENT,
		first_name varchar(50),
		last_name varchar(50),
		number BIGINT,
		balance INT,
		created_at TIMESTAMP,
		PRIMARY KEY(id)
	)`

	_, err := s.db.Exec(query)
	return err
}

func (s *MysqlStore) CreateAccount(acc *Account) error {
	query := `insert into accounts
	(first_name, last_name, number, balance, created_at)
	values(?, ?, ?, ?, ?)`

	resp, err := s.db.Query(
		query,
		acc.FirstName,
		acc.LastName,
		acc.Number,
		acc.Balance,
		acc.CreatedAt,
	)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *MysqlStore) UpdateAccount(*Account) error {
	return nil
}

func (s *MysqlStore) DeleteAccount(id int) error {
	return nil
}

func (s *MysqlStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}

func (s *MysqlStore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("select * from accounts")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)
		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
			&account.Number,
			&account.Balance,
			&account.CreatedAt,
		)

		if err != nil {
			return nil , err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

