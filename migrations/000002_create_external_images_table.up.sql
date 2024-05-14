CREATE TABLE external_image (
    id SERIAL UNIQUE NOT NULL,
    url TEXT NOT NULL,
    task_id UUID NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (task_id) REFERENCES task (id)
);