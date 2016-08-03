package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"github.com/miekg/dns"
	"log"
	"sync"
	"time"
)

type cache struct {
	m    map[string]*elem
	lock *sync.RWMutex
	ttl  int64
	max  int
}

type elem struct {
	m *dns.Msg
	t int64
}

func newCache(max int, ttl int64) *cache {
	return &cache{
		max:  max,
		ttl:  ttl,
		m:    map[string]*elem{},
		lock: new(sync.RWMutex),
	}
}

func key(m *dns.Msg) string {
	d := m.Question[0].Name
	b := []byte(d)
	b1 := make([]byte, 4)
	binary.BigEndian.PutUint16(b1[0:], m.Question[0].Qclass)
	binary.BigEndian.PutUint16(b1[2:], m.Question[0].Qtype)
	b = append(b, b1...)
	h := md5.New()
	h.Write(b)
	s1 := hex.EncodeToString(h.Sum(nil))
	return s1
}

func (c *cache) get(m *dns.Msg) *dns.Msg {
	c.lock.RLock()
	defer c.lock.RUnlock()
	k := key(m)
	if m1, ok := c.m[k]; ok {
		t := time.Now().Unix()
		if t < m1.t {
			return m1.m
		}
	}
	return nil
}

func (c *cache) set(m *dns.Msg) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.m) >= c.max {
		log.Printf("clean the old cache")
		c.cleanOld()
	}

	k := key(m)
	c.m[k] = &elem{
		m.Copy(),
		time.Now().Unix() + c.ttl,
	}
}

// must hold the write lock
func (c *cache) cleanOld() {
	t1 := time.Now().Unix()
	for k, v := range c.m {
		if v.t >= t1 {
			delete(c.m, k)
		}
	}
}
