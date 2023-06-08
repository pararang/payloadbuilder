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

INSERT INTO public.funder_actions
(funder_id, action_id)
VALUES(1, 1);

CREATE TABLE public.destination_fields (
	id serial NOT NULL,
	"key" varchar NOT NULL,
	data_type varchar NOT NULL,
	funder_action_id integer NOT NULL,
	CONSTRAINT destination_fields_pk PRIMARY KEY (id),
	CONSTRAINT destination_fields_fk_funder_action FOREIGN KEY (funder_action_id) REFERENCES public.funder_actions(id)
);
CREATE INDEX destination_fields_funder_action_id_idx ON public.destination_fields (funder_action_id);

ALTER TABLE public.destination_fields ADD preprocessor varchar NULL;

CREATE TABLE public.payload_maps (
	id serial NOT NULL,
	funder_action_id int NOT NULL,
	destination_field_id int NOT NULL,
	source_field_id int NULL,
	CONSTRAINT payload_maps_pk PRIMARY KEY (id),
	CONSTRAINT payload_maps_fk_funder_action FOREIGN KEY (funder_action_id) REFERENCES public.funder_actions(id),
	CONSTRAINT payload_maps_fk_destination_field FOREIGN KEY (destination_field_id) REFERENCES public.destination_fields(id)
);
CREATE INDEX payload_maps_funder_action_id_idx ON public.payload_maps (funder_action_id);


