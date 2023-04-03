package util

import (
	"fmt"
	"io"
	"time"
)

// New New
// rate 速度，rate/每秒
func New(rate int64) *Limiter {
	return &Limiter{
		rate:  rate,
		count: 0,
		t:     time.Now(),
	}
}

// Reader 返回一个带有Limiter的io.Reader
func Reader(r io.Reader, l *Limiter) io.Reader {
	return &reader{
		r: r,
		l: l,
	}
}

// ReadSeeker 返回一个带有Limiter的io.ReadSeeker
func ReadSeeker(rs io.ReadSeeker, l *Limiter) io.ReadSeeker {
	return &readSeeker{
		reader: reader{
			r: rs,
			l: l,
		},
		s: rs,
	}
}

// Writer 返回一个带有Limiter的io.Writer
func Writer(w io.Writer, l *Limiter) io.Writer {
	return &writer{
		w: w,
		l: l,
	}
}

// Limiter 速度限制器
type Limiter struct {
	rate  int64
	count int64 // 最大8G
	t     time.Time
}

// Wait 传入需要处理的数量，计算并等待需要经过的时间
func (l *Limiter) Wait(count int) {

}

type reader struct {
	r io.Reader
	l *Limiter
}

// Read Read
func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	r.l.Wait(n)
	return n, err
}

type readSeeker struct {
	reader
	s io.Seeker
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return rs.s.Seek(offset, whence)
}

type writer struct {
	w io.Writer
	l *Limiter
}

// Write Write
func (w *writer) Write(buf []byte) (int, error) {
	a := time.Now()
	uin := 1024
	nl := len(buf)
	flag := false

	now := time.Now()
	count, ln := 0, 0
	//start, end := 0, 0
	for {
		var ret []byte
		if (count+1)*uin < nl {
			ret = buf[count*uin : (count+1)*uin]
		} else {
			ret = buf[count*uin : nl]
			flag = true
		}
		count += 1
		ln += len(ret)
		w.w.Write(ret)
		if int64(ln) > w.l.rate {
			subM := time.Now().Sub(now)
			t := 1000 - subM.Milliseconds()
			if t > 0 {
				time.Sleep(time.Duration(t) * time.Millisecond)
				now = time.Now()
				ln = 0
			}
		}
		if flag == true {
			break
		}
	}
	fmt.Println(time.Now().Sub(a).Seconds())
	return nl, nil
}
