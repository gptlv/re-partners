-- +goose Up
-- +goose StatementBegin
CREATE TABLE pack_sizes (
    id BIGSERIAL PRIMARY KEY,
    size BIGINT NOT NULL UNIQUE
);

INSERT INTO pack_sizes (size) VALUES (250),(500),(1000),(2000),(5000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pack_sizes;
-- +goose StatementEnd
