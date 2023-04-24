-- Create company type
DROP TYPE IF EXISTS company_type;
CREATE TYPE company_type AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'SoleProprietorship');

-- Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY,
    name VARCHAR(15) NOT NULL,
    description VARCHAR(3000),
    employees INTEGER NOT NULL,
    registered BOOLEAN NOT NULL,
    type company_type NOT NULL
);