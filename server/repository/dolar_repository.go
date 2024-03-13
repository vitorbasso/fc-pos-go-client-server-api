package repository

import (
	"context"
	"database/sql"
	"time"
)

type DolarRepository struct {
	db *sql.DB
}

func NewDolarRepository(db *sql.DB) *DolarRepository {
	return &DolarRepository{db: db}
}

type USDBRL struct {
	Code       string `json:"code"`
	CodeIn     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func (d *DolarRepository) InsertDolar(dolar *USDBRL) error {
	stmt, err := d.db.Prepare("INSERT INTO usd_brl_quotations (code, codein, name, high, low, var_bid, pct_change, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	insertContext, insertCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer insertCancel()
	_, err = stmt.ExecContext(insertContext, dolar.Code, dolar.CodeIn, dolar.Name, dolar.High, dolar.Low, dolar.VarBid, dolar.PctChange, dolar.Bid, dolar.Ask, dolar.Timestamp, dolar.CreateDate)
	if err != nil {
		return err
	}
	return nil
}
