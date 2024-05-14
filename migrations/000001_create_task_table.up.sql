CREATE TYPE progress_status AS ENUM(
    'canceled',
    'in progress',
    'done'
);

CREATE TABLE task (
    id UUID UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    modified_at TIMESTAMP NOT NULL DEFAULT NOW(),
    header TEXT NOT NULL,
    text TEXT NOT NULL,
    deadline TIMESTAMP NOT NULL,
    progress_status progress_status NOT NULL,
    is_urgent BOOLEAN NOT NULL,
    is_important BOOLEAN NOT NULL,
    owner_id INTEGER NOT NULL,
    parent_id UUID,
    possible_deadline TIMESTAMP NOT NULL,
    weight INTEGER NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (parent_id) REFERENCES task(id)
)
