CREATE TABLE users
(
    id serial primary key,
    name varchar(255) not null,
    username varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE todo_lists
(
    id serial primary key,
    title varchar(255) not null,
    description varchar(255),
    user_id int not null,
    foreign key (user_id) references users(id) on delete cascade
);

CREATE TABLE todo_items
(   
    id serial primary key,
    title varchar(255) not null,
    description varchar(255),
    done boolean not null default false,
    list_id int not null,
    foreign key (list_id) references todo_lists(id) on delete cascade
);
