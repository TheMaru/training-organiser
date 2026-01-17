-- name: CreateCurriculumTopic :one
INSERT INTO curriculum_topics (name, description, color_code)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListCurriculumTopics :many
SELECT * FROM curriculum_topics ORDER BY name;

-- name: CreateCurriculumPlan :one
INSERT INTO curriculum_plans (name)
VALUES ($1)
RETURNING *;

-- name: ListCurriculumPlans :many
SELECT * FROM curriculum_plans ORDER BY created_at DESC;

-- name: AddSlotToPlan :one
INSERT INTO curriculum_plan_slots (plan_id, topic_id, sequence_order, duration)
VALUES ($1, $2, $3, $4::interval)
RETURNING *;

-- name: GetSlotsForPlan :many
SELECT * FROM curriculum_plan_slots
WHERE plan_id = $1
ORDER BY sequence_order;

-- name: CreateActivePeriod :one
INSERT INTO active_curriculum_periods (topic_id, start_date, end_date)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetActivePeriods :many
SELECT
  p.*,
  t.name as topic_name,
  t.description as topic_description,
  t.color_code
FROM active_curriculum_periods p
JOIN curriculum_topics t ON p.topic_id = t.id
WHERE p.end_date >= $1 AND p.start_date <= $2
ORDER BY p.start_date;

