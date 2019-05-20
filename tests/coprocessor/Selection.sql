use coprocessor;
SET time_zone='+00:00';

-- int（包含了与 > 2^63 的值比较的情况）
select * from t1 where col_int<-3 order by id;
select * from t1 where col_int < -9223372036854775808 order by id;
select * from t1 where col_int < 9223372036854775807 order by id;
select * from t1 where col_int < 18446744073709551615 order by id;
select * from t1 where col_int < null order by id;
select * from t1 where col_int<=-3 order by id;
select * from t1 where col_int <= -9223372036854775808 order by id;
select * from t1 where col_int <= 9223372036854775807 order by id;
select * from t1 where col_int <= 18446744073709551615 order by id;
select * from t1 where col_int <= null order by id;
select * from t1 where col_int = -3 order by id;
select * from t1 where col_int = -9223372036854775808 order by id;
select * from t1 where col_int = 9223372036854775807 order by id;
select * from t1 where col_int = 18446744073709551615 order by id;
select * from t1 where col_int = null order by id;
select * from t1 where col_int != -3 order by id;
select * from t1 where col_int != -9223372036854775808 order by id;
select * from t1 where col_int != 9223372036854775807 order by id;
select * from t1 where col_int != 18446744073709551615 order by id;
select * from t1 where col_int != null order by id;
select * from t1 where col_int >= -3 order by id;
select * from t1 where col_int >= -9223372036854775808 order by id;
select * from t1 where col_int >= 9223372036854775807 order by id;
select * from t1 where col_int >= 18446744073709551615 order by id;
select * from t1 where col_int >= null order by id;
select * from t1 where col_int > -3 order by id;
select * from t1 where col_int > -9223372036854775808 order by id;
select * from t1 where col_int > 9223372036854775807 order by id;
select * from t1 where col_int > 18446744073709551615 order by id;
select * from t1 where col_int > null order by id;
select * from t1 where col_int <=> -3 order by id;
select * from t1 where col_int <=> -9223372036854775808 order by id;
select * from t1 where col_int <=> 9223372036854775807 order by id;
select * from t1 where col_int <=> 18446744073709551615 order by id;
select * from t1 where col_int <=> null order by id;
select * from t1 where col_int is null order by id;


-- uint（包含了与负数比较的情况）
select * from t1 where col_uint<498 order by id;
select * from t1 where col_uint < -1 order by id;
select * from t1 where col_uint < 0 order by id;
select * from t1 where col_uint < 18446744073709551615 order by id;
select * from t1 where col_uint < null order by id;
select * from t1 where col_uint<=498 order by id;
select * from t1 where col_uint <= -1 order by id;
select * from t1 where col_uint <= 0 order by id;
select * from t1 where col_uint <= 18446744073709551615 order by id;
select * from t1 where col_uint <= null order by id;
select * from t1 where col_uint = 498 order by id;
select * from t1 where col_uint = -1 order by id;
select * from t1 where col_uint = 0 order by id;
select * from t1 where col_uint = 18446744073709551615 order by id;
select * from t1 where col_uint = null order by id;
select * from t1 where col_uint != 498 order by id;
select * from t1 where col_uint != -1 order by id;
select * from t1 where col_uint != 0 order by id;
select * from t1 where col_uint != 18446744073709551615 order by id;
select * from t1 where col_uint != null order by id;
select * from t1 where col_uint >= 498 order by id;
select * from t1 where col_uint >= -1 order by id;
select * from t1 where col_uint >= 0 order by id;
select * from t1 where col_uint >= 18446744073709551615 order by id;
select * from t1 where col_uint >= null order by id;
select * from t1 where col_uint > 498 order by id;
select * from t1 where col_uint > -1 order by id;
select * from t1 where col_uint > 0 order by id;
select * from t1 where col_uint > 18446744073709551615 order by id;
select * from t1 where col_uint > null order by id;
select * from t1 where col_uint <=> 498 order by id;
select * from t1 where col_uint <=> -1 order by id;
select * from t1 where col_uint <=> 0 order by id;
select * from t1 where col_uint <=> 18446744073709551615 order by id;
select * from t1 where col_uint <=> null order by id;
select * from t1 where col_uint is null order by id;

-- int 和 uint 比较
select * from t1 where col_int < col_uint order by id;
select * from t1 where col_int <= col_uint order by id;
select * from t1 where col_int = col_uint order by id;
select * from t1 where col_int != col_uint order by id;
select * from t1 where col_int > col_uint order by id;
select * from t1 where col_int >= col_uint order by id;
select * from t1 where col_int <=> col_uint order by id;
select * from t1 where col_uint < col_int order by id;
select * from t1 where col_uint <= col_int order by id;
select * from t1 where col_uint = col_int order by id;
select * from t1 where col_uint != col_int order by id;
select * from t1 where col_uint > col_int order by id;
select * from t1 where col_uint >= col_int order by id;
select * from t1 where col_uint <=> col_int order by id;

-- real
select * from t1 where col_real< -0.670441 order by id;
select * from t1 where col_real< null order by id;
select * from t1 where col_real<= -0.670441 order by id;
select * from t1 where col_real<= null order by id;
select * from t1 where col_real = -0.670441 order by id;
select * from t1 where col_real = null order by id;
select * from t1 where col_real != -0.670441 order by id;
select * from t1 where col_real != null order by id;
select * from t1 where col_real > -0.670441 order by id;
select * from t1 where col_real > null order by id;
select * from t1 where col_real >= -0.670441 order by id;
select * from t1 where col_real >= null order by id;
select * from t1 where col_real <=> -0.670441 order by id;
select * from t1 where col_real <=> null order by id;
select * from t1 where col_real is null order by id;

-- decimal （out of range和truncate的情况未包含）
select * from t1 where col_dec <  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec < -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec < 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec < null order by id;
select * from t1 where col_dec <=  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec <= -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec <= 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec <= null order by id;
select * from t1 where col_dec =  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec = -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec = 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec = null order by id;
select * from t1 where col_dec !=  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec != -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec != 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec != null order by id;
select * from t1 where col_dec >  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec > -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec > 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec > null order by id;
select * from t1 where col_dec >=  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec >= -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec >= 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec >= null order by id;
select * from t1 where col_dec <=>  -973679598173940274292704421.5490621702 order by id;
select * from t1 where col_dec <=> -999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec <=> 999999999999999999999999999999.9999999999 order by id;
select * from t1 where col_dec <=> null order by id;
select * from t1 where col_dec is null order by id;

-- duration
select * from t1 where col_duration < "40:55:49.0000" order by id;
select * from t1 where col_duration < "-838:59:59" order by id;
select * from t1 where col_duration < "838:59:59" order by id;
select * from t1 where col_duration < null order by id;
select * from t1 where col_duration <= "40:55:49.0000" order by id;
select * from t1 where col_duration <= "-838:59:59" order by id;
select * from t1 where col_duration <= "838:59:59" order by id;
select * from t1 where col_duration <= null order by id;
select * from t1 where col_duration = "40:55:49.0000" order by id;
select * from t1 where col_duration = "-838:59:59" order by id;
select * from t1 where col_duration = "838:59:59" order by id;
select * from t1 where col_duration = null order by id;
select * from t1 where col_duration != "40:55:49.0000" order by id;
select * from t1 where col_duration != "-838:59:59" order by id;
select * from t1 where col_duration != "838:59:59" order by id;
select * from t1 where col_duration != null order by id;
select * from t1 where col_duration > "40:55:49.0000" order by id;
select * from t1 where col_duration > "-838:59:59" order by id;
select * from t1 where col_duration > "838:59:59" order by id;
select * from t1 where col_duration > null order by id;
select * from t1 where col_duration >= "40:55:49.0000" order by id;
select * from t1 where col_duration >= "-838:59:59" order by id;
select * from t1 where col_duration >= "838:59:59" order by id;
select * from t1 where col_duration >= null order by id;
select * from t1 where col_duration <=> "40:55:49.0000" order by id;
select * from t1 where col_duration <=> "-838:59:59" order by id;
select * from t1 where col_duration <=> "838:59:59" order by id;
select * from t1 where col_duration <=> null order by id;
select * from t1 where col_duration is null order by id;

-- date（包含小于1000年的非法情况，不包含月或日为0的情况（会报错））
select * from t1 where col_date < cast("2000-12-31" as date) order by id;
select * from t1 where col_date < cast("1000-01-01" as date) order by id;
select * from t1 where col_date < cast("200-01-01" as date) order by id;
select * from t1 where col_date < cast("9999-12-31" as date) order by id;
select * from t1 where col_date < null order by id;
select * from t1 where col_date <= cast("2000-12-31" as date) order by id;
select * from t1 where col_date <= cast("1000-01-01" as date) order by id;
select * from t1 where col_date <= cast("200-01-01" as date) order by id;
select * from t1 where col_date <= cast("9999-12-31" as date) order by id;
select * from t1 where col_date <= null order by id;
select * from t1 where col_date = cast("2000-12-31" as date) order by id;
select * from t1 where col_date = cast("1000-01-01" as date) order by id;
select * from t1 where col_date = cast("200-01-01" as date) order by id;
select * from t1 where col_date = cast("9999-12-31" as date) order by id;
select * from t1 where col_date = null order by id;
select * from t1 where col_date != cast("2000-12-31" as date) order by id;
select * from t1 where col_date != cast("1000-01-01" as date) order by id;
select * from t1 where col_date != cast("200-01-01" as date) order by id;
select * from t1 where col_date != cast("9999-12-31" as date) order by id;
select * from t1 where col_date != null order by id;
select * from t1 where col_date > cast("2000-12-31" as date) order by id;
select * from t1 where col_date > cast("1000-01-01" as date) order by id;
select * from t1 where col_date > cast("200-01-01" as date) order by id;
select * from t1 where col_date > cast("9999-12-31" as date) order by id;
select * from t1 where col_date > null order by id;
select * from t1 where col_date >= cast("2000-12-31" as date) order by id;
select * from t1 where col_date >= cast("1000-01-01" as date) order by id;
select * from t1 where col_date >= cast("200-01-01" as date) order by id;
select * from t1 where col_date >= cast("9999-12-31" as date) order by id;
select * from t1 where col_date >= null order by id;
select * from t1 where col_date <=> cast("2000-12-31" as date) order by id;
select * from t1 where col_date <=> cast("1000-01-01" as date) order by id;
select * from t1 where col_date <=> cast("200-01-01" as date) order by id;
select * from t1 where col_date <=> cast("9999-12-31" as date) order by id;
select * from t1 where col_date <=> null order by id;
select * from t1 where col_date is null order by id;


-- datetime（包含小于1000年的非法情况，不包含月或日为0的情况（会报错））
select * from t1 where col_datetime < cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime < cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime < cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime < cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime < null order by id;
select * from t1 where col_datetime <= cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime <= cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime <= cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime <= cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime <= null order by id;
select * from t1 where col_datetime = cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime = cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime = cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime = cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime = null order by id;
select * from t1 where col_datetime != cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime != cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime != cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime != cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime != null order by id;
select * from t1 where col_datetime > cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime > cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime > cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime > cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime > null order by id;
select * from t1 where col_datetime >= cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime >= cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime >= cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime >= cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime >= null order by id;
select * from t1 where col_datetime <=> cast("1993-11-01 18:44:23" as datetime) order by id;
select * from t1 where col_datetime <=> cast('1000-01-01 00:00:00' as datetime) order by id;
select * from t1 where col_datetime <=> cast('9999-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime <=> cast('500-12-31 23:59:59' as datetime) order by id;
select * from t1 where col_datetime <=> null order by id;
select * from t1 where col_datetime is null order by id;

-- timestamp（包含了小于epoch time和大于2038年的非法时间）
select * from t1 where col_ts < "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts < "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts < '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts < "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts < "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts < null order by id;
select * from t1 where col_ts <= "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts <= "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts <= '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts <= "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts <= "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts <= null order by id;
select * from t1 where col_ts = "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts = "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts = '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts = "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts = "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts = null order by id;
select * from t1 where col_ts != "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts != "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts != '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts != "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts != "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts != null order by id;
select * from t1 where col_ts > "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts > "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts > '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts > "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts > "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts > null order by id;
select * from t1 where col_ts >= "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts >= "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts >= '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts >= "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts >= "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts >= null order by id;
select * from t1 where col_ts <=> "1989-05-15 19:12:50.9080" order by id;
select * from t1 where col_ts <=> "1970-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts <=> '2038-01-19 11:14:07.999999' order by id;
select * from t1 where col_ts <=> "1898-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts <=> "2100-01-01 08:00:01.000000" order by id;
select * from t1 where col_ts <=> null order by id;
select * from t1 where col_ts is null order by id;


-- date 和 datetime 比较
select * from t1 where col_date < col_datetime order by id;
select * from t1 where col_date <= col_datetime order by id;
select * from t1 where col_date = col_datetime order by id;
select * from t1 where col_date != col_datetime order by id;
select * from t1 where col_date > col_datetime order by id;
select * from t1 where col_date >= col_datetime order by id;
select * from t1 where col_date <=> col_datetime order by id;
select * from t1 where col_datetime < col_date order by id;
select * from t1 where col_datetime <= col_date order by id;
select * from t1 where col_datetime = col_date order by id;
select * from t1 where col_datetime != col_date order by id;
select * from t1 where col_datetime > col_date order by id;
select * from t1 where col_datetime >= col_date order by id;
select * from t1 where col_datetime <=> col_date order by id;

-- date 和 timestamp 比较
select * from t1 where col_date < col_ts order by id;
select * from t1 where col_date <= col_ts order by id;
select * from t1 where col_date = col_ts order by id;
select * from t1 where col_date != col_ts order by id;
select * from t1 where col_date > col_ts order by id;
select * from t1 where col_date >= col_ts order by id;
select * from t1 where col_date <=> col_ts order by id;
select * from t1 where col_ts < col_date order by id;
select * from t1 where col_ts <= col_date order by id;
select * from t1 where col_ts = col_date order by id;
select * from t1 where col_ts != col_date order by id;
select * from t1 where col_ts > col_date order by id;
select * from t1 where col_ts >= col_date order by id;
select * from t1 where col_ts <=> col_date order by id;

-- datetime 和 timestamp 比较
select * from t1 where col_datetime < col_ts order by id;
select * from t1 where col_datetime <= col_ts order by id;
select * from t1 where col_datetime = col_ts order by id;
select * from t1 where col_datetime != col_ts order by id;
select * from t1 where col_datetime > col_ts order by id;
select * from t1 where col_datetime >= col_ts order by id;
select * from t1 where col_datetime <=> col_ts order by id;
select * from t1 where col_ts < col_datetime order by id;
select * from t1 where col_ts <= col_datetime order by id;
select * from t1 where col_ts = col_datetime order by id;
select * from t1 where col_ts != col_datetime order by id;
select * from t1 where col_ts > col_datetime order by id;
select * from t1 where col_ts >= col_datetime order by id;
select * from t1 where col_ts <=> col_datetime order by id;

-- LogicalAnd(and) LogicalOr(or),LogicalNot(not)
select * from t1 where not(col_int<0 or col_int>100000) order by id;
select * from t1 where col_int<0 or col_int>100000 order by id;
select * from t1 where not col_int order by id;
