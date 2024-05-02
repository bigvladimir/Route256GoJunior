-- +goose Up
-- +goose StatementBegin
create table pvz(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    adress TEXT NOT NULL DEFAULT '',
    contacts TEXT NOT NULL DEFAULT ''
);
CREATE TABLE orders (
    order_id BIGSERIAL PRIMARY KEY NOT NULL,
    pvz_id INT NOT NULL,
    customer_id INT NOT NULL,
    storage_last_time TIMESTAMP NOT NULL,
    is_completed BOOLEAN NOT NULL,
    complete_time TIMESTAMP NOT NULL,
    is_refunded BOOLEAN NOT NULL,
    arrival_time TIMESTAMP NOT NULL,
    package_type TEXT,
    weight FLOAT NOT NULL,
    price INT NOT NULL
);
ALTER TABLE orders ADD CONSTRAINT fk_pvz_id FOREIGN KEY (pvz_id) REFERENCES pvz(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table orders;
drop table pvz;
-- +goose StatementEnd
