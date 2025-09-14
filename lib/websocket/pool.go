package websocket

import (
	"app/lib/logger"
	"context"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"go.uber.org/zap"
)

type WebsocketPool struct {
	conns    chan *websocket.Conn
	url      string
	maxConns int
	mu       sync.RWMutex
}

func NewWebsocketPool(url string, maxConns int) *WebsocketPool {
	return &WebsocketPool{
		conns:    make(chan *websocket.Conn, maxConns),
		url:      url,
		maxConns: maxConns,
	}
}

func (wsPool *WebsocketPool) getConnection(ctx context.Context) (*websocket.Conn, error) {
	select {
	case conn := <-wsPool.conns:
		// Reuse pooled connection without health check - we'll validate it during actual write
		// If connection is broken, write will fail and we'll create a fresh connection
		return conn, nil
	default:
		return wsPool.createConnection(ctx)
	}
}

func (wsPool *WebsocketPool) createConnection(ctx context.Context) (*websocket.Conn, error) {
	c, _, err := websocket.Dial(ctx, wsPool.url, nil)
	if err != nil {
		logger.LogError(ctx, "error websocket.Dial", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"websocket", "createConnection"}),
		}...)
		return nil, err
	}
	return c, nil
}

func (wsPool *WebsocketPool) returnConnection(conn *websocket.Conn) {
	select {
	case wsPool.conns <- conn:
	default:
		// Pool is full, close the connection
		conn.CloseNow()
	}
}

func (p *WebsocketPool) SendMessage(ctx context.Context, message Message) error {
	conn, err := p.getConnection(ctx)
	if err != nil {
		return err
	}

	writeCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := wsjson.Write(writeCtx, conn, message); err != nil {
		conn.CloseNow() // Don't return broken connection to pool
		logger.LogError(ctx, "error first try wsjson.Write, retry with fresh connection", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"websocket", "SendMessage"}),
		}...)

		// Retry with fresh connection
		newConn, err := p.createConnection(ctx)
		if err != nil {
			return err
		}

		err = wsjson.Write(writeCtx, newConn, message)
		if err != nil {
			newConn.CloseNow() // Don't return broken connection to pool
			logger.LogError(ctx, "error retry wsjson.Write", []zap.Field{
				zap.Error(err),
				zap.Strings("tags", []string{"websocket", "SendMessage"}),
			}...)
			return err
		}

		// Retry success, return the connection to pool
		p.returnConnection(newConn)
		return nil
	}

	// First attempt success, return the connection to pool
	p.returnConnection(conn)
	return nil
}

func (p *WebsocketPool) Close() {
	for {
		select {
		case conn := <-p.conns:
			conn.CloseNow()
		default:
			return
		}
	}
}
