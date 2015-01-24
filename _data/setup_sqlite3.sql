--
-- Database structure for SQLite3
--

DROP TABLE IF EXISTS `cachedfetch_cache`;

CREATE TABLE IF NOT EXISTS `cachedfetch_cache` (
	--
	-- context and fetch information
	--
	`url`               VARCHAR(255) DEFAULT '',
	`context_str`       VARCHAR(255) DEFAULT '',
	`context_time`      INT(11) DEFAULT 0,
	`fetched_time`      INT(11) DEFAULT 0,

	--
	-- response meta information
	--
	`status`            TEXT DEFAULT '',
	`status_code`       INT(5) DEFAULT 200,
	`proto`             TEXT DEFAULT '',
	`content_length`    INT(11) DEFAULT 0,
	`transfer_encoding` TEXT DEFAULT '',
	`header`            TEXT DEFAULT '',
	`trailer`           TEXT DEFAULT '',
	`request`           TEXT DEFAULT '',
	`tls`               TEXT DEFAULT '',

	--
	-- response body
	--
	`body`              BLOB,

	PRIMARY KEY (`url`, `context_str`, `context_time`)
);

--
-- Add extra index
--
CREATE INDEX `url` ON `cachedfetch_cache`(`url`);
CREATE INDEX `context`  ON `cachedfetch_cache`(`context_str`, `context_time`);
CREATE INDEX `context_str`  ON `cachedfetch_cache`(`context_str`);
CREATE INDEX `context_time` ON `cachedfetch_cache`(`context_time`);
