use test;

select * from t1;
select col_ts from t1;
select * from t1 where id<-100 or id>100;
select * from t2 where id<10 or id>100;
select * from t1 where id in(-1,2,-3,10,100);
select * from t2 where id in(1,2,3,10,100);
select * from t1 where id in(-1,2,-3,10,100) or (id<50 and id>20) or id <-5;
