-- Requires psql variable user_seed_password_hash (bcrypt hash for all test users).
-- Abort if not provided or empty to avoid creating users with blank passwords.
\if :{?user_seed_password_hash}
  \if :'user_seed_password_hash' = ''
    \echo 'ERROR: user_seed_password_hash is empty. Provide a bcrypt hash (USER_TOKEN_TEST).'
    \quit 1
  \endif
\else
  \echo 'ERROR: user_seed_password_hash not provided. Export USER_TOKEN_TEST or pass -v user_seed_password_hash=<hash>.'
  \quit 1
\endif

INSERT INTO aion_api.users (name, username, password, email, roles)
VALUES
    ('user_name_1','user1',:'user_seed_password_hash','user1@aion.com','user'),
    ('user_name_2','user2',:'user_seed_password_hash','user2@aion.com','user'),
    ('user_name_3','user3',:'user_seed_password_hash','user3@aion.com','user'),
    ('user_name_4','user4',:'user_seed_password_hash','user4@aion.com','user'),
    ('user_name_5','user5',:'user_seed_password_hash','user5@aion.com','user'),
    ('user_name_6','user6',:'user_seed_password_hash','user6@aion.com','user')
ON CONFLICT (username) DO NOTHING;
