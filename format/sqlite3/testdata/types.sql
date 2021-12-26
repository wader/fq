CREATE TABLE types (
    cint INT PRIMARY KEY,
    cinteger INTEGER,
    cvarchar VARCHAR(30),
    ctext TEXT,
    creal REAL,
    cblob BLOB,
    ctimestamp TIMESTAMP
);

INSERT INTO types VALUES(0, 0, 'var1', 'text1', 0.5, "blob1", '2022-03-10 10:00:00');
INSERT INTO types VALUES(1, 1, 'var2', 'test2', 1.5, "blob2", '2022-03-10 10:00:01');
INSERT INTO types VALUES(128, 128, 'var3', 'test3', 128.5, "blob3", '2022-03-10 10:00:02');
INSERT INTO types VALUES(-128, 128, 'var3', 'test3', -128.5, "blob3", '2022-03-10 10:00:03');
INSERT INTO types VALUES(9223372036854775807, 9223372036854775807, 'var4', 'test4', 9223372036854775807.5, "blob4", '2022-03-10 10:00:04');
INSERT INTO types VALUES(-9223372036854775808, -9223372036854775808, 'var5', 'test5', -9223372036854775808.5, "blob5", '2022-03-10 10:00:05');
