package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	res, err := s.db.Exec("INSERT into parcel (client, status, address, created_at) VALUES (:client, :status, :address, :createdAt)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("createdAt", p.CreatedAt),
	)

	if err != nil {
		return 0, err
	}

	lastId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	row := s.db.QueryRow("SELECT * FROM parcel WHERE number=:number", sql.Named("number", number))
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client=:clientId", sql.Named("clientId", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		parcel := Parcel{}

		err := rows.Scan(&parcel.Number, &parcel.Client, &parcel.Status, &parcel.Address, &parcel.CreatedAt)
		if err != nil {
			return res, err
		}

		res = append(res, parcel)
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	_, err := s.db.Exec("UPDATE parcel SET status=:status WHERE number=:number",
		sql.Named("status", status),
		sql.Named("number", number),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow("SELECT status FROM parcel WHERE number=:number", sql.Named("number", number))
	var p string
	err = row.Scan(&p)
	if err != nil {
		tx.Rollback()
		return err
	}

	if p == ParcelStatusRegistered {
		tx.Exec("UPDATE parcel SET address=:address WHERE number=:number",
			sql.Named("address", address),
			sql.Named("number", number),
		)
	}

	tx.Commit()
	return nil
}

func (s ParcelStore) Delete(number int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow("SELECT status FROM parcel WHERE number=:number", sql.Named("number", number))
	var p string
	err = row.Scan(&p)
	if err != nil {
		tx.Rollback()
		return err
	}

	if p == ParcelStatusRegistered {
		tx.Exec("DELETE FROM parcel WHERE number=:number", sql.Named("number", number))
	}

	tx.Commit()
	return nil
}
