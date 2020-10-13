-- +goose Up
-- +goose StatementBegin
CREATE TABLE reported_posts (
    post_uuid uuid REFERENCES post_references (post_uuid) 
        ON DELETE CASCADE,
    reported_uuid uuid NOT NULL,
    reported_reason text NOT NULL,
    PRIMARY KEY (reported_uuid)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reported_posts;
-- +goose StatementEnd