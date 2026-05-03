package postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// toPgUUID converts uuid.UUID to pgtype.UUID
func toPgUUID(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: u,
		Valid: true,
	}
}

// fromPgUUID converts pgtype UUID into uuid.UUID
func fromPgUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	return uuid.UUID(p.Bytes)
}

// fromPgTimestamptz converts from pgtype.Timestamptz to time.Time
func fromPgTimestamptz(ts pgtype.Timestamptz) (time.Time, bool) {
	if !ts.Valid {
		return time.Time{}, false
	}
	return ts.Time, true
}

// toPgTimestamptz converts from time.Time to pgtype.Timestamptz
func toPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  t,
		Valid: true,
	}
}
