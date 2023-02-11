package grpcserver

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestGRPCServer(t *testing.T) {
	ctx := context.Background()
	grpcSrv := New(&sync.WaitGroup{})
	grpcSrv.Start(ctx)

	grpcBase := grpc.NewServer()
	api.RegisterSSDServer(grpcBase, grpcSrv)

	dialer := func() func(context.Context, string) (net.Conn, error) {
		lis := bufconn.Listen(1024 * 1024)

		go func() {
			if err := grpcBase.Serve(lis); err != nil {
				require.NoError(t, err)
			}
		}()

		return func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}
	}

	getConn := func(ctx context.Context) *grpc.ClientConn {
		conn, err := grpc.DialContext(ctx, "",
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithContextDialer(dialer()))
		if err != nil {
			t.Fatal(err)
		}
		return conn
	}

	t.Run("dummy_sensor", func(t *testing.T) {
		conn := getConn(context.Background())
		defer conn.Close()

		var stats api.STATS
		client := api.NewSSDClient(conn)
		stats |= api.STATS_DUMMY
		req := &api.Request{
			N:       durationpb.New(10 * time.Millisecond),
			M:       durationpb.New(20 * time.Microsecond),
			Bitmask: stats,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		stream, err := client.Subsribe(ctx, req)
		require.NoError(t, err)

		resp, err := stream.Recv()
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
