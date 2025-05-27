-- ENUM type definitions
CREATE TYPE user_type_enum AS ENUM ('customer', 'organizer', 'admin');
CREATE TYPE user_status_enum AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE venue_type_enum AS ENUM ('indoor', 'outdoor', 'virtual', 'hybrid');
CREATE TYPE event_type_enum AS ENUM ('concert', 'sports', 'theater', 'conference', 'festival', 'other');
CREATE TYPE event_status_enum AS ENUM ('draft', 'published', 'cancelled', 'postponed', 'completed');
CREATE TYPE ticket_category_type_enum AS ENUM ('general', 'vip', 'early_bird', 'group', 'season');
CREATE TYPE ticket_status_enum AS ENUM ('available', 'reserved', 'sold', 'cancelled', 'used');
CREATE TYPE order_status_enum AS ENUM ('pending', 'processing', 'confirmed', 'cancelled', 'refunded', 'partially_refunded');
CREATE TYPE payment_type_enum AS ENUM ('credit_card', 'debit_card', 'paypal', 'bank_transfer', 'digital_wallet');
CREATE TYPE payment_status_enum AS ENUM ('pending', 'processing', 'completed', 'failed', 'cancelled', 'refunded', 'partially_refunded');
CREATE TYPE refund_status_enum AS ENUM ('pending', 'processing', 'completed', 'failed');
CREATE TYPE queue_status_enum AS ENUM ('waiting', 'active', 'expired', 'completed', 'cancelled');
CREATE TYPE reservation_status_enum AS ENUM ('active', 'expired', 'completed', 'cancelled');
CREATE TYPE notification_type_enum AS ENUM ('email', 'sms', 'push', 'in_app');
CREATE TYPE notification_status_enum AS ENUM ('pending', 'sent', 'delivered', 'failed', 'bounced');

-- Table definitions
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    user_type user_type_enum DEFAULT 'customer',
    status user_status_enum DEFAULT 'active',
    email_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP
);

CREATE TABLE venues (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    country VARCHAR(100) NOT NULL,
    capacity INT NOT NULL,
    venue_type venue_type_enum NOT NULL,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    organizer_id BIGINT NOT NULL,
    venue_id BIGINT,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_type event_type_enum NOT NULL,
    status event_status_enum DEFAULT 'draft',
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP,
    timezone VARCHAR(50) NOT NULL,
    is_recurring BOOLEAN DEFAULT FALSE,
    max_tickets_per_order INT DEFAULT 10,
    sale_start_date TIMESTAMP,
    sale_end_date TIMESTAMP,
    image_url VARCHAR(500),
    terms_and_conditions TEXT,
    age_restriction INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organizer_id) REFERENCES users(id),
    FOREIGN KEY (venue_id) REFERENCES venues(id)
);

CREATE TABLE ticket_categories (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    quantity_available INT NOT NULL,
    quantity_sold INT DEFAULT 0,
    max_per_order INT DEFAULT 10,
    sale_start_date TIMESTAMP,
    sale_end_date TIMESTAMP,
    is_transferable BOOLEAN DEFAULT TRUE,
    is_refundable BOOLEAN DEFAULT TRUE,
    category_type ticket_category_type_enum DEFAULT 'general',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id) ON DELETE CASCADE
);

CREATE TABLE tickets (
    id BIGSERIAL PRIMARY KEY,
    ticket_category_id BIGINT NOT NULL,
    ticket_number VARCHAR(50) UNIQUE NOT NULL,
    seat_section VARCHAR(50),
    seat_row VARCHAR(10),
    seat_number VARCHAR(10),
    status ticket_status_enum DEFAULT 'available',
    reserved_at TIMESTAMP NULL,
    reserved_expires_at TIMESTAMP NULL,
    qr_code VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_category_id) REFERENCES ticket_categories(id) ON DELETE CASCADE,
    UNIQUE (ticket_category_id, seat_section, seat_row, seat_number)
);

CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    status order_status_enum DEFAULT 'pending',
    total_amount DECIMAL(10, 2) NOT NULL,
    discount_amount DECIMAL(10, 2) DEFAULT 0.00,
    tax_amount DECIMAL(10, 2) DEFAULT 0.00,
    service_fee DECIMAL(10, 2) DEFAULT 0.00,
    final_amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    email_received VARCHAR(255) NOT NULL,
    notes TEXT,
    expires_at TIMESTAMP,
    confirmed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE order_items (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    ticket_id BIGINT NOT NULL,
    unit_price DECIMAL(10, 2) NOT NULL,
    quantity INT DEFAULT 1,
    subtotal DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id),
    UNIQUE (order_id, ticket_id)
);

CREATE TABLE order_status_history (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    previous_status VARCHAR(50),
    new_status VARCHAR(50) NOT NULL,
    reason TEXT,
    changed_by BIGINT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id)
);

CREATE TABLE payment_methods (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    payment_type payment_type_enum NOT NULL,
    provider VARCHAR(50) NOT NULL,
    last_four_digits VARCHAR(4),
    expiry_month INT,
    expiry_year INT,
    cardholder_name VARCHAR(255),
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    external_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE payments (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL,
    payment_method_id BIGINT,
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    status payment_status_enum DEFAULT 'pending',
    payment_intent_id VARCHAR(255),
    transaction_id VARCHAR(255),
    gateway_response TEXT,
    failure_reason VARCHAR(500),
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (payment_method_id) REFERENCES payment_methods(id)
);

CREATE TABLE refunds (
    id BIGSERIAL PRIMARY KEY,
    payment_id BIGINT NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    reason TEXT,
    status refund_status_enum DEFAULT 'pending',
    refund_id VARCHAR(255),
    gateway_response TEXT,
    processed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (payment_id) REFERENCES payments(id)
);

CREATE TABLE ticket_queues (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    position INT NOT NULL,
    status queue_status_enum DEFAULT 'waiting',
    entered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    activated_at TIMESTAMP,
    expires_at TIMESTAMP,
    completed_at TIMESTAMP,
    session_token VARCHAR(255) UNIQUE,
    estimated_wait_time INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (event_id, user_id, status)
);

CREATE TABLE ticket_reservations (
    id BIGSERIAL PRIMARY KEY,
    ticket_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    order_id BIGINT,
    reserved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    status reservation_status_enum DEFAULT 'active',
    reservation_token VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE TABLE queue_settings (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    is_enabled BOOLEAN DEFAULT FALSE,
    max_concurrent_users INT DEFAULT 1000,
    reservation_timeout_minutes INT DEFAULT 10,
    queue_activation_threshold INT DEFAULT 100,
    estimated_service_time_seconds INT DEFAULT 300,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (event_id) REFERENCES events(id),
    UNIQUE (event_id)
);

CREATE TABLE notification_templates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type notification_type_enum NOT NULL,
    subject VARCHAR(255),
    template_content TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (name, type)
);

CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    template_id BIGINT NOT NULL,
    order_id BIGINT,
    event_id BIGINT,
    type notification_type_enum NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    content TEXT NOT NULL,
    status notification_status_enum DEFAULT 'pending',
    sent_at TIMESTAMP,
    delivered_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (template_id) REFERENCES notification_templates(id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (event_id) REFERENCES events(id)
);

-- Indexes
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_user_type ON users(user_type);
CREATE INDEX idx_status_users ON users(status);

CREATE INDEX idx_city_country ON venues(city, country);
CREATE INDEX idx_venue_type ON venues(venue_type);

CREATE INDEX idx_organizer_id ON events(organizer_id);
CREATE INDEX idx_venue_id ON events(venue_id);
CREATE INDEX idx_start_date ON events(start_date);
CREATE INDEX idx_status_events ON events(status);
CREATE INDEX idx_event_type ON events(event_type);

CREATE INDEX idx_event_id_ticket_categories ON ticket_categories(event_id);
CREATE INDEX idx_category_type ON ticket_categories(category_type);

CREATE INDEX idx_ticket_category_id ON tickets(ticket_category_id);
CREATE INDEX idx_status_tickets ON tickets(status);
CREATE INDEX idx_reserved_expires_at ON tickets(reserved_expires_at);

CREATE INDEX idx_user_id_orders ON orders(user_id);
CREATE INDEX idx_status_orders ON orders(status);
CREATE INDEX idx_order_number ON orders(order_number);
CREATE INDEX idx_expires_at ON orders(expires_at);

CREATE INDEX idx_order_id_order_items ON order_items(order_id);
CREATE INDEX idx_ticket_id_order_items ON order_items(ticket_id);

CREATE INDEX idx_order_id_status_history ON order_status_history(order_id);
CREATE INDEX idx_changed_at ON order_status_history(changed_at);

CREATE INDEX idx_user_id_payment_methods ON payment_methods(user_id);
CREATE INDEX idx_payment_type ON payment_methods(payment_type);

CREATE INDEX idx_order_id_payments ON payments(order_id);
CREATE INDEX idx_status_payments ON payments(status);
CREATE INDEX idx_transaction_id ON payments(transaction_id);

CREATE INDEX idx_payment_id_refunds ON refunds(payment_id);
CREATE INDEX idx_status_refunds ON refunds(status);

CREATE INDEX idx_event_id_ticket_queues ON ticket_queues(event_id);
CREATE INDEX idx_user_id_ticket_queues ON ticket_queues(user_id);
CREATE INDEX idx_status_ticket_queues ON ticket_queues(status);
CREATE INDEX idx_position ON ticket_queues(position);

CREATE INDEX idx_ticket_id_reservations ON ticket_reservations(ticket_id);
CREATE INDEX idx_user_id_reservations ON ticket_reservations(user_id);
CREATE INDEX idx_expires_at_reservations ON ticket_reservations(expires_at);
CREATE INDEX idx_status_reservations ON ticket_reservations(status);

CREATE INDEX idx_type_notification_templates ON notification_templates(type);

CREATE INDEX idx_user_id_notifications ON notifications(user_id);
CREATE INDEX idx_status_notifications ON notifications(status);
CREATE INDEX idx_type_notifications ON notifications(type);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp_users
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_venues
    BEFORE UPDATE ON venues
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_events
    BEFORE UPDATE ON events
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_ticket_categories
    BEFORE UPDATE ON ticket_categories
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_tickets
    BEFORE UPDATE ON tickets
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_orders
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_payment_methods
    BEFORE UPDATE ON payment_methods
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_payments
    BEFORE UPDATE ON payments
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_ticket_queues
    BEFORE UPDATE ON ticket_queues
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_ticket_reservations
    BEFORE UPDATE ON ticket_reservations
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_queue_settings
    BEFORE UPDATE ON queue_settings
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp_notification_templates
    BEFORE UPDATE ON notification_templates
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();