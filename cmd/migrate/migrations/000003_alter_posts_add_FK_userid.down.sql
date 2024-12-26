DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_user'
    ) THEN
        ALTER TABLE posts DROP CONSTRAINT fk_user;
    END IF;
END $$;