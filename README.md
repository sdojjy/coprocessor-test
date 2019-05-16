# coprocessor test

## 生成随机测试数据

```bash
# --gen-data 为 true 时会根据传入的表结构生成随机数据， 否则会执行对比测试
# --ddl=ddl.sql 表结构文档，可以参考 ddl.sql
# --records=1000 每个表生成的数据条数， 默认1000
# 结果输出到标准输出，可根据需要重定向到文件
# example:
 go run main.go --gen-data=true --ddl=ddl.sql --records=1000 > rand-data.sql
```

测试文件组织结构

```bash
tests  # case 根目录
├── case1  # case1 目录
│   ├── ddl.sql  # 数据文件，执行case 1 前会执行dml.sql进行数据清理和插入，名字必须是dml.sql
│   ├── query.sql # sql 查询语句，程序会对比每一条语句的结果 
│   └── query2.sql
└── case2 # case 2 目录
    ├── ddl.sql
    └── query.sql

```


## 执行对比测试

```bash
# dsn1 和 dsn2 分别表示要对比执行的两个数据库，可以是 tidb 和 mysql
# dsn1 默认值为 root:@tcp(127.0.0.1:4000)/test?charset=utf8
# dsn1 默认值为 root:@tcp(127.0.0.1:3306)/test?charset=utf8
# test-dir 指定test case 的 base 目录， 默认值为当前目录的./tests/
# strict 是否做严格的比较(结果顺序一致), 默认为 true， 否则尝试不比较顺序
 go run main.go --dsn1=<tidb-dsn> --dsn2==<mysql-dsn> --test-dir=<case-dir> 
```

## TODO

- [ ] 随机数据生成优化
- [ ] 更多边界值支持
- [ ] 执行结果对比完善
