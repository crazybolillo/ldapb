package internal

import (
	"github.com/glauth/ldap"
	"net"
)

type Handler struct{}

func (h *Handler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	return ldap.LDAPResultSuccess, nil
}

func (h *Handler) Unbind(bindDN string, conn net.Conn) (ldap.LDAPResultCode, error) {
	return ldap.LDAPResultSuccess, nil
}

func (h *Handler) Search(boundDN string, req ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	return ldap.ServerSearchResult{}, nil
}

func (h *Handler) Abandon(boundDN string, conn net.Conn) error {
	return nil
}
