CREATE TABLE users (
    id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE photos(
    id BIGSERIAL PRIMARY KEY NOT NULL,
    image_url TEXT
);

CREATE TABLE user_details(
    user_id VARCHAR(50) UNIQUE NOT NULL,
    photo_id INT,
    email VARCHAR(255),
    role VARCHAR(100),
    about TEXT,
    FOREIGN KEY (photo_id) REFERENCES photos(id)
);

CREATE TABLE projects(
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    key VARCHAR(4) NOT NULL,
    description TEXT,
    created_by_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (created_by_id) REFERENCES users(id)
);

CREATE TYPE invitation_status AS ENUM('pending', 'accepted', 'rejected');

CREATE TABLE project_invitations(
    id BIGSERIAL PRIMARY KEY,
    receiver_id VARCHAR(50) NOT NULL,
    project_id VARCHAR(50) NOT NULL,
    sender_id VARCHAR(50) NOT NULL,
    status invitation_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (receiver_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (sender_id) REFERENCES users(id),
    CONSTRAINT check_self_invite_not_allowed CHECK (receiver_id != sender_id)
);

CREATE TABLE components(
    id BIGSERIAL PRIMARY KEY,
    project_id VARCHAR(50) NOT NULL,
    leader_id VARCHAR(50),
    name VARCHAR(75) NOT NULL,
    description TEXT,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (leader_id) REFERENCES users(id)
);

CREATE TABLE user_project_xref(
    user_id VARCHAR(50) NOT NULL,
    project_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    CONSTRAINT duplicate_entry_not_allowed UNIQUE (user_id, project_id)
);

CREATE TYPE ticket_priority AS ENUM ('0', '1', '2', '3', '4');

CREATE TYPE ticket_status AS ENUM (
    'open',
    'under development',
    'under review',
    'deployed to dev',
    'tested',
    'closed'
);

CREATE TYPE ticket_type AS ENUM (
    'bug',
    'devOps',
    'epic',
    'story',
    'sub-task',
    'technical-story',
    'design'
);

CREATE TABLE tickets (
    id VARCHAR(255) PRIMARY KEY,
    number INT NOT NULL,
    type ticket_type NOT NULL,
    priority ticket_priority NOT NULL DEFAULT '2',
    title VARCHAR(255) NOT NULL,
    story_points INT DEFAULT 0,
    description TEXT,
    status ticket_status NOT NULL DEFAULT 'open',
    created_by_id VARCHAR(50) NOT NULL,
    assignee_id VARCHAR(50) NULL,
    project_id VARCHAR(50) NOT NULL,
    component_id INT NULL,
    updated_by_id VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    FOREIGN KEY (assignee_id) REFERENCES users(id),
    FOREIGN KEY (updated_by_id) REFERENCES users(id),
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (component_id) REFERENCES components(id),
    CONSTRAINT duplicate_ticket_number_not_allowed UNIQUE (number, project_id)
);

CREATE
OR REPLACE FUNCTION f_generate_ticket_number () RETURNS trigger AS '
DECLARE BEGIN NEW.number = (
    SELECT
        COUNT(*)
    from
        tickets
    WHERE
        tickets.project_id = NEW.project_id
) + 1;

RETURN NEW;

END;

' LANGUAGE 'plpgsql';

CREATE TRIGGER trig_generate_ticket_number BEFORE
INSERT
    ON tickets FOR EACH ROW EXECUTE PROCEDURE f_generate_ticket_number();

CREATE view tickets_view AS
SELECT
    (p.key || '-' || t.number) as key,
    t.id,
    t.number,
    t.type,
    t.priority,
    t.title,
    t.story_points,
    t.description,
    t.status,
    t.component_id,
    t.created_at,
    t.updated_at,
    p.id as project_id,
    p.name as project_name,
    creator.id as creator_id,
    creator.username as creator_username,
    assignee.id as assignee_id,
    assignee.username as assignee_username
FROM
    tickets as t
    INNER JOIN projects as p ON p.id = t.project_id
    INNER JOIN users as creator ON creator.id = t.created_by_id
    LEFT JOIN users as assignee ON assignee.id = t.assignee_id;

CREATE TABLE comments(
    id BIGSERIAL PRIMARY KEY,
    ticket_id VARCHAR(50) NOT NULL,
    author_id VARCHAR(50) NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (ticket_id) REFERENCES tickets(id),
    FOREIGN KEY (author_id) REFERENCES users(id)
)