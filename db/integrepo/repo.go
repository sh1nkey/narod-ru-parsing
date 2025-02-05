package integrepo

import (
	"context"
	"data-sender/core"
	"data-sender/core/parsenarod"
	"data-sender/db"
	"github.com/pkg/errors"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog/log"
)

type dbRepo struct {
	conn core.Conn
}


func NewPostgresRepo(conn core.Conn) parsenarod.Repository {
	log.Info().Msg("creating user repository...")
	return &dbRepo{
		conn: conn,
	}
}


func (d *dbRepo) Create(ctx context.Context, url *parsenarod.CreateUrlReqDTO, txs ...core.UpdateOptions) error {
	tx := db.GetUpdateOptions(d.conn, txs...)

	_, err := tx.Exec(ctx, `
		INSERT INTO urls (url)
                      VALUES ($1);`,
					  url.Url,
		)
	if err != nil {
		return err
	}
	return nil
}




func (d *dbRepo) GetAll(ctx context.Context, limit int, offset int, options ...core.QueryOptions) ([]parsenarod.Url, error) {

	tx, _ := db.GetQueryOptions(d.conn, options...)

	urls := make([]parsenarod.Url, 0)
	log.Debug().
		Int("limit", limit).
		Int("offset", offset).
		Msg("получаем ссылки с параметрами")
		
	rows, err := tx.Query(ctx,
		`SELECT urls.id, urls.url, urls.description, urls.is_empty, urls.created_at FROM urls ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		if err == pgx.ErrNoRows {
			return urls, errors.WithStack(core.ErrNotFound)
		}
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	for rows.Next() {
		url := parsenarod.Url{}
		err = rows.Scan(&url.Id, &url.Url, &url.Description, &url.IsEmpty, &url.CreatedAt)
		if err != nil {
			if err == pgx.ErrNoRows {
				return nil, errors.WithStack(core.ErrNotFound)
			}
			return nil, errors.WithStack(err)
		}
		urls = append(urls, url)
	}

	return urls, nil
}



func (d *dbRepo) MarkAsEmpty(ctx context.Context, id uint64, options ...core.UpdateOptions) error {
	tx := db.GetUpdateOptions(d.conn, options...)

	_, err := tx.Exec(ctx, `UPDATE urls SET is_empty = true WHERE urls.id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *dbRepo) SetDescription(ctx context.Context, id uint64, description string, options ...core.UpdateOptions) error {
	tx := db.GetUpdateOptions(d.conn, options...)

	_, err := tx.Exec(ctx, `UPDATE urls SET description = $1, is_empty = false WHERE urls.id = $2`, description, id)
	if err != nil {
		return err
	}
	return nil
}