// Code generated by "protoc-gen-entpb". DO NOT EDIT

package horus

import (
	uuid "github.com/google/uuid"
)

func AccountById(k uuid.UUID) *GetAccountRequest {
	return &GetAccountRequest{Id: k[:]}
}
func AccountByIdV(k []byte) *GetAccountRequest {
	return &GetAccountRequest{Id: k}
}
func IdentityById(k uuid.UUID) *GetIdentityRequest {
	return &GetIdentityRequest{Id: k[:]}
}
func IdentityByIdV(k []byte) *GetIdentityRequest {
	return &GetIdentityRequest{Id: k}
}
func InvitationById(k uuid.UUID) *GetInvitationRequest {
	return &GetInvitationRequest{Id: k[:]}
}
func InvitationByIdV(k []byte) *GetInvitationRequest {
	return &GetInvitationRequest{Id: k}
}
func MembershipById(k uuid.UUID) *GetMembershipRequest {
	return &GetMembershipRequest{Id: k[:]}
}
func MembershipByIdV(k []byte) *GetMembershipRequest {
	return &GetMembershipRequest{Id: k}
}
func SiloById(k uuid.UUID) *GetSiloRequest {
	return &GetSiloRequest{Key: &GetSiloRequest_Id{Id: k[:]}}
}
func SiloByIdV(k []byte) *GetSiloRequest {
	return &GetSiloRequest{Key: &GetSiloRequest_Id{Id: k}}
}
func SiloByAlias(k string) *GetSiloRequest {
	return &GetSiloRequest{Key: &GetSiloRequest_Alias{Alias: k}}
}
func TeamById(k uuid.UUID) *GetTeamRequest {
	return &GetTeamRequest{Id: k[:]}
}
func TeamByIdV(k []byte) *GetTeamRequest {
	return &GetTeamRequest{Id: k}
}
func TokenById(k uuid.UUID) *GetTokenRequest {
	return &GetTokenRequest{Key: &GetTokenRequest_Id{Id: k[:]}}
}
func TokenByIdV(k []byte) *GetTokenRequest {
	return &GetTokenRequest{Key: &GetTokenRequest_Id{Id: k}}
}
func TokenByValue(k string) *GetTokenRequest {
	return &GetTokenRequest{Key: &GetTokenRequest_Value{Value: k}}
}
func UserById(k uuid.UUID) *GetUserRequest {
	return &GetUserRequest{Key: &GetUserRequest_Id{Id: k[:]}}
}
func UserByIdV(k []byte) *GetUserRequest {
	return &GetUserRequest{Key: &GetUserRequest_Id{Id: k}}
}
func UserByAlias(k string) *GetUserRequest {
	return &GetUserRequest{Key: &GetUserRequest_Alias{Alias: k}}
}
