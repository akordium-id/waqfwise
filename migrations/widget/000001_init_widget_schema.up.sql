-- migrations/widget/000001_init_widget_schema.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS mosques (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    domain VARCHAR(255),
    primary_color VARCHAR(7) NOT NULL DEFAULT '#1B5E20',
    logo_url TEXT,
    description TEXT,
    address TEXT,
    phone VARCHAR(50),
    email VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_mosques_slug ON mosques(slug);
CREATE INDEX idx_mosques_active ON mosques(is_active);

CREATE TABLE IF NOT EXISTS mosque_campaign_types (
    mosque_id UUID NOT NULL REFERENCES mosques(id) ON DELETE CASCADE,
    campaign_type VARCHAR(50) NOT NULL,
    PRIMARY KEY (mosque_id, campaign_type)
);
