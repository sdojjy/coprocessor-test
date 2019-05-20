use coprocessor;
SET time_zone='+00:00';

select * from t1  order by id;
select col_ts from t1 order by col_ts;
select id from t1 order by id;
select * from t1 where id<-100 or id>100 order by id;
select * from t2 where id<10 or id>100 order by id;
select * from t1 where id in(-1,2,-3,10,1000) order by id;
select * from t2 where id in(1,2,3,10,1000) order by id;
select * from t1 where id in(-1,2,-3,10,100,0,999) or (id<50 and id>20) or id <-5 order by id;





