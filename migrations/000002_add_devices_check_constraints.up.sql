ALTER TABLE devices ADD CONSTRAINT devices_value_check CHECK (value >= 0);
