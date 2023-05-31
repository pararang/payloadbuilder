CREATE TABLE public.funders (
	id serial NOT NULL,
	name varchar NOT NULL,
	code varchar NOT NULL,
	CONSTRAINT funders_pk PRIMARY KEY (id)
);


INSERT INTO public.funders
("name", code)
VALUES('Alami', 'alami');


CREATE TABLE public.actions (
	id serial NOT NULL,
	code varchar NOT NULL,
	description varchar NULL,
	CONSTRAINT actions_pk PRIMARY KEY (id)
);


INSERT INTO public.actions
(code, description)
VALUES('approval', 'register farmer to funder');

ALTER TABLE public.actions ADD CONSTRAINT actions_unique_code UNIQUE (code);
ALTER TABLE public.funders ADD CONSTRAINT funders_unique_code UNIQUE (code);
