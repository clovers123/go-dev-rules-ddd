package {{.Domain}}

import "gorm.io/gen"

type {{.TablePascal}} interface {
	// FindPage
	//
	// SELECT * FROM @@table
	// WHERE owned_by=@ownedBy AND is_deleted=false
	// AND (@keyword = '' OR name ILIKE @keyword)
	// ORDER BY create_time DESC LIMIT @limit OFFSET @offset
	FindPage(ownedBy string, keyword string, limit int, offset int) ([]*gen.T, error)
}
