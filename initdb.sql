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

CREATE TABLE public.funder_actions (
	id serial NOT NULL,
	funder_id integer NOT NULL,
	action_id integer NOT NULL,
	CONSTRAINT funder_actions_pk PRIMARY KEY (id),
	CONSTRAINT funder_actions_fk_funder FOREIGN KEY (funder_id) REFERENCES public.funders(id),
	CONSTRAINT funder_actions_fk_action FOREIGN KEY (action_id) REFERENCES public.actions(id)
);


