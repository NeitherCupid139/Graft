-- system_config_values stores administrator overrides only.
-- ConfigDefinition metadata and defaults remain module-registered runtime authority.

CREATE TABLE IF NOT EXISTS system_config_values (
    key TEXT PRIMARY KEY,
    override_value JSONB NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE system_config_values IS 'Administrator-provided system configuration overrides only.';
COMMENT ON COLUMN system_config_values.key IS 'Stable config definition key registered by a module.';
COMMENT ON COLUMN system_config_values.override_value IS 'Administrator override JSON. Module defaults are never copied here.';
COMMENT ON COLUMN system_config_values.updated_at IS 'Last override write timestamp.';
