package log

import (
	"context"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/internal/secret"
	"github.com/ydb-platform/ydb-go-sdk/v3/trace"
)

// Driver makes trace.Driver with logging events from details.
func Driver(l Logger, d trace.Detailer, opts ...Option) (t trace.Driver) {
	return internalDriver(wrapLogger(l, opts...), d)
}

func internalDriver(l Logger, d trace.Detailer) (t trace.Driver) { //nolint:gocyclo
	t.OnResolve = func(
		info trace.DriverResolveStartInfo,
	) func(
		trace.DriverResolveDoneInfo,
	) {
		if d.Details()&trace.DriverResolverEvents == 0 {
			return nil
		}
		ctx := with(context.Background(), TRACE, "ydb", "driver", "resolver", "update")
		target := info.Target
		addresses := info.Resolved
		l.Log(ctx, "start",
			String("target", target),
			Strings("resolved", addresses),
		)

		return func(info trace.DriverResolveDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					String("target", target),
					Strings("resolved", addresses),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					String("target", target),
					Strings("resolved", addresses),
					versionField(),
				)
			}
		}
	}
	t.OnInit = func(info trace.DriverInitStartInfo) func(trace.DriverInitDoneInfo) {
		if d.Details()&trace.DriverEvents == 0 {
			return nil
		}
		endpoint := info.Endpoint
		database := info.Database
		secure := info.Secure
		ctx := with(*info.Context, DEBUG, "ydb", "driver", "resolver", "init")
		l.Log(ctx, "start",
			String("endpoint", endpoint),
			String("database", database),
			Bool("secure", secure),
		)
		start := time.Now()

		return func(info trace.DriverInitDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					String("endpoint", endpoint),
					String("database", database),
					Bool("secure", secure),
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, ERROR), "failed",
					Error(info.Error),
					String("endpoint", endpoint),
					String("database", database),
					Bool("secure", secure),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnClose = func(info trace.DriverCloseStartInfo) func(trace.DriverCloseDoneInfo) {
		if d.Details()&trace.DriverEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "resolver", "close")
		l.Log(ctx, "start")
		start := time.Now()

		return func(info trace.DriverCloseDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnConnDial = func(info trace.DriverConnDialStartInfo) func(trace.DriverConnDialDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "dial")
		endpoint := info.Endpoint
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()

		return func(info trace.DriverConnDialDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					Stringer("endpoint", endpoint),
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnConnStateChange = func(info trace.DriverConnStateChangeStartInfo) func(trace.DriverConnStateChangeDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(context.Background(), TRACE, "ydb", "driver", "conn", "state", "change")
		endpoint := info.Endpoint
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
			Stringer("state", info.State),
		)
		start := time.Now()

		return func(info trace.DriverConnStateChangeDoneInfo) {
			l.Log(ctx, "done",
				Stringer("endpoint", endpoint),
				latencyField(start),
				Stringer("state", info.State),
			)
		}
	}
	t.OnConnPark = func(info trace.DriverConnParkStartInfo) func(trace.DriverConnParkDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "park")
		endpoint := info.Endpoint
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()

		return func(info trace.DriverConnParkDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					Stringer("endpoint", endpoint),
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnConnClose = func(info trace.DriverConnCloseStartInfo) func(trace.DriverConnCloseDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "close")
		endpoint := info.Endpoint
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()

		return func(info trace.DriverConnCloseDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					Stringer("endpoint", endpoint),
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnConnInvoke = func(info trace.DriverConnInvokeStartInfo) func(trace.DriverConnInvokeDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "invoke")
		endpoint := info.Endpoint
		method := string(info.Method)
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
			String("method", method),
		)
		start := time.Now()

		return func(info trace.DriverConnInvokeDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					Stringer("endpoint", endpoint),
					String("method", method),
					latencyField(start),
					Stringer("metadata", metadata(info.Metadata)),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					String("method", method),
					latencyField(start),
					Stringer("metadata", metadata(info.Metadata)),
					versionField(),
				)
			}
		}
	}
	t.OnConnNewStream = func(
		info trace.DriverConnNewStreamStartInfo,
	) func(
		trace.DriverConnNewStreamRecvInfo,
	) func(
		trace.DriverConnNewStreamDoneInfo,
	) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "new", "stream")
		endpoint := info.Endpoint
		method := string(info.Method)
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
			String("method", method),
		)
		start := time.Now()

		return func(info trace.DriverConnNewStreamRecvInfo) func(trace.DriverConnNewStreamDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "intermediate receive",
					Stringer("endpoint", endpoint),
					String("method", method),
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "intermediate fail",
					Error(info.Error),
					Stringer("endpoint", endpoint),
					String("method", method),
					latencyField(start),
					versionField(),
				)
			}

			return func(info trace.DriverConnNewStreamDoneInfo) {
				if info.Error == nil {
					l.Log(ctx, "done",
						Stringer("endpoint", endpoint),
						String("method", method),
						latencyField(start),
						Stringer("metadata", metadata(info.Metadata)),
					)
				} else {
					l.Log(WithLevel(ctx, WARN), "failed",
						Error(info.Error),
						Stringer("endpoint", endpoint),
						String("method", method),
						latencyField(start),
						Stringer("metadata", metadata(info.Metadata)),
						versionField(),
					)
				}
			}
		}
	}
	t.OnConnBan = func(info trace.DriverConnBanStartInfo) func(trace.DriverConnBanDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "ban")
		endpoint := info.Endpoint
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
			NamedError("cause", info.Cause),
			versionField(),
		)
		start := time.Now()

		return func(info trace.DriverConnBanDoneInfo) {
			l.Log(WithLevel(ctx, WARN), "done",
				Stringer("endpoint", endpoint),
				latencyField(start),
				Stringer("state", info.State),
				versionField(),
			)
		}
	}
	t.OnConnAllow = func(info trace.DriverConnAllowStartInfo) func(trace.DriverConnAllowDoneInfo) {
		if d.Details()&trace.DriverConnEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "conn", "allow")
		endpoint := info.Endpoint
		l.Log(ctx, "start",
			Stringer("endpoint", endpoint),
		)
		start := time.Now()

		return func(info trace.DriverConnAllowDoneInfo) {
			l.Log(ctx, "done",
				Stringer("endpoint", endpoint),
				latencyField(start),
				Stringer("state", info.State),
			)
		}
	}
	t.OnRepeaterWakeUp = func(info trace.DriverRepeaterWakeUpStartInfo) func(trace.DriverRepeaterWakeUpDoneInfo) {
		if d.Details()&trace.DriverRepeaterEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "repeater", "wake", "up")
		name := info.Name
		event := info.Event
		l.Log(ctx, "start",
			String("name", name),
			String("event", event),
		)
		start := time.Now()

		return func(info trace.DriverRepeaterWakeUpDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					String("name", name),
					String("event", event),
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, ERROR), "failed",
					Error(info.Error),
					String("name", name),
					String("event", event),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnBalancerInit = func(info trace.DriverBalancerInitStartInfo) func(trace.DriverBalancerInitDoneInfo) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "balancer", "init")
		l.Log(ctx, "start")
		start := time.Now()

		return func(info trace.DriverBalancerInitDoneInfo) {
			l.Log(WithLevel(ctx, INFO), "done",
				latencyField(start),
			)
		}
	}
	t.OnBalancerClose = func(info trace.DriverBalancerCloseStartInfo) func(trace.DriverBalancerCloseDoneInfo) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "balancer", "close")
		l.Log(ctx, "start")
		start := time.Now()

		return func(info trace.DriverBalancerCloseDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					latencyField(start),
				)
			} else {
				l.Log(WithLevel(ctx, WARN), "failed",
					Error(info.Error),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnBalancerChooseEndpoint = func(
		info trace.DriverBalancerChooseEndpointStartInfo,
	) func(
		trace.DriverBalancerChooseEndpointDoneInfo,
	) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "balancer", "choose", "endpoint")
		l.Log(ctx, "start")
		start := time.Now()

		return func(info trace.DriverBalancerChooseEndpointDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					latencyField(start),
					Stringer("endpoint", info.Endpoint),
				)
			} else {
				l.Log(WithLevel(ctx, ERROR), "failed",
					Error(info.Error),
					latencyField(start),
					versionField(),
				)
			}
		}
	}
	t.OnBalancerUpdate = func(
		info trace.DriverBalancerUpdateStartInfo,
	) func(
		trace.DriverBalancerUpdateDoneInfo,
	) {
		if d.Details()&trace.DriverBalancerEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "balancer", "update")
		l.Log(ctx, "start",
			Bool("needLocalDC", info.NeedLocalDC),
		)
		start := time.Now()

		return func(info trace.DriverBalancerUpdateDoneInfo) {
			l.Log(ctx, "done",
				latencyField(start),
				Stringer("endpoints", endpoints(info.Endpoints)),
				Stringer("added", endpoints(info.Added)),
				Stringer("dropped", endpoints(info.Dropped)),
				String("detectedLocalDC", info.LocalDC),
			)
		}
	}
	t.OnGetCredentials = func(info trace.DriverGetCredentialsStartInfo) func(trace.DriverGetCredentialsDoneInfo) {
		if d.Details()&trace.DriverCredentialsEvents == 0 {
			return nil
		}
		ctx := with(*info.Context, TRACE, "ydb", "driver", "credentials", "get")
		l.Log(ctx, "start")
		start := time.Now()

		return func(info trace.DriverGetCredentialsDoneInfo) {
			if info.Error == nil {
				l.Log(ctx, "done",
					latencyField(start),
					String("token", secret.Token(info.Token)),
				)
			} else {
				l.Log(WithLevel(ctx, ERROR), "done",
					Error(info.Error),
					latencyField(start),
					String("token", secret.Token(info.Token)),
					versionField(),
				)
			}
		}
	}

	return t
}
