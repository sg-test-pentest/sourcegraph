ALTER TABLE github_apps
    ADD COLUMN IF NOT EXISTS app_url TEXT NOT NULL DEFAULT '';