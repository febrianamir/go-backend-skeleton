package repository

import (
	"app/config"
	"app/lib"
	"app/lib/cache"
	"app/lib/logger"
	"app/lib/mailer"
	"app/lib/task"
	"context"
)

type TrxKey struct{}

type Repository struct {
	config    *config.Config
	db        *lib.Database
	mailer    *mailer.SMTP
	publisher *task.Publisher
	cache     *cache.Cache
}

func NewRepository(config *config.Config, db *lib.Database, mailer *mailer.SMTP, publisher *task.Publisher, cache *cache.Cache) Repository {
	return Repository{
		config:    config,
		db:        db,
		mailer:    mailer,
		publisher: publisher,
		cache:     cache,
	}
}

func (repo *Repository) Transaction(ctx context.Context, fn func(context.Context) error) error {
	trx := repo.db.Begin()

	ctx = context.WithValue(ctx, TrxKey{}, &lib.Database{DB: trx})
	if err := fn(ctx); err != nil {
		trx.Rollback()
		return err
	}

	return trx.Commit().Error
}

// prepareDBWithContext prepares a database connection with proper context setup for logging and tracing.
// It performs three main operations:
// 1. Adds the operation name to context for SQL logging identification
// 2. Retrieves the database connection from context (transaction) or falls back to default
// 3. Configures the database connection to use the enriched context
//
// Parameters:
//   - ctx: The incoming request context that may contain transaction or other values
//   - operation: The name of the database operation (e.g., "GetUser", "CreateOrder")
//     used for logging and debugging purposes
//
// Returns:
//   - context.Context: The enriched context with operation name added
//   - *lib.Database: The configured database connection ready for queries
//
// Usage example:
//
//	ctx, tx := repo.prepareDBWithContext(ctx, "GetUser")
//	err := tx.Where("id = ?", userID).First(&user).Error
func (repo *Repository) prepareDBWithContext(ctx context.Context, operation string) (context.Context, *lib.Database) {
	ctx = context.WithValue(ctx, logger.CtxRepoName, operation)
	tx, ok := ctx.Value(TrxKey{}).(*lib.Database)
	if !ok {
		tx = repo.db
	}
	tx.DB = tx.WithContext(ctx)
	return ctx, tx
}
