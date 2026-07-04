CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE user_role AS ENUM (
    'citizen',
    'officer',
    'department_admin',
    'super_admin'
);

CREATE TYPE staff_position AS ENUM (
'field_officer',
'department_admin',
'supervisor'
);

CREATE TYPE report_status AS ENUM (
    'pending',
    'verified',
    'assigned',
    'in_progress',
    'resolved',
    'rejected'
);

CREATE TYPE report_priority AS ENUM (
    'low',
    'medium',
    'high',
    'critical'
);

CREATE TYPE attachment_type AS ENUM (
    'evidence',
    'resolution'
);

CREATE TYPE news_category AS ENUM (
    'announcement',
    'news',
    'event',
    'emergency'
);

CREATE TYPE news_scope AS ENUM (
    'global',
    'district'
);

CREATE TABLE districts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE departments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    home_district_id UUID REFERENCES districts(id),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number TEXT CHECK (phone_number ~ '^\+?[0-9]{9,15}$'),
    role user_role NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE staff_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department_id UUID NOT NULL REFERENCES departments(id),
    employee_number VARCHAR(50) UNIQUE NOT NULL,
    position staff_position,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    description TEXT,

    CONSTRAINT uq_category_department UNIQUE(department_id,name),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    incident_district_id UUID NOT NULL REFERENCES districts(id),
    assigned_staff_id UUID REFERENCES staff_profiles(id),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    report_number VARCHAR(50) UNIQUE,

    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    address_landmark TEXT,

    status report_status DEFAULT 'pending',
    priority report_priority DEFAULT 'medium',

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE report_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    type attachment_type NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE report_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    status report_status NOT NULL,
    notes TEXT,
    is_internal BOOLEAN DEFAULT FALSE,
    actor_id UUID REFERENCES users(id),

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE regional_news (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID REFERENCES departments(id),
    district_id UUID REFERENCES districts(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category news_category,
    banner_url TEXT,
    target_scope news_scope DEFAULT 'global',
    is_pinned BOOLEAN DEFAULT FALSE,
    created_by UUID REFERENCES users(id),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE emergency_contacts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    department_id UUID REFERENCES departments(id),
    district_id UUID REFERENCES districts(id),
    name VARCHAR(100) NOT NULL,
    phone_number TEXT CHECK (phone_number ~ '^\+?[0-9]{3,20}$'),
    description TEXT,
    icon_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE weather_cache (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    district_id UUID UNIQUE REFERENCES districts(id),
    temperature NUMERIC(5,2),
    humidity INTEGER,
    weather VARCHAR(100),
    air_quality INTEGER,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_reports_status
ON reports(status);

CREATE INDEX idx_reports_priority
ON reports(priority);

CREATE INDEX idx_reports_category
ON reports(category_id);

CREATE INDEX idx_reports_staff
ON reports(assigned_staff_id);

CREATE INDEX idx_reports_district
ON reports(incident_district_id);

CREATE INDEX idx_news_district
ON regional_news(district_id);

CREATE INDEX idx_staff_department
ON staff_profiles(department_id);

CREATE INDEX idx_weather_district
ON weather_cache(district_id);