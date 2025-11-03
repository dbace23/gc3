 

DROP TABLE IF EXISTS exercise_logs CASCADE;
DROP TABLE IF EXISTS exercises CASCADE;
DROP TABLE IF EXISTS workouts CASCADE;
DROP TABLE IF EXISTS users CASCADE;
 
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL UNIQUE,
  username VARCHAR(50) NOT NULL UNIQUE,
  full_name VARCHAR(120) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  weight_kg INTEGER NOT NULL CHECK (weight_kg > 0),
  height_cm INTEGER NOT NULL CHECK (height_cm > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
 
CREATE TABLE workouts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_workouts_user ON workouts(user_id);
 
CREATE TABLE exercises (
  id BIGSERIAL PRIMARY KEY,
  workout_id BIGINT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
  name VARCHAR(120) NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_exercises_workout ON exercises(workout_id);

 
CREATE TABLE exercise_logs (
  id BIGSERIAL PRIMARY KEY,
  exercise_id BIGINT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  set_count INTEGER NOT NULL CHECK (set_count > 0),
  rep_count INTEGER NOT NULL CHECK (rep_count > 0),
  weight INTEGER NOT NULL CHECK (weight >= 0),
  logged_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_logs_user ON exercise_logs(user_id);
CREATE INDEX idx_logs_exercise ON exercise_logs(exercise_id);

 
INSERT INTO users (email, username, full_name, password_hash, weight_kg, height_cm)
VALUES ('demo@gym.com', 'demo', 'Demo User', 'bcrypt-hash-here', 70, 175);

INSERT INTO workouts (user_id, name, description)
VALUES (1, 'Push Day', 'Chest, shoulders, triceps');

INSERT INTO exercises (workout_id, name, description)
VALUES (1, 'Bench Press', 'Flat bench press with barbell');

INSERT INTO exercise_logs (exercise_id, user_id, set_count, rep_count, weight, logged_at)
VALUES (1, 1, 3, 10, 60, NOW());
