--
-- PostgreSQL database dump
--

\restrict mqAAhSOZFUb680QqMfSh4qWogMqOaSiRHoUi9ulfdlRTciQktkSEedXWNasZODR

-- Dumped from database version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.11 (Ubuntu 16.11-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: comments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.comments (
    id integer NOT NULL,
    user_id integer,
    recipe_id integer,
    comment text,
    rating double precision
);


ALTER TABLE public.comments OWNER TO postgres;

--
-- Name: comments_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.comments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.comments_id_seq OWNER TO postgres;

--
-- Name: comments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.comments_id_seq OWNED BY public.comments.id;


--
-- Name: meal_type; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.meal_type (
    id integer NOT NULL,
    name character varying(50)
);


ALTER TABLE public.meal_type OWNER TO postgres;

--
-- Name: meal_type_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.meal_type_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.meal_type_id_seq OWNER TO postgres;

--
-- Name: meal_type_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.meal_type_id_seq OWNED BY public.meal_type.id;


--
-- Name: recipe_likes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.recipe_likes (
    id integer NOT NULL,
    user_id integer NOT NULL,
    recipe_id integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.recipe_likes OWNER TO postgres;

--
-- Name: recipe_likes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.recipe_likes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recipe_likes_id_seq OWNER TO postgres;

--
-- Name: recipe_likes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.recipe_likes_id_seq OWNED BY public.recipe_likes.id;


--
-- Name: recipes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.recipes (
    id integer NOT NULL,
    user_id integer,
    name_recipe text,
    description text,
    meal_type_id integer,
    img_url character varying,
    guest_name character varying(100),
    steps character varying,
    is_active boolean DEFAULT false
);


ALTER TABLE public.recipes OWNER TO postgres;

--
-- Name: recipes_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.recipes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.recipes_id_seq OWNER TO postgres;

--
-- Name: recipes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.recipes_id_seq OWNED BY public.recipes.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id integer NOT NULL,
    name character varying(100),
    email text,
    password character varying,
    url_photo character varying,
    google_id integer
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: comments id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments ALTER COLUMN id SET DEFAULT nextval('public.comments_id_seq'::regclass);


--
-- Name: meal_type id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.meal_type ALTER COLUMN id SET DEFAULT nextval('public.meal_type_id_seq'::regclass);


--
-- Name: recipe_likes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipe_likes ALTER COLUMN id SET DEFAULT nextval('public.recipe_likes_id_seq'::regclass);


--
-- Name: recipes id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipes ALTER COLUMN id SET DEFAULT nextval('public.recipes_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: comments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.comments (id, user_id, recipe_id, comment, rating) FROM stdin;
1	2	1	Best pasta I have ever had	5
2	2	1	Needs more cheese!!!!	3.5
3	1	2	Good lunch	4
4	1	2	Too much tomato sauce	3
5	6	3	La mejor comida de El Salvador	5
6	6	3	Las mejores estan en Santa Ana	5
7	6	3	test #2	4
8	6	3	test #3	3
9	6	3	test #4	2
10	6	1	test #1	3
11	6	1	test #2	3.5
12	6	1	test #3	1.5
13	6	4	The best sauce indeed	5
14	6	2	Good food, kinda spicy	4
15	6	2	test #4	3.5
16	3	1	The best thing among Italy's pastas	5
17	6	6	Best test!!	5
18	3	6	Last test of the day #1	4
\.


--
-- Data for Name: meal_type; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.meal_type (id, name) FROM stdin;
1	breakfast
2	lunch
3	dinner
4	snack
\.


--
-- Data for Name: recipe_likes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.recipe_likes (id, user_id, recipe_id, created_at) FROM stdin;
\.


--
-- Data for Name: recipes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.recipes (id, user_id, name_recipe, description, meal_type_id, img_url, guest_name, steps, is_active) FROM stdin;
1	1	Spaghetti Carbonara	Classic roman pasta with pancetta, egg and cheese	2	\N	\N	Take a carbonara sauce	t
3	\N	Pupusas de queso	Ricas comida tipica salvadoreña	1	\N	Juana Pérez	Pon queso entre la masa, y cocinas	t
4	\N	Alfredo pasta	One of the best recipes of all time	2	\N	Fernando Gómez	put a load of cheese into the pasta	t
2	2	Chiken Tikka Masala	Creamy tomato-based indian curry with grilled chicken	2	\N	\N	Take a Tikka masala sauce	t
6	\N	Test #1	This is a small description	2	https://res.cloudinary.com/dii7ichy9/image/upload/v1751491540/lgtmedhuyzi8o7rmqpte.png	Kira L.	put your desire input here	t
5	6	Edit test sucess!	This is a small description for the edit endpoint	2	https://res.cloudinary.com/dii7ichy9/image/upload/v1751491221/snrhrosrh1ivwntuiexd.png	\N	put your desire input here	t
7	\N	Testing axios	asas	1		Fernando Gómez	adasd	t
8	\N	Testing post once	adsdaasd	2		Fernando Gómez	asdasdad	t
9	\N	Pasta Alfredo	The best pasta you'll ever have	3		John Doe	1-Use your favorite pasta\r\n2-Boil the pasta\r\n3-Usa parmigiano cheese\r\n4-Enjoy	t
10	\N	Test new like features	test	1		Jane Doe	there's no steps	t
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password, url_photo, google_id) FROM stdin;
1	Jane Doe	janedoe@gmail.com	123456789	\N	\N
2	John Doe	johndoe@gmail.com	123456789	\N	\N
3	Ana Castillo	anacastillo@gmail.com	$2a$10$Y8gw3Gms6pKABfGr.56R9.utHvCfM35Bt6VvF3cE9jB2B6I/ZQZ3.	\N	\N
4	Tizio	acaso@gmail.com	$2a$10$rV8Zzio5KTt8R2nUHz.geO/qCYc32Z9.aQuBDXN/WAiRGdlkVGB4a	\N	\N
5	Allison Juarez	aj@gmail.com	$2a$10$XvZXF7Of0pVPM7AIqRldPefj..E0yAEVxEZhaNg7TfxX8AifU28Ji	https://res.cloudinary.com/dii7ichy9/image/upload/v1750720511/exzotygyafplsflhrufs.png	\N
6	Tom Martin	tm@gmail.com	$2a$10$YsO/IX6gDDPlLi3C//rNC.v445P/I0z6fCjmn9F5HYuRmwtwuBPj.	https://res.cloudinary.com/dii7ichy9/image/upload/v1750972491/roykcsnaaadh0hmrn8l3.png	\N
\.


--
-- Name: comments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.comments_id_seq', 18, true);


--
-- Name: meal_type_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.meal_type_id_seq', 4, true);


--
-- Name: recipe_likes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.recipe_likes_id_seq', 1, false);


--
-- Name: recipes_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.recipes_id_seq', 10, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 6, true);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: meal_type meal_type_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.meal_type
    ADD CONSTRAINT meal_type_pkey PRIMARY KEY (id);


--
-- Name: recipe_likes recipe_likes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipe_likes
    ADD CONSTRAINT recipe_likes_pkey PRIMARY KEY (id);


--
-- Name: recipe_likes recipe_likes_user_id_recipe_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipe_likes
    ADD CONSTRAINT recipe_likes_user_id_recipe_id_key UNIQUE (user_id, recipe_id);


--
-- Name: recipes recipes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_recipe_likes_recipe; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_recipe_likes_recipe ON public.recipe_likes USING btree (recipe_id);


--
-- Name: idx_recipe_likes_user; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_recipe_likes_user ON public.recipe_likes USING btree (user_id);


--
-- Name: comments comments_recipe_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_recipe_id_fkey FOREIGN KEY (recipe_id) REFERENCES public.recipes(id);


--
-- Name: comments comments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: recipe_likes recipe_likes_recipe_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipe_likes
    ADD CONSTRAINT recipe_likes_recipe_id_fkey FOREIGN KEY (recipe_id) REFERENCES public.recipes(id) ON DELETE CASCADE;


--
-- Name: recipe_likes recipe_likes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipe_likes
    ADD CONSTRAINT recipe_likes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: recipes recipes_meal_type_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_meal_type_id_fkey FOREIGN KEY (meal_type_id) REFERENCES public.meal_type(id);


--
-- Name: recipes recipes_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.recipes
    ADD CONSTRAINT recipes_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

\unrestrict mqAAhSOZFUb680QqMfSh4qWogMqOaSiRHoUi9ulfdlRTciQktkSEedXWNasZODR

