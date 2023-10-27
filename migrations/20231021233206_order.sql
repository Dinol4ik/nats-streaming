-- +goose Up
-- +goose StatementBegin
CREATE TABLE "order"
(
    id   UUID NOT NULL  PRIMARY KEY ,
    order_uid VARCHAR(19) NOT NULL UNIQUE ,
    order_info  JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "order"
-- +goose StatementEnd
