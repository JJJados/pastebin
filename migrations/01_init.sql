-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_references (
    post_uuid uuid NOT NULL,
    read_access_uuid uuid NOT NULL,
    admin_access_uuid uuid NOT NULL,
    PRIMARY KEY (post_uuid)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE posts (
    post_uuid uuid REFERENCES post_references (post_uuid) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    title text NOT NULL,
    content text NOT NULL,
    public_access BOOLEAN NOT NULL DEFAULT TRUE,
    reported BOOLEAN NOT NULL DEFAULT FALSE,
    created TIMESTAMP NOT NULL DEFAULT now(),
    updated TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE posts;
DROP TABLE post_references;
-- +goose StatementEnd