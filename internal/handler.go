package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/crazybolillo/eryth/pkg/model"
	"github.com/glauth/ldap"
	"log/slog"
	"net"
	"net/http"
)

type Handler struct {
	Config  *Config
	Session *Session
}

func (h *Handler) verifyCredentials(dn, password string) error {
	backend, ok := h.Config.Backends[dn]
	if !ok {
		return fmt.Errorf("no backend found for dn: %s", dn)
	}

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s", dn, password)))
	expect, err := hex.DecodeString(backend.Password)
	if err != nil {
		return fmt.Errorf("failed to decode backend's password: %w", err)
	}

	if !bytes.Equal(expect, hash[:]) {
		return fmt.Errorf("invalid credentials")
	}

	return nil
}

func (h *Handler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	slog.Info("Processing Bind", "dn", bindDN, "remote", conn.RemoteAddr())

	err := h.verifyCredentials(bindDN, bindSimplePw)
	if err != nil {
		slog.Info("Bind rejected", "reason", err)
		return ldap.LDAPResultInvalidCredentials, nil
	}
	h.Session.Add(conn, bindDN)

	return ldap.LDAPResultSuccess, nil
}

func (h *Handler) Unbind(bindDN string, conn net.Conn) (ldap.LDAPResultCode, error) {
	slog.Info("Processing Unbind", "dn", bindDN, "remote", conn.RemoteAddr())
	h.Session.Remove(conn)

	return ldap.LDAPResultSuccess, nil
}

func (h *Handler) Search(boundDN string, req ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	slog.Info("Processing Search", "dn", boundDN, "remote", conn.RemoteAddr())

	// Return LDAP Root DES
	if boundDN == "" {
		return ldap.ServerSearchResult{
			Entries: []*ldap.Entry{
				{
					DN: "",
					Attributes: []*ldap.EntryAttribute{
						{Name: "supportedLDAPVersion", Values: []string{"3"}},
					},
				},
			},
			ResultCode: ldap.LDAPResultSuccess,
		}, nil
	}

	dn := h.Session.Get(conn)
	if dn != boundDN {
		slog.Info("Unauthorized search request performed", "remote", conn.RemoteAddr())
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, nil
	}

	backend, ok := h.Config.Backends[dn]
	if !ok {
		slog.Info("No backend found, aborting search", "dn", dn)
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, nil
	}

	res, err := http.Get(backend.URL)
	if err != nil {
		slog.Error("Failed to request information from backend", "url", backend.URL, "reason", err)
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, err
	}

	var page model.ContactPage
	decode := json.NewDecoder(res.Body)
	err = decode.Decode(&page)
	if err != nil {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, err
	}

	entries := make([]*ldap.Entry, page.Retrieved)
	for idx, contact := range page.Contacts {
		entries[idx] = &ldap.Entry{
			DN: fmt.Sprintf("uname=%s,%s", contact.ID, req.BaseDN),
			Attributes: []*ldap.EntryAttribute{
				{Name: "cn", Values: []string{contact.Name}},
				{Name: "telephoneNumber", Values: []string{contact.Phone}},
			},
		}
	}

	return ldap.ServerSearchResult{
		Entries:    entries,
		ResultCode: ldap.LDAPResultSuccess,
	}, nil
}

func (h *Handler) Abandon(boundDN string, conn net.Conn) error {
	slog.Info("Processing Abandon", "dn", boundDN, "remote", conn.RemoteAddr())
	return nil
}
