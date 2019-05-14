use test;

select * from student;

select * from student where id =1 or id =2;

select * from student where id < 100;
select * from student where id = 100 or (id > 101 and id <200) or id =399;
select * from student where mail="aaa@pingcap.com" or mail="bbb@pingcap.com";
select mail from student where mail>="aaa@pingcap.com" and mail<"bbb@pingcap.com";
select mail from student where mail="a@pingcap.com" or (mail>="aaa@pingcap.com" and mail<"bbb@pingcap.com") or mail="d@pingcap.com";
select * from student limit 10;
select avg(score) from student;
select count(score) from student group by score;
select * from student where id = 1 and last_name = 2;
