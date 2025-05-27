DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS notification_templates;
DROP TABLE IF EXISTS queue_settings;
DROP TABLE IF EXISTS ticket_reservations;
DROP TABLE IF EXISTS ticket_queues;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS ticket_categories;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS venues;
DROP TABLE IF EXISTS users;

-- Drop custom ENUM types
DROP TYPE IF EXISTS notification_status_enum;
DROP TYPE IF EXISTS notification_type_enum;
DROP TYPE IF EXISTS reservation_status_enum;
DROP TYPE IF EXISTS queue_status_enum;
DROP TYPE IF EXISTS refund_status_enum;
DROP TYPE IF EXISTS payment_status_enum;
DROP TYPE IF EXISTS payment_type_enum;
DROP TYPE IF EXISTS order_status_enum;
DROP TYPE IF EXISTS ticket_status_enum;
DROP TYPE IF EXISTS ticket_category_type_enum;
DROP TYPE IF EXISTS event_status_enum;
DROP TYPE IF EXISTS event_type_enum;
DROP TYPE IF EXISTS venue_type_enum;
DROP TYPE IF EXISTS user_status_enum;
DROP TYPE IF EXISTS user_type_enum;

-- Drop triggers for updated_at
DROP TRIGGER IF EXISTS set_timestamp_users ON users;
DROP TRIGGER IF EXISTS set_timestamp_venues ON venues;
DROP TRIGGER IF EXISTS set_timestamp_events ON events;
DROP TRIGGER IF EXISTS set_timestamp_ticket_categories ON ticket_categories;
DROP TRIGGER IF EXISTS set_timestamp_tickets ON tickets;
DROP TRIGGER IF EXISTS set_timestamp_orders ON orders;
DROP TRIGGER IF EXISTS set_timestamp_payment_methods ON payment_methods;
DROP TRIGGER IF EXISTS set_timestamp_payments ON payments;
DROP TRIGGER IF EXISTS set_timestamp_ticket_queues ON ticket_queues;
DROP TRIGGER IF EXISTS set_timestamp_ticket_reservations ON ticket_reservations;
DROP TRIGGER IF EXISTS set_timestamp_queue_settings ON queue_settings;
DROP TRIGGER IF EXISTS set_timestamp_notification_templates ON notification_templates;

-- Drop the trigger function
DROP FUNCTION IF EXISTS trigger_set_timestamp();
