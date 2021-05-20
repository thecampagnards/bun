package bun

import (
	"context"
	"database/sql"
	"reflect"
)

type scanModel struct {
	hookStubs

	db *DB

	dest      []interface{}
	scanIndex int
}

var _ model = (*scanModel)(nil)

func newScanModel(db *DB, dest []interface{}) *scanModel {
	return &scanModel{
		db:   db,
		dest: dest,
	}
}

func (m *scanModel) ScanRows(ctx context.Context, rows *sql.Rows) (int, error) {
	if !rows.Next() {
		return 0, errNoRows(rows)
	}

	dest := makeDest(m, len(m.dest))

	m.scanIndex = 0
	if err := rows.Scan(dest...); err != nil {
		return 0, err
	}

	return 1, nil
}

func (m *scanModel) ScanRow(ctx context.Context, rows *sql.Rows) error {
	return rows.Scan(m.dest...)
}

func (m *scanModel) Scan(src interface{}) error {
	dest := reflect.ValueOf(m.dest[m.scanIndex])
	m.scanIndex++

	scanner := m.db.dialect.Scanner(dest.Type())
	return scanner(dest, src)
}
