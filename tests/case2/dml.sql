use test;

drop table if exists `student`;

CREATE TABLE if not exists `student` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `last_name` varchar(255) NOT NULL,
  `first_name` varchar(255) DEFAULT NULL,
  `age` int(11) DEFAULT NULL,
  `score` int(11) DEFAULT NULL,
  `mail` varchar(255) DEFAULT NULL,
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `name` (`last_name`,`first_name`),
  KEY `score` (`score`),
  KEY `age` (`age`),
  UNIQUE KEY `mail` (`mail`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

insert into `student`(id,last_name, first_name, age, score, mail, create_time) values (1, "2", "aaa", 12, 99, "aaa@pingcap.com", "2016-11-21 05:05:05");
insert into `student`(id,last_name, first_name, age, score, mail, create_time) values (2, "bbb", "bbb", 12, 65535, "bbb@pingcap.com", "2016-11-21 05:05:05");
insert into `student`(id,last_name, first_name, age, score, mail, create_time) values (104, "ccc", "ccc", 12, 65534, "ccc@pingcap.com", "1970-01-01 05:05:05");
insert into `student`(id,last_name, first_name, age, score, mail, create_time) values (100, "ddd", "ddd", -1, -1, "d@pingcap.com", "1969-11-21 05:05:05");
insert into `student`(id,last_name, first_name, age, score, mail, create_time) values (399, "eee", "eee", -2, 99, "eee@pingcap.com", "3098-11-21 05:05:05");