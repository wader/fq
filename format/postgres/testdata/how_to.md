## How to generate PostgreSQL test data

### 1) install postgres
Use this links to add repository and install PostgreSQL:
- https://www.postgresql.org/download/linux/
- https://www.postgresql.org/download/linux/redhat/
- https://www.postgresql.org/download/linux/debian/

You may install PostgreSQL with OS repository. But it may contains not all versions of PostgreSQL.

### 2) use default intallation location of postgres.
If you don't want default location. 
You may remove cluster and init it again.

\# Detect PGDATA location:
```shell
ps aux | grep postgre

postgres     887  0.0  0.0 4334456 13608 ?       Ss   07:48   0:00 /usr/pgsql-14/bin/postmaster -D /u02/data
```
Result is:
- `/u02/data` - is PGDATA location
- `/usr/pgsql-14/bin` - location of PostgreSQL bin

### 3) init tables

\# use postgres user
```shell
sudo su postgres
```

\# then init pgbench tables
```shell
/usr/pgsql-14/bin/pgbench -i
```

\# run simple test
```shell
/usr/pgsql-14/bin/pgbench -T 60 -j 100
```

### 4) run commands in psql

\# start psql
```shell
/usr/pgsql-14/bin/psql
```

\# display tables
```psql
\dt+

Schema |       Name       | Type  |  Owner   | Persistence | Access method |  Size  | Description
--------+------------------+-------+----------+-------------+---------------+--------+-------------
public | pgbench_accounts | table | postgres | permanent   | heap          | 13 MB  |
public | pgbench_branches | table | postgres | permanent   | heap          | 40 kB  |
public | pgbench_history  | table | postgres | permanent   | heap          | 968 kB |
public | pgbench_tellers  | table | postgres | permanent   | heap          | 40 kB  |
```

\# run CHECKPOINT to avoid partially written data
```psql
CHECKPOINT;
```

\# detect location of tables
```psql
SELECT pg_relation_filepath('pgbench_history');
pg_relation_filepath
----------------------
base/13746/24599
```

`base/13746/24599` - location of table in PGDATA

### 5) use root account to copy table file

\# login in root
```shell
sudo su
```

\# go to PGDATA
```shell
cd /u02/data
```

\# copy table file
```shell
cp base/13746/24599 /home/user
```

\# change persmissions
```shell
chwon user:user /home/user/24599
chmod 644 /home/user/24599
```

### 6) Copy tables files to local with scp
```shell
scp user@192.168.0.100:~/24599 .
```

### 7) You may want to cut 2 pages (8192 * 2) from table
```shell
head -c 16384 ./245991 > ./24599_2pages
```

### 8) Locate pg_control file
`global/pg_control` inside PGDATA

### 9) Locate btree index 
Get info about index:
```psql
\d pgbench_accounts
              Table "public.pgbench_accounts"
  Column  |     Type      | Collation | Nullable | Default 
----------+---------------+-----------+----------+---------
 aid      | integer       |           | not null | 
 bid      | integer       |           |          | 
 abalance | integer       |           |          | 
 filler   | character(84) |           |          | 
Indexes:
    "pgbench_accounts_pkey" PRIMARY KEY, btree (aid)
```
Then get path of pgbench_accounts_pkey:
```psql
select pg_relation_filepath('pgbench_accounts_pkey');
 pg_relation_filepath 
----------------------
 base/13746/24596
```
`base/13746/24596` - is a path inside PGDATA of btree index pgbench_accounts_pkey.