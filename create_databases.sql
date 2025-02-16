-- Create databases
CREATE DATABASE flowstate_demo;
CREATE DATABASE flowstate_test;

-- Connect to flowstate_demo (you'll need to run the same for flowstate_test)
\c flowstate_demo

-- Create schema
CREATE SCHEMA IF NOT EXISTS main;

-- Create users table
CREATE TABLE main.users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

-- Create flows table
CREATE TABLE main.flows (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(255) NOT NULL REFERENCES main.users(username),
    content JSONB
);

-- Create flows_access table
CREATE TABLE main.flows_access (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    user_id INTEGER NOT NULL REFERENCES main.users(id),
    flow_id INTEGER NOT NULL REFERENCES main.flows(id),
    access VARCHAR(50) NOT NULL,
    UNIQUE(user_id, flow_id)
);

-- Create indexes
CREATE INDEX idx_users_username ON main.users(username);
CREATE INDEX idx_flows_owner ON main.flows(owner);
CREATE INDEX idx_flows_access_user_id ON main.flows_access(user_id);
CREATE INDEX idx_flows_access_flow_id ON main.flows_access(flow_id);

-- Add soft delete indexes
CREATE INDEX idx_users_deleted_at ON main.users(deleted_at);
CREATE INDEX idx_flows_deleted_at ON main.flows(deleted_at);
CREATE INDEX idx_flows_access_deleted_at ON main.flows_access(deleted_at);