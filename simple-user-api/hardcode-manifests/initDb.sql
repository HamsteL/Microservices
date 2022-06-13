create table users(
    id int GENERATED ALWAYS AS IDENTITY,
    first_name varchar(500),
    last_name varchar(500),
    update_date timestamp
);

insert into users (first_name, last_name) values ('Dmitry', 'Bogadukhov'), ('Petr', 'Smirnov')