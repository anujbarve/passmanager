// internal/session/session.go
package session

import (
	"sync"
	"time"

	"passmanager/internal/crypto"
	"passmanager/internal/database"
)

type Session struct {
	mu              sync.RWMutex
	isAuthenticated bool
	cryptoService   *crypto.CryptoService
	dbClient        *database.PocketBaseClient
	lastActivity    time.Time
	timeout         time.Duration
	salt            []byte
}

var (
	currentSession *Session
	once           sync.Once
)

func GetSession() *Session {
	once.Do(func() {
		currentSession = &Session{
			timeout: 5 * time.Minute,
		}
	})
	return currentSession
}

func (s *Session) Login(client *database.PocketBaseClient, cryptoSvc *crypto.CryptoService, salt []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.isAuthenticated = true
	s.dbClient = client
	s.cryptoService = cryptoSvc
	s.salt = salt
	s.lastActivity = time.Now()
}

func (s *Session) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cryptoService != nil {
		s.cryptoService.SecureClear()
	}
	s.isAuthenticated = false
	s.cryptoService = nil
	s.dbClient = nil
	s.salt = nil
}

func (s *Session) IsAuthenticated() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isAuthenticated {
		return false
	}

	// Check for timeout
	if time.Since(s.lastActivity) > s.timeout {
		s.mu.RUnlock()
		s.Logout()
		s.mu.RLock()
		return false
	}

	return true
}

func (s *Session) UpdateActivity() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastActivity = time.Now()
}

func (s *Session) GetCrypto() *crypto.CryptoService {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.lastActivity = time.Now()
	return s.cryptoService
}

func (s *Session) GetDB() *database.PocketBaseClient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	s.lastActivity = time.Now()
	return s.dbClient
}

func (s *Session) GetSalt() []byte {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.salt
}

func (s *Session) SetTimeout(duration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timeout = duration
}

func (s *Session) GetTimeRemaining() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	remaining := s.timeout - time.Since(s.lastActivity)
	if remaining < 0 {
		return 0
	}
	return remaining
}