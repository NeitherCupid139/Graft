This directory is the plugin-owned migration boundary for `server/plugins/user`.

Historical mixed migrations remain under `server/internal/ent/migrate/migrations` until the
`users` / `refresh_sessions` / `user_roles` ownership history can be split without rewriting
the live Atlas chain. New user-only migration versions should land here.

`202605190001_user_plugin_boundary_checkpoint.sql` is the first forward-only checkpoint in this
directory. It is intentionally a no-op baseline marker so the plugin-owned migration chain becomes
runnable without claiming ownership of the historical mixed Atlas files.
