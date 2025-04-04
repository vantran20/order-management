CREATE TABLE IF NOT EXISTS users
(
    id            BIGINT PRIMARY KEY,
    name          TEXT                     NOT NULL CONSTRAINT users_name_check CHECK (name <> ''::TEXT),
    email         TEXT                     NOT NULL CONSTRAINT users_email_check CHECK (email <> ''::TEXT),
    password      TEXT                     NOT NULL CONSTRAINT users_password_check CHECK (password <> ''::TEXT),
    status        TEXT                     NOT NULL CHECK (status <> ''::text),
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    UNIQUE (email),
    UNIQUE (password)
    );
CREATE INDEX IF NOT EXISTS users_email_index ON users(email);

CREATE TABLE IF NOT EXISTS public.products
(
    id          BIGINT PRIMARY KEY,
    name        TEXT                     NOT NULL CONSTRAINT products_name_check CHECK (name <> ''::text),
    description TEXT                     NOT NULL CHECK (description <> ''::text),
    status      TEXT                     NOT NULL CHECK (status <> ''::text),
    price       FLOAT                    NOT NULL CHECK (price > 0::FLOAT),
    stock       BIGINT                   NOT NULL CHECK (stock >= 0),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.orders
(
    id          BIGINT PRIMARY KEY,
    user_id     BIGINT                   NOT NULL REFERENCES public.users (id),
    status      TEXT                     NOT NULL CHECK (status <> ''::text),
    total_cost  FLOAT                    NOT NULL CHECK (total_cost >= 0::FLOAT),
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public.order_items
(
    id         BIGINT PRIMARY KEY,
    order_id   BIGINT                   NOT NULL REFERENCES public.orders (id),
    product_id BIGINT                   NOT NULL REFERENCES public.products (id),
    quantity   BIGINT                   NOT NULL CHECK (quantity >= 0),
    price      FLOAT                    NOT NULL CHECK (price >= 0::FLOAT),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
 );
CREATE UNIQUE INDEX IF NOT EXISTS order_item_uidx_order_id_product_id ON public.order_items (order_id, product_id);
