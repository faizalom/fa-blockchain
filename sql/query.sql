CREATE TABLE IF NOT EXISTS utxos (
    block_height INT NOT NULL DEFAULT 0,
	tx_id VARCHAR(128) NOT NULL,
	out_index INT NOT NULL,
	amount INT NOT NULL,
	pub_key_hash VARCHAR(256) NOT NULL,
	status VARCHAR(10) NOT NULL DEFAULT 'Unspent',
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (tx_id, out_index)
);

INSERT INTO `utxos` (`block_height`, `tx_id`, `out_index`, `amount`, `pub_key_hash`, `status`, `created_at`) VALUES
(1, '6ac66397f8f1f6c43a3f338eaa21d1fc49cb799f1d85ac26eb6a14d3438d8959', 0, 100, '659d7a7a33e7531c52944f4ecffe3129da5e3668', 'Unspent', '2025-01-01 00:00:00');
