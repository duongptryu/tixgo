-- Create templates table
CREATE TABLE IF NOT EXISTS templates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    subject VARCHAR(500),
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('email', 'sms', 'push')),
    status VARCHAR(50) NOT NULL DEFAULT 'draft' CHECK (status IN ('active', 'inactive', 'draft')),
    variables TEXT[], -- Array of variable names used in the template
    description TEXT,
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_templates_slug ON templates(slug);
CREATE INDEX IF NOT EXISTS idx_templates_type ON templates(type);
CREATE INDEX IF NOT EXISTS idx_templates_status ON templates(status);
CREATE INDEX IF NOT EXISTS idx_templates_created_by ON templates(created_by);
CREATE INDEX IF NOT EXISTS idx_templates_created_at ON templates(created_at);

-- Add foreign key constraint to users table (assuming it exists)
-- ALTER TABLE templates ADD CONSTRAINT fk_templates_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE;

-- Add comments for documentation
COMMENT ON TABLE templates IS 'Email/SMS/Push notification templates for the platform';
COMMENT ON COLUMN templates.name IS 'Human-readable name of the template';
COMMENT ON COLUMN templates.slug IS 'Unique slug identifier for the template';
COMMENT ON COLUMN templates.subject IS 'Subject line for email templates (can be empty for SMS/Push)';
COMMENT ON COLUMN templates.content IS 'Template content with placeholders';
COMMENT ON COLUMN templates.type IS 'Type of template: email, sms, or push';
COMMENT ON COLUMN templates.status IS 'Template status: active, inactive, or draft';
COMMENT ON COLUMN templates.variables IS 'Array of variable names that can be used in the template';
COMMENT ON COLUMN templates.description IS 'Description of what this template is used for';
COMMENT ON COLUMN templates.created_by IS 'ID of the user who created this template'; 