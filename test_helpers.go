package otf

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func NewTestOrganization(t *testing.T) *Organization {
	org, err := NewOrganization(OrganizationCreateOptions{
		Name: String(uuid.NewString()),
	})
	require.NoError(t, err)
	return org
}

func NewTestWorkspace(t *testing.T, org *Organization, opts WorkspaceCreateOptions) *Workspace {
	if opts.Name == "" {
		opts.Name = uuid.NewString()
	}
	ws, err := NewWorkspace(org, opts)
	require.NoError(t, err)
	return ws
}

func NewTestConfigurationVersion(t *testing.T, ws *Workspace, opts ConfigurationVersionCreateOptions) *ConfigurationVersion {
	cv, err := NewConfigurationVersion(ws.ID(), opts)
	require.NoError(t, err)
	return cv
}

func NewTestUser(t *testing.T, opts ...NewUserOption) *User {
	return NewUser(uuid.NewString(), opts...)
}

func NewTestTeam(t *testing.T, org *Organization, opts ...NewTeamOption) *Team {
	return NewTeam(uuid.NewString(), org, opts...)
}
