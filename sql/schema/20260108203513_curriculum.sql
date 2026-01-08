-- +goose Up
-- +goose StatementBegin
CREATE TABLE curriculum_topics (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  description TEXT,
  color_code TEXT
);

CREATE TABLE curriculum_plans (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL
);

CREATE TABLE curriculum_plan_slots (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  plan_id UUID NOT NULL REFERENCES curriculum_plans(id) ON DELETE CASCADE,
  topic_id UUID NOT NULL REFERENCES curriculum_topics(id) ON DELETE CASCADE,
  sequence_order INTEGER NOT NULL,
  duration INTERVAL NOT NULL
);

CREATE TABLE active_curriculum_periods (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  topic_id UUID NOT NULL REFERENCES curriculum_topics(id) ON DELETE CASCADE,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  EXCLUDE USING gist (daterange(start_date, end_date, '[]') WITH &&)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE active_curriculum_periods;
DROP TABLE curriculum_plan_slots;
DROP TABLE curriculum_plans;
DROP TABLE curriculum_topics;
-- +goose StatementEnd
