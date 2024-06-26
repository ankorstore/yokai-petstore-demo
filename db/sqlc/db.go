// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package sqlc

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.createOwnerStmt, err = db.PrepareContext(ctx, createOwner); err != nil {
		return nil, fmt.Errorf("error preparing query CreateOwner: %w", err)
	}
	if q.createOwnerPetStmt, err = db.PrepareContext(ctx, createOwnerPet); err != nil {
		return nil, fmt.Errorf("error preparing query CreateOwnerPet: %w", err)
	}
	if q.deleteOwnerStmt, err = db.PrepareContext(ctx, deleteOwner); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteOwner: %w", err)
	}
	if q.deleteOwnerPetStmt, err = db.PrepareContext(ctx, deleteOwnerPet); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteOwnerPet: %w", err)
	}
	if q.getOwnerStmt, err = db.PrepareContext(ctx, getOwner); err != nil {
		return nil, fmt.Errorf("error preparing query GetOwner: %w", err)
	}
	if q.getOwnerPetStmt, err = db.PrepareContext(ctx, getOwnerPet); err != nil {
		return nil, fmt.Errorf("error preparing query GetOwnerPet: %w", err)
	}
	if q.listOwnerPetsStmt, err = db.PrepareContext(ctx, listOwnerPets); err != nil {
		return nil, fmt.Errorf("error preparing query ListOwnerPets: %w", err)
	}
	if q.listOwnersStmt, err = db.PrepareContext(ctx, listOwners); err != nil {
		return nil, fmt.Errorf("error preparing query ListOwners: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.createOwnerStmt != nil {
		if cerr := q.createOwnerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createOwnerStmt: %w", cerr)
		}
	}
	if q.createOwnerPetStmt != nil {
		if cerr := q.createOwnerPetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing createOwnerPetStmt: %w", cerr)
		}
	}
	if q.deleteOwnerStmt != nil {
		if cerr := q.deleteOwnerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteOwnerStmt: %w", cerr)
		}
	}
	if q.deleteOwnerPetStmt != nil {
		if cerr := q.deleteOwnerPetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteOwnerPetStmt: %w", cerr)
		}
	}
	if q.getOwnerStmt != nil {
		if cerr := q.getOwnerStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getOwnerStmt: %w", cerr)
		}
	}
	if q.getOwnerPetStmt != nil {
		if cerr := q.getOwnerPetStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getOwnerPetStmt: %w", cerr)
		}
	}
	if q.listOwnerPetsStmt != nil {
		if cerr := q.listOwnerPetsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listOwnerPetsStmt: %w", cerr)
		}
	}
	if q.listOwnersStmt != nil {
		if cerr := q.listOwnersStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listOwnersStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                 DBTX
	tx                 *sql.Tx
	createOwnerStmt    *sql.Stmt
	createOwnerPetStmt *sql.Stmt
	deleteOwnerStmt    *sql.Stmt
	deleteOwnerPetStmt *sql.Stmt
	getOwnerStmt       *sql.Stmt
	getOwnerPetStmt    *sql.Stmt
	listOwnerPetsStmt  *sql.Stmt
	listOwnersStmt     *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                 tx,
		tx:                 tx,
		createOwnerStmt:    q.createOwnerStmt,
		createOwnerPetStmt: q.createOwnerPetStmt,
		deleteOwnerStmt:    q.deleteOwnerStmt,
		deleteOwnerPetStmt: q.deleteOwnerPetStmt,
		getOwnerStmt:       q.getOwnerStmt,
		getOwnerPetStmt:    q.getOwnerPetStmt,
		listOwnerPetsStmt:  q.listOwnerPetsStmt,
		listOwnersStmt:     q.listOwnersStmt,
	}
}
