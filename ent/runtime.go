// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/google/uuid"
	"khepri.dev/horus/ent/account"
	"khepri.dev/horus/ent/identity"
	"khepri.dev/horus/ent/invitation"
	"khepri.dev/horus/ent/membership"
	"khepri.dev/horus/ent/silo"
	"khepri.dev/horus/ent/team"
	"khepri.dev/horus/ent/token"
	"khepri.dev/horus/ent/user"
	"khepri.dev/horus/schema"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	accountFields := schema.Account{}.Fields()
	_ = accountFields
	// accountDescAlias is the schema descriptor for alias field.
	accountDescAlias := accountFields[2].Descriptor()
	// account.DefaultAlias holds the default value on creation for the alias field.
	account.DefaultAlias = accountDescAlias.Default.(func() string)
	// account.AliasValidator is a validator for the "alias" field. It is called by the builders before save.
	account.AliasValidator = func() func(string) error {
		validators := accountDescAlias.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(alias string) error {
			for _, fn := range fns {
				if err := fn(alias); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// accountDescName is the schema descriptor for name field.
	accountDescName := accountFields[3].Descriptor()
	// account.NameValidator is a validator for the "name" field. It is called by the builders before save.
	account.NameValidator = func() func(string) error {
		validators := accountDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// accountDescDescription is the schema descriptor for description field.
	accountDescDescription := accountFields[4].Descriptor()
	// account.DefaultDescription holds the default value on creation for the description field.
	account.DefaultDescription = accountDescDescription.Default.(string)
	// account.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	account.DescriptionValidator = accountDescDescription.Validators[0].(func(string) error)
	// accountDescCreatedDate is the schema descriptor for created_date field.
	accountDescCreatedDate := accountFields[6].Descriptor()
	// account.DefaultCreatedDate holds the default value on creation for the created_date field.
	account.DefaultCreatedDate = accountDescCreatedDate.Default.(func() time.Time)
	// accountDescID is the schema descriptor for id field.
	accountDescID := accountFields[0].Descriptor()
	// account.DefaultID holds the default value on creation for the id field.
	account.DefaultID = accountDescID.Default.(func() uuid.UUID)
	identityFields := schema.Identity{}.Fields()
	_ = identityFields
	// identityDescKind is the schema descriptor for kind field.
	identityDescKind := identityFields[1].Descriptor()
	// identity.KindValidator is a validator for the "kind" field. It is called by the builders before save.
	identity.KindValidator = identityDescKind.Validators[0].(func(string) error)
	// identityDescVerifier is the schema descriptor for verifier field.
	identityDescVerifier := identityFields[2].Descriptor()
	// identity.VerifierValidator is a validator for the "verifier" field. It is called by the builders before save.
	identity.VerifierValidator = identityDescVerifier.Validators[0].(func(string) error)
	// identityDescName is the schema descriptor for name field.
	identityDescName := identityFields[3].Descriptor()
	// identity.DefaultName holds the default value on creation for the name field.
	identity.DefaultName = identityDescName.Default.(string)
	// identityDescCreatedDate is the schema descriptor for created_date field.
	identityDescCreatedDate := identityFields[4].Descriptor()
	// identity.DefaultCreatedDate holds the default value on creation for the created_date field.
	identity.DefaultCreatedDate = identityDescCreatedDate.Default.(func() time.Time)
	// identityDescID is the schema descriptor for id field.
	identityDescID := identityFields[0].Descriptor()
	// identity.IDValidator is a validator for the "id" field. It is called by the builders before save.
	identity.IDValidator = identityDescID.Validators[0].(func(string) error)
	invitationFields := schema.Invitation{}.Fields()
	_ = invitationFields
	// invitationDescInvitee is the schema descriptor for invitee field.
	invitationDescInvitee := invitationFields[1].Descriptor()
	// invitation.InviteeValidator is a validator for the "invitee" field. It is called by the builders before save.
	invitation.InviteeValidator = invitationDescInvitee.Validators[0].(func(string) error)
	// invitationDescCreatedDate is the schema descriptor for created_date field.
	invitationDescCreatedDate := invitationFields[2].Descriptor()
	// invitation.DefaultCreatedDate holds the default value on creation for the created_date field.
	invitation.DefaultCreatedDate = invitationDescCreatedDate.Default.(func() time.Time)
	// invitationDescID is the schema descriptor for id field.
	invitationDescID := invitationFields[0].Descriptor()
	// invitation.DefaultID holds the default value on creation for the id field.
	invitation.DefaultID = invitationDescID.Default.(func() uuid.UUID)
	membershipFields := schema.Membership{}.Fields()
	_ = membershipFields
	// membershipDescCreatedDate is the schema descriptor for created_date field.
	membershipDescCreatedDate := membershipFields[2].Descriptor()
	// membership.DefaultCreatedDate holds the default value on creation for the created_date field.
	membership.DefaultCreatedDate = membershipDescCreatedDate.Default.(func() time.Time)
	// membershipDescID is the schema descriptor for id field.
	membershipDescID := membershipFields[0].Descriptor()
	// membership.DefaultID holds the default value on creation for the id field.
	membership.DefaultID = membershipDescID.Default.(func() uuid.UUID)
	siloFields := schema.Silo{}.Fields()
	_ = siloFields
	// siloDescAlias is the schema descriptor for alias field.
	siloDescAlias := siloFields[1].Descriptor()
	// silo.AliasValidator is a validator for the "alias" field. It is called by the builders before save.
	silo.AliasValidator = func() func(string) error {
		validators := siloDescAlias.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(alias string) error {
			for _, fn := range fns {
				if err := fn(alias); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// siloDescName is the schema descriptor for name field.
	siloDescName := siloFields[2].Descriptor()
	// silo.NameValidator is a validator for the "name" field. It is called by the builders before save.
	silo.NameValidator = func() func(string) error {
		validators := siloDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// siloDescDescription is the schema descriptor for description field.
	siloDescDescription := siloFields[3].Descriptor()
	// silo.DefaultDescription holds the default value on creation for the description field.
	silo.DefaultDescription = siloDescDescription.Default.(string)
	// silo.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	silo.DescriptionValidator = siloDescDescription.Validators[0].(func(string) error)
	// siloDescCreatedDate is the schema descriptor for created_date field.
	siloDescCreatedDate := siloFields[4].Descriptor()
	// silo.DefaultCreatedDate holds the default value on creation for the created_date field.
	silo.DefaultCreatedDate = siloDescCreatedDate.Default.(func() time.Time)
	// siloDescID is the schema descriptor for id field.
	siloDescID := siloFields[0].Descriptor()
	// silo.DefaultID holds the default value on creation for the id field.
	silo.DefaultID = siloDescID.Default.(func() uuid.UUID)
	teamFields := schema.Team{}.Fields()
	_ = teamFields
	// teamDescAlias is the schema descriptor for alias field.
	teamDescAlias := teamFields[2].Descriptor()
	// team.DefaultAlias holds the default value on creation for the alias field.
	team.DefaultAlias = teamDescAlias.Default.(func() string)
	// team.AliasValidator is a validator for the "alias" field. It is called by the builders before save.
	team.AliasValidator = func() func(string) error {
		validators := teamDescAlias.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(alias string) error {
			for _, fn := range fns {
				if err := fn(alias); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// teamDescName is the schema descriptor for name field.
	teamDescName := teamFields[3].Descriptor()
	// team.NameValidator is a validator for the "name" field. It is called by the builders before save.
	team.NameValidator = func() func(string) error {
		validators := teamDescName.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(name string) error {
			for _, fn := range fns {
				if err := fn(name); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	// teamDescDescription is the schema descriptor for description field.
	teamDescDescription := teamFields[4].Descriptor()
	// team.DefaultDescription holds the default value on creation for the description field.
	team.DefaultDescription = teamDescDescription.Default.(string)
	// team.DescriptionValidator is a validator for the "description" field. It is called by the builders before save.
	team.DescriptionValidator = teamDescDescription.Validators[0].(func(string) error)
	// teamDescCreatedDate is the schema descriptor for created_date field.
	teamDescCreatedDate := teamFields[7].Descriptor()
	// team.DefaultCreatedDate holds the default value on creation for the created_date field.
	team.DefaultCreatedDate = teamDescCreatedDate.Default.(func() time.Time)
	// teamDescID is the schema descriptor for id field.
	teamDescID := teamFields[0].Descriptor()
	// team.DefaultID holds the default value on creation for the id field.
	team.DefaultID = teamDescID.Default.(func() uuid.UUID)
	tokenFields := schema.Token{}.Fields()
	_ = tokenFields
	// tokenDescValue is the schema descriptor for value field.
	tokenDescValue := tokenFields[1].Descriptor()
	// token.ValueValidator is a validator for the "value" field. It is called by the builders before save.
	token.ValueValidator = tokenDescValue.Validators[0].(func(string) error)
	// tokenDescType is the schema descriptor for type field.
	tokenDescType := tokenFields[2].Descriptor()
	// token.TypeValidator is a validator for the "type" field. It is called by the builders before save.
	token.TypeValidator = tokenDescType.Validators[0].(func(string) error)
	// tokenDescName is the schema descriptor for name field.
	tokenDescName := tokenFields[3].Descriptor()
	// token.DefaultName holds the default value on creation for the name field.
	token.DefaultName = tokenDescName.Default.(string)
	// tokenDescDateCreated is the schema descriptor for date_created field.
	tokenDescDateCreated := tokenFields[4].Descriptor()
	// token.DefaultDateCreated holds the default value on creation for the date_created field.
	token.DefaultDateCreated = tokenDescDateCreated.Default.(func() time.Time)
	// tokenDescID is the schema descriptor for id field.
	tokenDescID := tokenFields[0].Descriptor()
	// token.DefaultID holds the default value on creation for the id field.
	token.DefaultID = tokenDescID.Default.(func() uuid.UUID)
	userFields := schema.User{}.Fields()
	_ = userFields
	// userDescName is the schema descriptor for name field.
	userDescName := userFields[1].Descriptor()
	// user.NameValidator is a validator for the "name" field. It is called by the builders before save.
	user.NameValidator = userDescName.Validators[0].(func(string) error)
	// userDescCreatedDate is the schema descriptor for created_date field.
	userDescCreatedDate := userFields[2].Descriptor()
	// user.DefaultCreatedDate holds the default value on creation for the created_date field.
	user.DefaultCreatedDate = userDescCreatedDate.Default.(func() time.Time)
	// userDescID is the schema descriptor for id field.
	userDescID := userFields[0].Descriptor()
	// user.DefaultID holds the default value on creation for the id field.
	user.DefaultID = userDescID.Default.(func() uuid.UUID)
}
