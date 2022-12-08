package utils

import "database/sql"

func NullString(s string) (ns sql.NullString) {
	if s != "" {
		ns.String = s
		ns.Valid = true
	}

	return ns
}

func NullFloat64(f float64) (ns sql.NullFloat64) {
	if f != 0 {
		ns.Float64 = f
		ns.Valid = true
	}
	return ns
}

func FormatNullTime(nt sql.NullTime, format string) string {
	if nt.Valid {
		return nt.Time.Format(format)
	}
	return ""
}