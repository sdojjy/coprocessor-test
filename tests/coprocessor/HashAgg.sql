use coprocessor;
SET time_zone='+00:00';

-- count
select col_uint, count(*) from t1 group by col_uint order by col_uint;
-- sum
select col_uint, sum(col_real) from t1 group by col_uint order by col_uint;
-- UInt64:
SELECT col_uint, SUM(col_uint) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- Int64:
select col_uint, sum(col_int) from t1 group by col_uint order by col_uint LIMIT 20;
-- Real:
select col_uint, sum(col_real) from t1 group by col_uint order by col_uint;
-- Decimal:
select col_uint, sum(col_dec) from t1 group by col_uint order by col_uint LIMIT 20;
-- BYTES:
select col_uint, sum(col_dec) from t1 group by col_uint order by col_uint LIMIT 20;
-- DateTime ignore these cases, see tidb issue: https://github.com/pingcap/tidb/issues/10543
-- select col_uint, sum(col_datetime) from t1 group by col_uint order by col_uint LIMIT 20;
-- Date:
-- select col_uint, sum(col_date) from t1 group by col_uint order by col_uint LIMIT 20;
-- Duration:
-- select col_uint, sum(t1.col_duration ) from t1 group by col_uint order by col_uint LIMIT 20;

-- avg
-- REAL:
select col_uint,avg(col_real) from t1 group by col_uint order by col_uint;
-- INT64
SELECT col_uint, AVG(col_int) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- UINT64
SELECT col_uint, AVG(col_uint) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- DECIMAL
SELECT col_uint, AVG(col_dec) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- BYTES:
SELECT col_uint, AVG(str_unq_idx) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- DateTime
-- SELECT col_uint, AVG(col_datetime) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- Date:
-- SELECT col_uint, AVG(col_date) FROM t1 GROUP BY col_uint ORDER BY col_uint LIMIT 20;
-- Duration:
-- SELECT col_uint, AVG(t1.col_duration ) from t1 group by col_uint order by col_uint LIMIT 20;



-- group by 多列: 对上方测试添加 `GROUP BY col_uint, col_int`
-- count
select col_uint, count(*) from t1 group by col_uint, col_int order by col_uint;
-- 2. sum/avg TODO
--	a. REAL:
SELECT col_uint,AVG(col_real) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint;
SELECT col_uint,SUM(col_real) FROM t1 GROUP BY col_uint, col_int ORDER BY  col_uint;
-- INT64
SELECT col_uint, AVG(col_int) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
SELECT col_uint, SUM(col_int) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- UINT64
SELECT col_uint, AVG(col_uint) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
SELECT col_uint, SUM(col_uint) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- DECIMAL
SELECT col_uint, AVG(col_dec) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
SELECT col_uint, SUM(col_dec) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- BYTES:
SELECT col_uint, AVG(str_unq_idx) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
SELECT col_uint, SUM(str_unq_idx) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- DateTime
-- SELECT col_uint, AVG(col_datetime) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- SELECT col_uint, SUM(col_datetime) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- Date:
-- SELECT col_uint, AVG(col_date) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;
-- SELECT col_uint, SUM(col_date) FROM t1 GROUP BY col_uint, col_int ORDER BY col_uint LIMIT 20;

-- Duration:
-- SELECT col_uint, AVG(t1.col_duration ) from t1 group by col_uint, col_int order by col_uint LIMIT 20;
-- SELECT col_uint, AVG(t1.col_duration ) from t1 group by col_uint, col_int order by col_uint LIMIT 20;

