use coprocessor;
SET time_zone='+00:00';

-- group by 单列索引
-- count
select int_idx,count(*) from t1 group by int_idx order by int_idx;
select str_idx1_1,count(*) from t1 group by str_idx1_1 order by str_idx1_1;

-- 2. sum/agv 与 Hash 重复，不测

-- group by 多列索引

-- count
select str_idx1_1,str_idx1_2,count(*) from t1 group by str_idx1_1,str_idx1_2 order by str_idx1_1,str_idx1_2;
SELECT int_idx, COUNT(*) FROM  t1 GROUP BY int_idx LIMIT 20;

-- 2. sum/avg
-- BYTES:
SELECT COUNT(*), AVG(str_idx1_1) FROM t1 GROUP BY str_idx1_1, str_idx1_2;
SELECT COUNT(*), SUM(str_idx1_1) FROM t1 GROUP BY str_idx1_1, str_idx1_2;
-- INT:
SELECT COUNT(*), SUM(str_idx1_1) FROM t1 GROUP BY str_idx1_1, str_idx1_2;