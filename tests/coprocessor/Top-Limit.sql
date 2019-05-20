use coprocessor;
SET time_zone='+00:00';

SELECT * FROM t1 LIMIT 10;
SELECT * FROM t1 ORDER BY t1.id LIMIT 10;
SELECT col_int FROM t1 ORDER BY col_int LIMIT 23;
SELECT * FROM t1 LIMIT 2000;
SELECT * FROM t1 order by id limit 10;
SELECT t1.col_int FROM t1 ORDER BY t1.col_int DESC LIMIT 5;
SELECT t2.id FROM t2 ORDER BY t2.id DESC LIMIT 10;
SELECT id FROM t1 ORDER BY t1.col_uint LIMIT 10;

SELECT t1.col_datetime FROM t1 ORDER BY t1.col_datetime LIMIT  10;
SELECT t1.col_datetime FROM t1 ORDER BY t1.col_datetime DESC LIMIT 10;
SELECT t1.col_date FROM t1 ORDER BY t1.col_date LIMIT  10;
SELECT t1.col_date FROM t1 ORDER BY t1.col_date DESC LIMIT  10;

SELECT t1.col_dec FROM t1 ORDER BY t1.col_dec LIMIT  10;
SELECT t1.col_dec FROM t1 ORDER BY t1.col_dec DESC LIMIT  10;
-- Duration:
SELECT t1.col_duration FROM  t1 ORDER BY t1.col_duration LIMIT 10;
SELECT t1.col_duration FROM  t1 ORDER BY t1.col_duration DESC LIMIT 10;
-- BYTES
SELECT t1.id FROM t1 WHERE str_unq_idx IS NULL LIMIT 30;
SELECT str_unq_idx FROM t1 WHERE str_unq_idx IS NULL ORDER BY str_unq_idx LIMIT 30;
