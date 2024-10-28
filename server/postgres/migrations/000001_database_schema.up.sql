-- Users table
CREATE TABLE IF NOT EXISTS "users" (
  "id" SERIAL NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "email" VARCHAR(255) NOT NULL,
  "password" TEXT NOT NULL,
  "avatar" TEXT,
  "location" TEXT,
  "bio" TEXT,
  "interests" TEXT,
  "role" VARCHAR(50),
  "is_email_verified" BOOLEAN DEFAULT FALSE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "users_pkey" PRIMARY KEY ("id")
);

-- Relationships table
CREATE TABLE IF NOT EXISTS "relationships" (
  "id" SERIAL NOT NULL,
  "follower_id" INTEGER,
  "following_id" INTEGER,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  UNIQUE("follower_id", "following_id"),
  CONSTRAINT "relationships_pkey" PRIMARY KEY ("id")
);

-- Communities table
CREATE TABLE IF NOT EXISTS "communities" (
  "id" SERIAL NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "banner" VARCHAR(255),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "communities_pkey" PRIMARY KEY ("id")
);

-- Community members joining table
CREATE TABLE IF NOT EXISTS "community_members" (
  "community_id" INTEGER,
  "user_id" INTEGER,
  "is_moderator" BOOLEAN DEFAULT FALSE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL
);

-- Banned users joining table
CREATE TABLE IF NOT EXISTS "community_banned_users" (
  "community_id" INTEGER,
  "user_id" INTEGER,
  "banned_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Posts table
CREATE TABLE IF NOT EXISTS "posts" (
  "id" SERIAL NOT NULL,
  "content" TEXT,
  "file_url" VARCHAR(255),
  "community_id" INTEGER,
  "user_id" INTEGER,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "posts_pkey" PRIMARY KEY ("id")
);

-- Post likes joining table
CREATE TABLE IF NOT EXISTS "post_likes" (
  "post_id" INTEGER,
  "user_id" INTEGER,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("post_id", "user_id")
);

-- Saved posts joining table
CREATE TABLE IF NOT EXISTS "saved_posts" (
  "user_id" INTEGER,
  "post_id" INTEGER,
  "saved_at" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("user_id", "post_id")
);

-- Comments table
CREATE TABLE IF NOT EXISTS "comments" (
  "id" SERIAL NOT NULL,
  "body" TEXT NOT NULL,
  "user_id" INTEGER,
  "post_id" INTEGER,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "comments_pkey" PRIMARY KEY ("id")
);

-- Rules table
CREATE TABLE IF NOT EXISTS "rules" (
  "id" SERIAL NOT NULL,
  "community_id" INTEGER,
  "rule" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "rules_pkey" PRIMARY KEY ("id")
);

-- Reports table
CREATE TABLE IF NOT EXISTS "reports" (
  "id" SERIAL NOT NULL,
  "post_id" INTEGER,
  "community_id" INTEGER,
  "reported_by" INTEGER,
  "report_reason" TEXT,
  "report_date" TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT "reports_pkey" PRIMARY KEY ("id")
);

-- Pending posts table
CREATE TABLE IF NOT EXISTS "pending_posts" (
  "id" SERIAL NOT NULL,
  "content" TEXT,
  "file_url" VARCHAR(255),
  "community_id" INTEGER,
  "user_id" INTEGER,
  "status" VARCHAR(50),
  "comnfirmation_token" VARCHAR(255),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "pending_posts_pkey" PRIMARY KEY ("id")
);

-- Suspicious logins table
CREATE TABLE IF NOT EXISTS "suspicious_logins" (
  "id" SERIAL NOT NULL,
  "user_id" INTEGER,
  "email" VARCHAR(255),
  "ip" VARCHAR(45),
  "country" VARCHAR(100),
  "city" VARCHAR(100),
  "browser" VARCHAR(255),
  "platform" VARCHAR(255),
  "os" VARCHAR(100),
  "device" VARCHAR(255),
  "device_type" VARCHAR(255),
  "unverified_attempts" INTEGER DEFAULT 0,
  "is_trusted" BOOLEAN DEFAULT FALSE,
  "is_blocked" BOOLEAN DEFAULT FALSE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "suspicious_logins_pkey" PRIMARY KEY ("id")
);

-- User preferences table
CREATE TABLE IF NOT EXISTS "preferences" (
  "id" SERIAL NOT NULL,
  "user_id" INTEGER,
  "enable_context_based_auth" BOOLEAN DEFAULT FALSE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT "preferences_pkey" PRIMARY KEY ("id")
);

-- Tokens table
CREATE TABLE IF NOT EXISTS "tokens" (
  "id" SERIAL NOT NULL,
  "user_id" INTEGER,
  "refresh_token" VARCHAR(255),
  "access_token" VARCHAR(255),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "tokens_pkey" PRIMARY KEY ("id")
);

-- Admin users table
CREATE TABLE IF NOT EXISTS "admins" (
  "id" SERIAL NOT NULL,
  "username" VARCHAR(255) NOT NULL,
  "password" TEXT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "admins_pkey" PRIMARY KEY ("id")
);

-- Admins tokens table
CREATE TABLE IF NOT EXISTS "admin_tokens" (
  "id" SERIAL NOT NULL,
  "admin_id" INTEGER,
  "access_token" VARCHAR(255),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "admin_tokens_pkey" PRIMARY KEY ("id")
);

-- Logs table
CREATE TABLE IF NOT EXISTS "logs" (
  "id" SERIAL NOT NULL,
  "email" VARCHAR(255),
  "context" TEXT,
  "message" TEXT,
  "type" VARCHAR(50),
  "level" VARCHAR(50),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "logs_pkey" PRIMARY KEY ("id")
);

-- Configuration table
CREATE TABLE IF NOT EXISTS "configs" (
  "id" SERIAL NOT NULL,
  "use_perspective_api" BOOLEAN DEFAULT FALSE,
  "category_filtering_service_provider" VARCHAR(255),
  "category_filtering_request_timeout" INTEGER,
  CONSTRAINT "configs_pkey" PRIMARY KEY ("id")
);

-- Emails table
CREATE TABLE IF NOT EXISTS "emails" (
  "id" SERIAL NOT NULL,
  "email" VARCHAR(255),
  "verification_code" VARCHAR(100),
  "message_id" VARCHAR(255),
  "for" VARCHAR(255),
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "expires_at" TIMESTAMPTZ NOT NULL,
  CONSTRAINT "emails_pkey" PRIMARY KEY ("id")
);

CREATE UNIQUE INDEX "users_email_key" ON "users"("email");

CREATE UNIQUE INDEX "admins_username_key" ON "admins"("username");

ALTER TABLE "relationships" ADD CONSTRAINT "relationships_follower_id_fkey" FOREIGN KEY ("follower_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "relationships" ADD CONSTRAINT "relationships_following_id_fkey" FOREIGN KEY ("following_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "community_members" ADD CONSTRAINT "community_members_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "community_members" ADD CONSTRAINT "community_members_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "community_banned_users" ADD CONSTRAINT "community_banned_users_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "community_banned_users" ADD CONSTRAINT "community_banned_users_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "posts" ADD CONSTRAINT "posts_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "posts" ADD CONSTRAINT "posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "post_likes" ADD CONSTRAINT "post_likes_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "posts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "post_likes" ADD CONSTRAINT "post_likes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "saved_posts" ADD CONSTRAINT "saved_posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "saved_posts" ADD CONSTRAINT "saved_posts_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "posts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "comments" ADD CONSTRAINT "comments_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "comments" ADD CONSTRAINT "comments_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "posts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "rules" ADD CONSTRAINT "rules_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "reports" ADD CONSTRAINT "reports_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "posts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "reports" ADD CONSTRAINT "reports_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "reports" ADD CONSTRAINT "reports_reported_by_fkey" FOREIGN KEY ("reported_by") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "pending_posts" ADD CONSTRAINT "pending_posts_community_id_fkey" FOREIGN KEY ("community_id") REFERENCES "communities"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
ALTER TABLE "pending_posts" ADD CONSTRAINT "pending_posts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "suspicious_logins" ADD CONSTRAINT "suspicious_logins_user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "preferences" ADD CONSTRAINT "preferences_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "tokens" ADD CONSTRAINT "tokens_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

ALTER TABLE "admin_tokens" ADD CONSTRAINT "admin_tokens_admin_id_fkey" FOREIGN KEY ("admin_id") REFERENCES "admins"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
