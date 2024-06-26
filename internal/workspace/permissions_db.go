package workspace

import (
	"context"

	"github.com/tofutf/tofutf/internal"
	"github.com/tofutf/tofutf/internal/rbac"
	"github.com/tofutf/tofutf/internal/sql"
	"github.com/tofutf/tofutf/internal/sql/pggen"
)

func (db *pgdb) SetWorkspacePermission(ctx context.Context, workspaceID, teamID string, role rbac.Role) error {
	return db.Query(ctx, func(ctx context.Context, q pggen.Querier) error {
		_, err := q.UpsertWorkspacePermission(ctx, pggen.UpsertWorkspacePermissionParams{
			WorkspaceID: sql.String(workspaceID),
			TeamID:      sql.String(teamID),
			Role:        sql.String(role.String()),
		})
		if err != nil {
			return sql.Error(err)
		}

		return nil
	})
}

func (db *pgdb) GetWorkspacePolicy(ctx context.Context, workspaceID string) (internal.WorkspacePolicy, error) {
	return sql.Query(ctx, db.Pool, func(ctx context.Context, q pggen.Querier) (internal.WorkspacePolicy, error) {
		// Retrieve not only permissions but the workspace too, so that:
		// (1) we ensure that workspace exists and return not found if not
		// (2) we retrieve the name of the organization, which is part of a policy
		ws, err := q.FindWorkspaceByID(ctx, sql.String(workspaceID))
		if err != nil {
			return internal.WorkspacePolicy{}, sql.Error(err)
		}

		perms, err := q.FindWorkspacePermissionsByWorkspaceID(ctx, sql.String(workspaceID))
		if err != nil {
			return internal.WorkspacePolicy{}, sql.Error(err)
		}

		policy := internal.WorkspacePolicy{
			Organization:      ws.OrganizationName.String,
			WorkspaceID:       workspaceID,
			GlobalRemoteState: ws.GlobalRemoteState.Bool,
		}
		for _, perm := range perms {
			role, err := rbac.WorkspaceRoleFromString(perm.Role.String)
			if err != nil {
				return internal.WorkspacePolicy{}, err
			}
			policy.Permissions = append(policy.Permissions, internal.WorkspacePermission{
				TeamID: perm.TeamID.String,
				Role:   role,
			})
		}
		return policy, nil
	})
}

func (db *pgdb) UnsetWorkspacePermission(ctx context.Context, workspaceID, team string) error {
	return db.Query(ctx, func(ctx context.Context, q pggen.Querier) error {
		_, err := q.DeleteWorkspacePermissionByID(ctx, sql.String(workspaceID), sql.String(team))
		if err != nil {
			return sql.Error(err)
		}

		return nil
	})
}
