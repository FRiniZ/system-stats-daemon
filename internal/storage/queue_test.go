package storage

import (
	"fmt"
	"testing"
	"time"

	api "github.com/FRiniZ/system-stats-daemon/api/stub"
	"github.com/stretchr/testify/require"
)

type DummySensor struct {
	AvgValue int32
}

func (d *DummySensor) Add(a *DummySensor) {
	d.AvgValue += a.AvgValue
}

func (d *DummySensor) Div(n int32) {
	d.AvgValue /= n
}

func (d *DummySensor) MakeResponse() *api.Responce {
	return nil
}

func TestQueue(t *testing.T) {

	for i := 0; i < 10; i++ {
		queue := New(15)
		t.Run(fmt.Sprintf("loadtesting_N%d", i), func(t *testing.T) {
			t.Parallel()
			for i := 0; i < 100; i++ {
				sensor := DummySensor{
					AvgValue: 1,
				}
				queue.Push(&sensor, time.Now())
			}
			count := int(0)
			in := queue.GetElementsAfter(time.Now().Add(-10 * time.Minute))
			for s := range in {
				if _, ok := s.(*DummySensor); ok {
					count++
				}
			}
			require.Equal(t, count, int(15))
		})
	}

	t.Run("checking_limits_by_size", func(t *testing.T) {
		queue := New(2)
		sensor1 := DummySensor{
			AvgValue: 1,
		}
		sensor2 := DummySensor{
			AvgValue: 2,
		}
		sensor3 := DummySensor{
			AvgValue: 3,
		}

		queue.Push(&sensor1, time.Now())
		queue.Push(&sensor2, time.Now())
		queue.Push(&sensor3, time.Now())

		require.Equal(t, int(2), queue.size)
		for e := queue.list.Front(); e != nil; e = e.Next() {
			elm := e.Value.(element)
			sensor := elm.data.(*DummySensor)
			require.NotEqual(t, sensor1.AvgValue, sensor.AvgValue)
		}

		queue.SetSize("testing", 3)
		queue.Push(&sensor1, time.Now())

		elm := queue.list.Front().Value.(element)
		sensor := elm.data.(*DummySensor)
		require.Equal(t, sensor1.AvgValue, sensor.AvgValue)

		elm = queue.list.Front().Next().Value.(element)
		sensor = elm.data.(*DummySensor)
		require.Equal(t, sensor3.AvgValue, sensor.AvgValue)

		elm = queue.list.Back().Value.(element)
		sensor = elm.data.(*DummySensor)
		require.Equal(t, sensor2.AvgValue, sensor.AvgValue)
	})

	t.Run("checking_GetElementsAfter", func(t *testing.T) {
		queue := New(3)
		sensor1 := DummySensor{
			AvgValue: 1,
		}
		sensor2 := DummySensor{
			AvgValue: 2,
		}
		sensor3 := DummySensor{
			AvgValue: 3,
		}

		curTime := time.Now()
		queue.Push(&sensor1, curTime.Add(-15*time.Second))
		queue.Push(&sensor2, curTime.Add(-10*time.Second))
		queue.Push(&sensor3, curTime.Add(-5*time.Second))

		count := int(0)
		for e := queue.list.Front(); e != nil; e = e.Next() {
			count++
		}
		require.Equal(t, count, int(3))

		count = int(0)
		in := queue.GetElementsAfter(curTime.Add(-6 * time.Second))
		for s := range in {
			if _, ok := s.(*DummySensor); ok {
				count++
			}
		}
		require.Equal(t, count, int(1))

		count = int(0)
		in = queue.GetElementsAfter(curTime.Add(-11 * time.Second))
		for s := range in {
			if _, ok := s.(*DummySensor); ok {
				count++
			}
		}
		require.Equal(t, count, int(2))

		count = int(0)
		in = queue.GetElementsAfter(curTime.Add(-16 * time.Second))
		for s := range in {
			if _, ok := s.(*DummySensor); ok {
				count++
			}
		}
		require.Equal(t, count, int(3))
	})
}
