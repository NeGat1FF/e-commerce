CREATE OR REPLACE FUNCTION reset_password(
    reset_token TEXT,
    new_password TEXT
) RETURNS TEXT AS $$
DECLARE
    user_uid UUID;
BEGIN
    -- Step 1: Validate the token
    SELECT user_id INTO user_uid
    FROM password_resets
    WHERE token = reset_token
      AND expires_at > NOW();

    IF NOT FOUND THEN
        RETURN 'Invalid or expired reset token.';
    END IF;

    -- Step 2: Check if the new password matches the current password
    IF EXISTS (
        SELECT 1 FROM users WHERE id = user_uid AND password = new_password
    ) THEN
        RETURN 'New password cannot be the same as the old password.';
    END IF;

    -- Step 3: Update the password
    UPDATE users
    SET password = new_password
    WHERE id = user_uid;

    -- Step 4: Delete the token after use
    DELETE FROM password_resets WHERE token = reset_token;


    RETURN 'Password updated successfully.';
END;
$$ LANGUAGE plpgsql;
