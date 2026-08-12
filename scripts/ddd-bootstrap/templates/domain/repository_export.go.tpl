package repository

import (
	"context"

	"gorm.io/gorm"
)

func New{{.TablePascal}}Dao(db *gorm.DB) {
	Q = Use(db)
}

func Get{{.TablePascal}}Dao(ctx context.Context, tx *Query) I{{.TablePascal}}Do {
	if tx != nil {
		return tx.{{.TablePascal}}.WithContext(ctx)
	}
	return Q.{{.TablePascal}}.WithContext(ctx)
}
