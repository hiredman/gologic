package gologic

func mzero () *Stream {
        return nil
}

func unit (a S) *Stream {
        var x = new(Stream)
        x.first = a
        x.rest = func () *Stream {
                return mzero()
        }
        return x
}

func choice (a S, s func () *Stream) *Stream {
        var x = new(Stream)
        x.first = a
        x.rest = s
        return x
}

func stream_concat(s1 *Stream, s2 func () *Stream) *Stream {
        if s1 == mzero() {
                return s2()
        } else {
                return choice(s1.first, func () *Stream {
                        return stream_concat(s1.rest(), s2)
                })
        }
}

func stream_interleave (s1 *Stream, s2 *Stream) *Stream {
        if s1 == mzero() {
                return s2
        } else {
                return choice(s1.first, func () *Stream {
                        return stream_interleave(s2,s1.rest())
                })
        }
}
