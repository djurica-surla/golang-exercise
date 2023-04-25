-- Create company type
DROP TYPE IF EXISTS company_type;
CREATE TYPE company_type AS ENUM ('Corporation', 'NonProfit', 'Cooperative', 'SoleProprietorship');

-- Create companies table
CREATE TABLE IF NOT EXISTS companies (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(15) UNIQUE NOT NULL,
    description VARCHAR(3000) NOT NULL,
    employees INTEGER NOT NULL,
    registered BOOLEAN NOT NULL,
    type company_type NOT NULL
);