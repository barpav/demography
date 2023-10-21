CREATE TABLE people (
    id bigserial PRIMARY KEY,
    surname varchar(150) NOT NULL,
    person_name varchar(150) NOT NULL,
    patronymic varchar(150),
    age smallint,
    gender gender,
    country varchar(5)
);