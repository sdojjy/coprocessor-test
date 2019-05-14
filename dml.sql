CREATE TABLE `t1` (
 `id` int8 NOT NULL AUTO_INCREMENT,
 `str_idx1_1` varchar(255) NOT NULL,
 `str_idx1_2` varchar(255) DEFAULT NULL,
 `int_idx` int(11) DEFAULT NULL,
 `str_unq_idx` varchar(255) DEFAULT NULL,
 `col_real` double DEFAULT NULL,
 `col_dec` decimal(40, 10),
 `col_duration` time(4) Comment 'Duration',
 `col_date` date,
 `col_datetime` datetime DEFAULT CURRENT_TIMESTAMP,
 `col_ts` timestamp(4),
 `col_int` int8 DEFAULT 0,
 `col_uint` int8 unsigned DEFAULT null,
 PRIMARY KEY (`id`) Comment 'int 主键',
 KEY `unified_idx` (`str_idx1_1`,`str_idx1_2`) Comment '联合索引',
 KEY `int_idx` (`int_idx`) Comment '单列索引',
 UNIQUE KEY `unq_idx` (`str_unq_idx`) Comment '唯一索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

Create table `t2`(
    `id` int8 unsigned NOT NULL,
    `col_json` json,
    primary key(id)
);