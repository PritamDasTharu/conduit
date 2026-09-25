package database

import (
	"fmt"

	authModel "github.com/PritamDasTharu/conduit/internal/modules/authentication/model"
	channelModel "github.com/PritamDasTharu/conduit/internal/modules/channel/model"
	messageModel "github.com/PritamDasTharu/conduit/internal/modules/message/model"
	workspaceModel "github.com/PritamDasTharu/conduit/internal/modules/workspace/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return fmt.Errorf("enable pgcrypto extension: %w", err)
	}

	// GORM AutoMigrate schema
	err := db.AutoMigrate(
		&authModel.User{},
		&workspaceModel.Workspace{},
		&workspaceModel.Permission{},
		&workspaceModel.Role{},
		&workspaceModel.RolePermission{},
		&workspaceModel.WorkspaceMember{},
		&channelModel.Channel{},
		&channelModel.ChannelMember{},
		&messageModel.Message{},
	)
	if err != nil {
		return fmt.Errorf("gorm automigrate: %w", err)
	}

	partialIndexes := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_users_active_email 
		 ON users (LOWER(email)) WHERE deleted_at IS NULL;`,

		// Unique active workspace slug 
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_workspaces_active_slug 
		 ON workspaces (LOWER(slug)) WHERE deleted_at IS NULL;`,

		// One active workspace membership per user
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_workspace_members_active_user 
		 ON workspace_members (workspace_id, user_id) WHERE deleted_at IS NULL;`,

		// Unique active role name within a workspace
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_active_workspace_name 
		 ON roles (workspace_id, LOWER(name)) WHERE deleted_at IS NULL;`,

		// Unique active channel slug within a workspace
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_channels_active_workspace_slug 
		 ON channels (workspace_id, LOWER(slug)) WHERE deleted_at IS NULL;`,

		// One active channel membership per user
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_channel_members_active_user 
		 ON channel_members (channel_id, user_id) WHERE left_at IS NULL;`,
	}

	for _, stmt := range partialIndexes {
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("create partial index: %w", err)
		}
	}

	tenantSafetyConstraints := []string{
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_workspace_id_id 
		 ON roles (workspace_id, id);`,

		`DO $$ 
		BEGIN 
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint WHERE conname = 'fk_workspace_members_role_same_workspace'
			) THEN 
				ALTER TABLE workspace_members 
				ADD CONSTRAINT fk_workspace_members_role_same_workspace 
				FOREIGN KEY (workspace_id, role_id) 
				REFERENCES roles(workspace_id, id) 
				ON DELETE RESTRICT; 
			END IF; 
		END $$;`,
	}

	for _, stmt := range tenantSafetyConstraints {
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("apply tenant safety constraint: %w", err)
		}
	}

	triggerStatements := []string{
		`CREATE OR REPLACE FUNCTION set_updated_at_column()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;`,

		`DO $$
		DECLARE
			tbl TEXT;
			tables TEXT[] := ARRAY['users', 'workspaces', 'roles', 'permissions', 'workspace_members', 'channels', 'messages'];
		BEGIN
			FOREACH tbl IN ARRAY tables LOOP
				IF NOT EXISTS (
					SELECT 1 FROM pg_trigger WHERE tgname = 'trg_set_updated_at_' || tbl
				) THEN
					EXECUTE format('
						CREATE TRIGGER trg_set_updated_at_%I
						BEFORE UPDATE ON %I
						FOR EACH ROW
						EXECUTE FUNCTION set_updated_at_column();', tbl, tbl);
				END IF;
			END LOOP;
		END $$;`,
	}

	for _, stmt := range triggerStatements {
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("apply updated_at trigger: %w", err)
		}
	}

	return nil
}
