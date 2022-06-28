create table users(
                      id int GENERATED ALWAYS AS IDENTITY,
                      first_name varchar(500),
                      last_name varchar(500),
                      email varchar(500),
                      password_hash varchar(1000),
                      update_date timestamp
);

insert into users (first_name, last_name, email, password_hash) values ('Dmitry', 'Bogadukhov', 'test1@email.com', 'passHash1'), ('Petr', 'Smirnov', 'test1@email.com', 'passHash2')