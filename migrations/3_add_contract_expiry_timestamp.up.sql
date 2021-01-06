ALTER TABLE contract ADD expiry_timestamp REAL NOT NULL DEFAULT 0;
UPDATE contract SET expiry_timestamp = contract_expiry_timestamp(text_id, props);
