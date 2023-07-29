CREATE TABLE IF NOT EXISTS parents(
  "id" serial PRIMARY KEY,
  "uid" int NOT NULL,
  "account_number" varchar(255) NOT NULL,
  "nik" varchar(255) NOT NULL UNIQUE,
  "username" varchar(255) NOT NULL UNIQUE,
  "email" varchar(255) NOT NULL UNIQUE,
  "password" varchar(255) NOT NULL,
  "pin" varchar(255) NOT NULL,
  "phone_number" varchar(255) NOT NULL,
  "full_name" varchar(255) NOT NULL,
  "domisili" varchar(255) NOT NULL,
  "tanggal_lahir" varchar(255) NOT NULL,
  "jenis_kelamin" smallint NOT NULL, -- 0: laki-laki, 1: perempuan
  "alamat" varchar(255) NOT NULL,
  "rt_rw" varchar(255) NOT NULL,
  "kelurahan" varchar(255) NOT NULL,
  "kecamatan" varchar(255) NOT NULL,
  "pekerjaan" varchar(255) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS kids(
  "id" serial PRIMARY KEY,
  "parent_id" int NOT NULL,
  "account_number" varchar(255) NOT NULL,
  "nik" varchar(255) NOT NULL,
  "full_name" varchar(255) NOT NULL,
  "domisili" varchar(255) NOT NULL,
  "tanggal_lahir" varchar(255) NOT NULL,
  "jenis_kelamin" smallint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  FOREIGN KEY (parent_id) REFERENCES parents(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS goals(
  "id" serial PRIMARY KEY,
  "kid_id" int NOT NULL,
  "account_number" varchar(255) NOT NULL,
  "title" varchar(255) NOT NULL,
  "target_amount" decimal NOT NULL,
  "status" smallint NOT NULL, -- 0: ongoing, 1: achieved, 2: overdue, 3: canceled
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  FOREIGN KEY (kid_id) REFERENCES kids(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS quests(
  "id" serial PRIMARY KEY,
  "parent_id" int NOT NULL,
  "title" varchar(255) NOT NULL,
  "description" varchar(255) NOT NULL,
  "reward" decimal NOT NULL,
  "status" smallint NOT NULL, -- 0: available, 1: ongoing, 2: done, 3: canceled, 4: expired
  "start_date" timestamptz NOT NULL,
  "end_date" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  FOREIGN KEY (parent_id) REFERENCES parents(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS kid_assigned_quests(
  "kid_id" int NOT NULL,
  "quest_id" int NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  PRIMARY KEY (kid_id, quest_id),
  FOREIGN KEY (kid_id) REFERENCES kids(id) ON DELETE CASCADE,
  FOREIGN KEY (quest_id) REFERENCES quests(id) ON DELETE CASCADE
);

