// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AccountsColumns holds the columns for the "accounts" table.
	AccountsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "alias", Type: field.TypeString},
		{Name: "name", Type: field.TypeString, Size: 64},
		{Name: "description", Type: field.TypeString, Size: 256, Default: ""},
		{Name: "role", Type: field.TypeEnum, Enums: []string{"OWNER", "MEMBER"}},
		{Name: "created_date", Type: field.TypeTime},
		{Name: "silo_id", Type: field.TypeUUID},
		{Name: "user_accounts", Type: field.TypeUUID},
	}
	// AccountsTable holds the schema information for the "accounts" table.
	AccountsTable = &schema.Table{
		Name:       "accounts",
		Columns:    AccountsColumns,
		PrimaryKey: []*schema.Column{AccountsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "accounts_silos_members",
				Columns:    []*schema.Column{AccountsColumns[6]},
				RefColumns: []*schema.Column{SilosColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "accounts_users_accounts",
				Columns:    []*schema.Column{AccountsColumns[7]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "account_silo_id_alias",
				Unique:  true,
				Columns: []*schema.Column{AccountsColumns[6], AccountsColumns[1]},
			},
		},
	}
	// IdentitiesColumns holds the columns for the "identities" table.
	IdentitiesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true},
		{Name: "kind", Type: field.TypeString},
		{Name: "verifier", Type: field.TypeString},
		{Name: "name", Type: field.TypeString, Default: ""},
		{Name: "created_date", Type: field.TypeTime},
		{Name: "user_identities", Type: field.TypeUUID},
	}
	// IdentitiesTable holds the schema information for the "identities" table.
	IdentitiesTable = &schema.Table{
		Name:       "identities",
		Columns:    IdentitiesColumns,
		PrimaryKey: []*schema.Column{IdentitiesColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "identities_users_identities",
				Columns:    []*schema.Column{IdentitiesColumns[5]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.NoAction,
			},
		},
	}
	// InvitationsColumns holds the columns for the "invitations" table.
	InvitationsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "invitee", Type: field.TypeString},
		{Name: "created_date", Type: field.TypeTime},
		{Name: "expired_date", Type: field.TypeTime},
		{Name: "accepted_date", Type: field.TypeTime},
		{Name: "declined_date", Type: field.TypeTime},
		{Name: "canceled_date", Type: field.TypeTime},
		{Name: "account_invitations", Type: field.TypeUUID},
		{Name: "silo_invitations", Type: field.TypeUUID},
	}
	// InvitationsTable holds the schema information for the "invitations" table.
	InvitationsTable = &schema.Table{
		Name:       "invitations",
		Columns:    InvitationsColumns,
		PrimaryKey: []*schema.Column{InvitationsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "invitations_accounts_invitations",
				Columns:    []*schema.Column{InvitationsColumns[7]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.NoAction,
			},
			{
				Symbol:     "invitations_silos_invitations",
				Columns:    []*schema.Column{InvitationsColumns[8]},
				RefColumns: []*schema.Column{SilosColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// MembershipsColumns holds the columns for the "memberships" table.
	MembershipsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "role", Type: field.TypeEnum, Enums: []string{"OWNER", "MEMBER"}},
		{Name: "created_date", Type: field.TypeTime},
		{Name: "account_memberships", Type: field.TypeUUID},
		{Name: "team_members", Type: field.TypeUUID},
	}
	// MembershipsTable holds the schema information for the "memberships" table.
	MembershipsTable = &schema.Table{
		Name:       "memberships",
		Columns:    MembershipsColumns,
		PrimaryKey: []*schema.Column{MembershipsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "memberships_accounts_memberships",
				Columns:    []*schema.Column{MembershipsColumns[3]},
				RefColumns: []*schema.Column{AccountsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "memberships_teams_members",
				Columns:    []*schema.Column{MembershipsColumns[4]},
				RefColumns: []*schema.Column{TeamsColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// SilosColumns holds the columns for the "silos" table.
	SilosColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "alias", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString, Size: 64},
		{Name: "description", Type: field.TypeString, Size: 256, Default: ""},
		{Name: "date_created", Type: field.TypeTime},
	}
	// SilosTable holds the schema information for the "silos" table.
	SilosTable = &schema.Table{
		Name:       "silos",
		Columns:    SilosColumns,
		PrimaryKey: []*schema.Column{SilosColumns[0]},
	}
	// TeamsColumns holds the columns for the "teams" table.
	TeamsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "alias", Type: field.TypeString, Unique: true},
		{Name: "name", Type: field.TypeString, Size: 64},
		{Name: "description", Type: field.TypeString, Size: 256, Default: ""},
		{Name: "inter_visibility", Type: field.TypeEnum, Enums: []string{"PRIVATE", "PUBLIC"}},
		{Name: "intra_visibility", Type: field.TypeEnum, Enums: []string{"PRIVATE", "PUBLIC"}},
		{Name: "created_date", Type: field.TypeTime},
		{Name: "silo_id", Type: field.TypeUUID},
	}
	// TeamsTable holds the schema information for the "teams" table.
	TeamsTable = &schema.Table{
		Name:       "teams",
		Columns:    TeamsColumns,
		PrimaryKey: []*schema.Column{TeamsColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "teams_silos_teams",
				Columns:    []*schema.Column{TeamsColumns[7]},
				RefColumns: []*schema.Column{SilosColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
		Indexes: []*schema.Index{
			{
				Name:    "team_silo_id_alias",
				Unique:  true,
				Columns: []*schema.Column{TeamsColumns[7], TeamsColumns[1]},
			},
		},
	}
	// TokensColumns holds the columns for the "tokens" table.
	TokensColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "value", Type: field.TypeString, Unique: true},
		{Name: "type", Type: field.TypeString},
		{Name: "name", Type: field.TypeString, Default: ""},
		{Name: "date_created", Type: field.TypeTime},
		{Name: "date_expired", Type: field.TypeTime},
		{Name: "token_children", Type: field.TypeUUID, Nullable: true},
		{Name: "user_tokens", Type: field.TypeUUID},
	}
	// TokensTable holds the schema information for the "tokens" table.
	TokensTable = &schema.Table{
		Name:       "tokens",
		Columns:    TokensColumns,
		PrimaryKey: []*schema.Column{TokensColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "tokens_tokens_children",
				Columns:    []*schema.Column{TokensColumns[6]},
				RefColumns: []*schema.Column{TokensColumns[0]},
				OnDelete:   schema.SetNull,
			},
			{
				Symbol:     "tokens_users_tokens",
				Columns:    []*schema.Column{TokensColumns[7]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID, Unique: true},
		{Name: "name", Type: field.TypeString, Unique: true},
		{Name: "created_date", Type: field.TypeTime},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AccountsTable,
		IdentitiesTable,
		InvitationsTable,
		MembershipsTable,
		SilosTable,
		TeamsTable,
		TokensTable,
		UsersTable,
	}
)

func init() {
	AccountsTable.ForeignKeys[0].RefTable = SilosTable
	AccountsTable.ForeignKeys[1].RefTable = UsersTable
	IdentitiesTable.ForeignKeys[0].RefTable = UsersTable
	InvitationsTable.ForeignKeys[0].RefTable = AccountsTable
	InvitationsTable.ForeignKeys[1].RefTable = SilosTable
	MembershipsTable.ForeignKeys[0].RefTable = AccountsTable
	MembershipsTable.ForeignKeys[1].RefTable = TeamsTable
	TeamsTable.ForeignKeys[0].RefTable = SilosTable
	TokensTable.ForeignKeys[0].RefTable = TokensTable
	TokensTable.ForeignKeys[1].RefTable = UsersTable
}
