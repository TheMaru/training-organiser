package api

import "database/sql"

func stringToPointer(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}

	return sql.NullString{String: *s, Valid: true}
}
