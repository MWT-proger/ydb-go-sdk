package ydbsql

import (
	"context"
	"crypto/tls"
	"database/sql/driver"
	"fmt"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/meta/credentials"
	internal "github.com/ydb-platform/ydb-go-sdk/v3/internal/table"
	"sync"

	"github.com/ydb-platform/ydb-go-sdk/v3/table"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"

	"github.com/ydb-platform/ydb-go-sdk/v3/config"
	"github.com/ydb-platform/ydb-go-sdk/v3/internal/dial"

	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
)

type ConnectorOption func(*connector)

func WithDialer(d dial.Dialer) ConnectorOption {
	return func(c *connector) {
		c.dialer = d
		if c.dialer.DriverConfig == nil {
			c.dialer.DriverConfig = new(config.Config)
		}
	}
}

func withClient(pool internal.ClientAsPool) ConnectorOption {
	return func(c *connector) {
		c.pool = pool
	}
}

func WithConnectParams(params ydb.ConnectParams) ConnectorOption {
	return func(c *connector) {
		c.endpoint = params.Endpoint()
		if c.dialer.DriverConfig == nil {
			c.dialer.DriverConfig = &config.Config{}
		}
		c.dialer.DriverConfig.Database = params.Database()
		if params.UseTLS() {
			c.dialer.TLSConfig = &tls.Config{}
		} else {
			c.dialer.TLSConfig = nil
		}
	}
}

func WithConnectionString(connection string) ConnectorOption {
	return WithConnectParams(ydb.MustConnectionString(connection))
}

func WithEndpoint(addr string) ConnectorOption {
	return func(c *connector) {
		c.endpoint = addr
	}
}

func WithDriverConfig(config config.Config) ConnectorOption {
	return func(c *connector) {
		*(c.dialer.DriverConfig) = config
	}
}

func WithCredentials(creds credentials.Credentials) ConnectorOption {
	return func(c *connector) {
		c.dialer.DriverConfig.Credentials = creds
	}
}

func WithDatabase(db string) ConnectorOption {
	return func(c *connector) {
		c.dialer.DriverConfig.Database = db
	}
}

func WithTraceDriver(t trace.Driver) ConnectorOption {
	return func(c *connector) {
		c.dialer.DriverConfig.Trace = t
	}
}

func WithTrace(t trace.Table) ConnectorOption {
	return func(c *connector) {
		c.trace = t
	}
}

func WithDefaultTxControl(txControl *table.TransactionControl) ConnectorOption {
	return func(c *connector) {
		c.defaultTxControl = txControl
	}
}

func WithDefaultExecDataQueryOption(opts ...options.ExecuteDataQueryOption) ConnectorOption {
	return func(c *connector) {
		c.dataOpts = append(c.dataOpts, opts...)
	}
}

func WithDefaultExecScanQueryOption(opts ...options.ExecuteScanQueryOption) ConnectorOption {
	return func(c *connector) {
		c.scanOpts = append(c.scanOpts, opts...)
	}
}

func Connector(opts ...ConnectorOption) driver.Connector {
	c := &connector{
		dialer: dial.Dialer{
			DriverConfig: config.New(),
		},
		defaultTxControl: table.TxControl(
			table.BeginTx(
				table.WithSerializableReadWrite(),
			),
			table.CommitTx(),
		),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// USE CONNECTOR ONLY
type connector struct {
	dialer   dial.Dialer
	endpoint string

	trace trace.Table

	mu    sync.Mutex
	ready chan struct{}
	pool  internal.ClientAsPool // Used as a template for created connections.

	defaultTxControl *table.TransactionControl

	dataOpts []options.ExecuteDataQueryOption
	scanOpts []options.ExecuteScanQueryOption
}

func (c *connector) init(ctx context.Context) (err error) {
	// in driver database/sql/sql.go:1228 connect run under mutex, but don't rely on it here
	c.mu.Lock()
	defer c.mu.Unlock()

	// Setup some more on less generic reasonable pool limit to prevent
	// session overflow on the YDB servers.
	//
	// Note that it must be controlled from outside by making
	// database/sql.DB.SetMaxIdleConns() call. Unfortunately, we can not
	// receive that limit here and we do not want to force user to
	// configure it twice (and pass it as an option to connector).
	if c.pool == nil {
		c.pool, err = c.dial(ctx)
	}
	//c.pool.Builder = c.client
	return
}

func (c *connector) dial(ctx context.Context) (internal.ClientAsPool, error) {
	d, err := c.dialer.Dial(ctx, c.endpoint)
	if err != nil {
		if c == nil {
			panic("nil connector")
		}
		if c.dialer.DriverConfig == nil {
			panic("nil driver config")
		}
		if c.dialer.DriverConfig.Credentials == nil {
			panic("nil credentials")
		}
		if stringer, ok := c.dialer.DriverConfig.Credentials.(fmt.Stringer); ok {
			return nil, fmt.Errorf("dial error: %w (credentials: %s)", err, stringer.String())
		}
		return nil, fmt.Errorf("dial error: %w", err)
	}
	return internal.NewClientAsPool(d, ContextTableConfig(ctx)), nil
}

func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	s, err := c.pool.Create(ctx)
	if err != nil {
		return nil, err
	}
	if s == nil {
		panic("ydbsql: abnormal result of pool.Create()")
	}
	return &sqlConn{
		connector: c,
		session:   s,
	}, nil
}

func (c *connector) Driver() driver.Driver {
	return &Driver{c}
}

func (c *connector) unwrap(ctx context.Context) (table.Client, error) {
	if err := c.init(ctx); err != nil {
		return nil, err
	}
	return c.pool, nil
}

// Driver is an adapter to allow the use table client as sql.Driver instance.
// The main purpose of this types is exported is an ability to call Unwrap()
// method on it to receive raw *table.client instance.
type Driver struct {
	c *connector
}

func (d *Driver) Close(ctx context.Context) error {
	return d.c.pool.Close(ctx)
}

// Open returns a new connection to the ydb.
func (d *Driver) Open(string) (driver.Conn, error) {
	return nil, ErrDeprecated
}

func (d *Driver) OpenConnector(string) (driver.Connector, error) {
	return d.c, nil
}

func (d *Driver) Unwrap(ctx context.Context) (table.Client, error) {
	return d.c.unwrap(ctx)
}
