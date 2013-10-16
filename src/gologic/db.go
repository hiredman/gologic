package gologic

// a db is sort of like an clojure agent containing a linked list
func Db () DB {
        var db *DBCons = nil
        c := make(chan func (*DBCons) *DBCons)
        go func () {
                for {
                        f := <- c
                        db = f(db)
                }
        }()
        return DB{c}
}

func (db DB) write(r db_record) DB {
        db.c <- func (d *DBCons) *DBCons {
                return &DBCons{r,d}
        }
	return db
}

func (d DB) Deref() DBValue {
        x := make(chan *DBCons)
        d.c <- func (d *DBCons) *DBCons {
                x <- d
                return d
        }
        return DBValue{<- x}
}

func (d DB) Assert (entity interface{}, attribute interface{}, value interface{}) {
        d.write(db_record{Entity:entity,Attribute:attribute,Value:value})
}

func (d DBValue) Find (entity interface{}, attribute interface{}, value interface{}) Goal {
        r := db_record{Entity:entity,Attribute:attribute,Value:value}
        return func (s S) R {
                g := Fail()
                for e := d.d; e != nil; e = e.cdr {
                        g = Or(g,Unify(r,e.car))
                }
                return g(s)
        }
}
