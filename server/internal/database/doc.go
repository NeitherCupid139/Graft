// Package database opens the PostgreSQL-backed GORM connection for the core runtime.
//
// Database ownership stays in core so plugins can depend on a stable runtime
// handle instead of constructing their own storage connections.
package database
