CREATE TABLE aion_api.days
(
    id            SERIAL PRIMARY KEY,                               -- Unique Day ID
    user_id       INT NOT NULL,                                     -- Foreign Key to User
    date          DATE NOT NULL,                                    -- Day's Date
    mood          TEXT,                                             -- Mood of the day //TODO: transformar para uma nova tabela, adicionando updatedAt (isso vai me ajudar a verificar quando mudo meu estado durante o dia).
    energy        INT CHECK (energy BETWEEN 1 AND 5),               -- Energy (1-5) //TODO: transformar para uma nova tabela, adicionando updatedAt (isso vai me ajudar a verificar quando mudo minha energia durante o dia).
    water_intake  INT CHECK (water_intake >= 0),                     -- Water intake //TODO: transformar para uma nova tabela, adicionando updatedAt (isso vai me ajudar a verificar quando ingiro água durante o dia).
    intention     TEXT,                                             -- Daily intentions //TODO: transformar para uma nova tabela, adicionando updatedAt (isso vai me ajudar a verificar quando mudo minhas intenções durante o dia).
    FOREIGN KEY (user_id) REFERENCES aion_api.users (id) ON DELETE CASCADE -- Relationship to User
);