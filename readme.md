# 中国行政区划爬虫

## 数据源

[中华人民共和国国家统计局](http://www.stats.gov.cn/sj/tjbz/tjyqhdmhcxhfdm/2022/index.html)


## 数据文件

最新数据量 664486 （2022年10月31日）

本仓库提供已爬取的最新数据的 SQL 文件，文件名为`2023.sql`

### SQL 格式数据

```sql
-- ----------------------------
-- Table structure for regions
-- ----------------------------
DROP TABLE IF EXISTS `regions`;
CREATE TABLE `regions` (
  `code` bigint unsigned NOT NULL DEFAULT '0',
  `name` varchar(255) NOT NULL DEFAULT '',
  `level` tinyint unsigned NOT NULL DEFAULT '0',
  `p_code` bigint unsigned NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

```

> 预览部分数据

```bash
bunzip2 2023.sql.tar.bz2
tar -zxvf 2023.sql.tar
less 2023.sql
```

## Run

你也可以自己执行该脚本进行爬取

```bash
go mod tidy

go run main.go
```