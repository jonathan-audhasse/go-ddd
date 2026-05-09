package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ToPgUUID converts uuid.UUID to pgtype.UUID
func ToPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// FromPgUUID converts pgtype UUID into uuid.UUID
func FromPgUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	return uuid.UUID(p.Bytes)
}

// FromPgTimestamptz converts from pgtype.Timestamptz to time.Time
func FromPgTimestamptz(ts pgtype.Timestamptz) (time.Time, bool) {
	if !ts.Valid {
		return time.Time{}, false
	}
	return ts.Time, true
}

// ToPgTimestamptz converts from time.Time to pgtype.Timestamptz
func ToPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
