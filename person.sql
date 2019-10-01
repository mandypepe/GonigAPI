-- auto-generated definition
create table person
(
  Id         int auto_increment
    primary key,
  first_name varchar(200) not null,
  last_name  varchar(200) not null
);


INSERT INTO plazavea.person (Id, first_name, last_name) VALUES (1, 'emilio', 'soto');
INSERT INTO plazavea.person (Id, first_name, last_name) VALUES (2, 'alfonso', 'soria');
INSERT INTO plazavea.person (Id, first_name, last_name) VALUES (3, 'andres', 'paredes');