CREATE TABLE IF NOT EXISTS advert (
	id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name varchar(255),
	description varchar(255),
	photos text[],
	price integer
);


