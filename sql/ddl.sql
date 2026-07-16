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
    id BIGSERIAL PRIMARY KEY,
    bmkg_adm4_code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    home_district_id BIGINT REFERENCES districts(id),
    nik CHAR(16) UNIQUE NOT NULL CHECK (nik ~ '^[0-9]{16}$'),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone_number TEXT CHECK (phone_number ~ '^\+?[0-9]{9,15}$'),
    role user_role NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE staff_profiles (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department_id BIGINT NOT NULL REFERENCES departments(id),
    employee_number VARCHAR(50) UNIQUE NOT NULL,
    position staff_position,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,

    CONSTRAINT uq_category_department UNIQUE(department_id,name),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE reports (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    category_id BIGINT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    incident_district_id BIGINT NOT NULL REFERENCES districts(id),
    assigned_staff_id BIGINT REFERENCES staff_profiles(id),
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
    id BIGSERIAL PRIMARY KEY,
    report_id BIGINT NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,
    type attachment_type NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE report_histories (
    id BIGSERIAL PRIMARY KEY,
    report_id BIGINT NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    status report_status NOT NULL,
    notes TEXT,
    is_internal BOOLEAN DEFAULT FALSE,
    actor_id BIGINT REFERENCES users(id),

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE regional_news (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT REFERENCES departments(id),
    district_id BIGINT REFERENCES districts(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category news_category,
    banner_url TEXT,
    target_scope news_scope DEFAULT 'global',
    is_pinned BOOLEAN DEFAULT FALSE,
    created_by BIGINT REFERENCES users(id),

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE emergency_contacts (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT REFERENCES departments(id),
    district_id BIGINT REFERENCES districts(id),
    name VARCHAR(100) NOT NULL,
    phone_number TEXT CHECK (phone_number ~ '^\+?[0-9]{3,20}$'),
    description TEXT,
    icon_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE weather_forecasts (
    id BIGSERIAL PRIMARY KEY,
    district_id BIGINT NOT NULL REFERENCES districts(id),
    forecast_time TIMESTAMP NOT NULL,
    temperature DECIMAL(5,2) NOT NULL,
    humidity INTEGER NOT NULL,
    weather VARCHAR(100) NOT NULL,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),

    UNIQUE (
        district_id,
        forecast_time
    )
);

CREATE TABLE water_statuses (
    id BIGSERIAL PRIMARY KEY,
    district_id BIGINT NOT NULL REFERENCES districts(id),
    status VARCHAR(30) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    estimated_duration INTEGER NOT NULL,
    estimated_recovery_at TIMESTAMP NOT NULL,
    reason TEXT,
    created_by BIGINT NOT NULL REFERENCES users(id),

    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_water_status_district
ON water_statuses(district_id);

CREATE INDEX idx_water_status_created
ON water_statuses(created_at DESC);

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